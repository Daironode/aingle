use crate::conductor::ConductorHandle;
use crate::core::ribosome::ZomeCallInvocation;
use crate::core::workflow::incoming_dgd_ops_workflow::IncomingDgdOpsWorkspace;
use crate::test_utils::host_fn_caller::*;
use crate::test_utils::new_invocation;
use crate::test_utils::new_zome_call;
use crate::test_utils::setup_app;
use crate::test_utils::wait_for_integration;
use fallible_iterator::FallibleIterator;
use aingle_hash::AnyDgdHash;
use aingle_hash::DgdOpHash;
use aingle_hash::EntryHash;
use aingle_hash::HeaderHash;
use aingle_lmdb::env::EnvironmentWrite;
use aingle_lmdb::fresh_reader_test;
use aingle_middleware_bytes::SerializedBytes;
use aingle_state::dgd_op_integration::IntegratedDgdOpsValue;
use aingle_state::element_buf::ElementBuf;
use aingle_state::validation_db::ValidationLimboValue;
use aingle_types::prelude::*;
use aingle_wasm_test_utils::TestWasm;

use aingle_zome_types::Entry;
use aingle_zome_types::ValidationStatus;
use std::convert::TryFrom;
use std::convert::TryInto;
use std::time::Duration;
use tracing::*;

#[tokio::test(threaded_scheduler)]
async fn app_validation_workflow_test() {
    observability::test_run_open().ok();

    let dna_file = DnaFile::new(
        DnaDef {
            name: "app_validation_workflow_test".to_string(),
            uuid: "ba1d046d-ce29-4778-914b-47e6010d2faf".to_string(),
            properties: SerializedBytes::try_from(()).unwrap(),
            zomes: vec![
                TestWasm::Validate.into(),
                TestWasm::ValidateLink.into(),
                TestWasm::Create.into(),
            ]
            .into(),
        },
        vec![
            TestWasm::Validate.into(),
            TestWasm::ValidateLink.into(),
            TestWasm::Create.into(),
        ],
    )
    .await
    .unwrap();

    let alice_agent_id = fake_agent_pubkey_1();
    let alice_cell_id = CellId::new(dna_file.dna_hash().to_owned(), alice_agent_id.clone());
    let alice_installed_cell = InstalledCell::new(alice_cell_id.clone(), "alice_handle".into());

    let bob_agent_id = fake_agent_pubkey_2();
    let bob_cell_id = CellId::new(dna_file.dna_hash().to_owned(), bob_agent_id.clone());
    let bob_installed_cell = InstalledCell::new(bob_cell_id.clone(), "bob_handle".into());

    let (_tmpdir, _app_api, handle) = setup_app(
        vec![(
            "test_app",
            vec![(alice_installed_cell, None), (bob_installed_cell, None)],
        )],
        vec![dna_file.clone()],
    )
    .await;

    let expected_count = run_test(
        alice_cell_id.clone(),
        bob_cell_id.clone(),
        handle.clone(),
        &dna_file,
    )
    .await;
    run_test_entry_def_id(
        alice_cell_id,
        bob_cell_id,
        handle.clone(),
        &dna_file,
        expected_count,
    )
    .await;

    let shutdown = handle.take_shutdown_handle().await.unwrap();
    handle.shutdown().await;
    shutdown.await.unwrap();
}

// These are the expected invalid ops
fn expected_invalid_entry(
    (hash, i, el): &(DgdOpHash, IntegratedDgdOpsValue, Element),
    line: u32,
    invalid_header_hash: &HeaderHash,
    invalid_entry_hash: &AnyDgdHash,
) -> bool {
    let s = format!("\nline:{}\n{:?}\n{:?}\n{:?}", line, hash, i, el);
    match &i.op {
        // A Store entry that matches these hashes
        DgdOpLight::StoreEntry(hh, _, eh)
            if eh == invalid_entry_hash && hh == invalid_header_hash =>
        {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s)
        }
        // And the store element
        DgdOpLight::StoreElement(hh, _, _) if hh == invalid_header_hash => {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s);
        }
        _ => return false,
    }
    true
}

// All others must be valid
fn others((hash, i, el): &(DgdOpHash, IntegratedDgdOpsValue, Element), line: u32) {
    let s = format!("\nline:{}\n{:?}\n{:?}\n{:?}", line, hash, i, el);
    match &i.op {
        // Register agent activity will be invalid if the previous header is invalid
        // This is very hard to track in these tests and this op also doesn't
        // go through app validation so it's more productive to skip it
        DgdOpLight::RegisterAgentActivity(_, _) => {}
        _ => assert_eq!(i.validation_status, ValidationStatus::Valid, "{}", s),
    }
}

// Now we expect an invalid link
fn expected_invalid_link(
    (hash, i, el): &(DgdOpHash, IntegratedDgdOpsValue, Element),
    line: u32,
    invalid_link_hash: &HeaderHash,
) -> bool {
    let s = format!("\nline:{}\n{:?}\n{:?}\n{:?}", line, hash, i, el);
    match &i.op {
        // Invalid link
        DgdOpLight::RegisterAddLink(hh, _) if hh == invalid_link_hash => {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s)
        }
        // The store element for this CreateLink header is also rejected
        DgdOpLight::StoreElement(hh, _, _) if hh == invalid_link_hash => {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s)
        }
        _ => return false,
    }
    true
}

// Now we're trying to remove an invalid link
fn expected_invalid_remove_link(
    (hash, i, el): &(DgdOpHash, IntegratedDgdOpsValue, Element),
    line: u32,
    invalid_remove_hash: &HeaderHash,
) -> bool {
    let s = format!("\nline:{}\n{:?}\n{:?}\n{:?}", line, hash, i, el);

    // To make it simple we want to skip this op
    if let DgdOpLight::RegisterAgentActivity(_, _) = &i.op {
        return false;
    }

    // Get the hash of the entry that makes the link invalid
    let sb = SerializedBytes::try_from(&MaybeLinkable::NeverLinkable).unwrap();
    let invalid_link_entry_hash = EntryHash::with_data_sync(&Entry::app(sb).unwrap());

    // Link adds with these base / target are invalid
    if let Header::CreateLink(la) = el.header() {
        if invalid_link_entry_hash == la.base_address
            || invalid_link_entry_hash == la.target_address
        {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s);
            return true;
        }
    }
    match &i.op {
        // The store element for the DeleteLink is invalid
        DgdOpLight::StoreElement(hh, _, _) if hh == invalid_remove_hash => {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s)
        }
        // The remove link op is also invalid
        DgdOpLight::RegisterRemoveLink(hh, _) if hh == invalid_remove_hash => {
            assert_eq!(i.validation_status, ValidationStatus::Rejected, "{}", s)
        }
        _ => return false,
    }
    true
}

async fn run_test(
    alice_cell_id: CellId,
    bob_cell_id: CellId,
    handle: ConductorHandle,
    dna_file: &DnaFile,
) -> usize {
    // Check if the correct number of ops are integrated
    // every 100 ms for a maximum of 10 seconds but early exit
    // if they are there.
    let num_attempts = 100;
    let delay_per_attempt = Duration::from_millis(100);

    let invocation =
        new_zome_call(&bob_cell_id, "always_validates", (), TestWasm::Validate).unwrap();
    handle.call_zome(invocation).await.unwrap().unwrap();

    // Integration should have 3 ops in it
    // Plus another 16 for genesis + init
    // Plus 2 for Cap Grant
    let expected_count = 3 + 16 + 2;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for (_, i, _) in &int {
            assert_eq!(i.validation_status, ValidationStatus::Valid);
        }

        assert_eq!(int.len(), expected_count);
    }

    let (invalid_header_hash, invalid_entry_hash) =
        commit_invalid(&bob_cell_id, &handle, dna_file).await;
    let invalid_entry_hash: AnyDgdHash = invalid_entry_hash.into();

    // Integration should have 3 ops in it
    // StoreEntry should be invalid.
    // RegisterAgentActivity doesn't run app validation
    // So they will be valid.
    let expected_count = 3 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            if !expected_invalid_entry(v, line!(), &invalid_header_hash, &invalid_entry_hash) {
                others(v, line!())
            }
        }

        assert_eq!(int.len(), expected_count);
    }

    let invocation =
        new_zome_call(&bob_cell_id, "add_valid_link", (), TestWasm::ValidateLink).unwrap();
    handle.call_zome(invocation).await.unwrap().unwrap();

    // Integration should have 6 ops in it
    let expected_count = 6 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            if !expected_invalid_entry(v, line!(), &invalid_header_hash, &invalid_entry_hash) {
                others(v, line!())
            }
        }
        assert_eq!(int.len(), expected_count);
    }
    let invocation =
        new_invocation(&bob_cell_id, "add_invalid_link", (), TestWasm::ValidateLink).unwrap();
    let invalid_link_hash: HeaderHash =
        call_zome_directly(&bob_cell_id, &handle, dna_file, invocation)
            .await
            .decode()
            .unwrap();

    // Integration should have 9 ops in it
    let expected_count = 9 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            if !expected_invalid_entry(v, line!(), &invalid_header_hash, &invalid_entry_hash)
                && !expected_invalid_link(v, line!(), &invalid_link_hash)
            {
                others(v, line!())
            }
        }
        assert_eq!(int.len(), expected_count);
    }

    let invocation = new_invocation(
        &bob_cell_id,
        "remove_valid_link",
        (),
        TestWasm::ValidateLink,
    )
    .unwrap();
    call_zome_directly(&bob_cell_id, &handle, dna_file, invocation).await;

    // Integration should have 9 ops in it
    let expected_count = 9 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            if !expected_invalid_entry(v, line!(), &invalid_header_hash, &invalid_entry_hash)
                && !expected_invalid_link(v, line!(), &invalid_link_hash)
            {
                others(v, line!())
            }
        }
        assert_eq!(int.len(), expected_count);
    }

    let invocation = new_invocation(
        &bob_cell_id,
        "remove_invalid_link",
        (),
        TestWasm::ValidateLink,
    )
    .unwrap();
    let invalid_remove_hash: HeaderHash =
        call_zome_directly(&bob_cell_id, &handle, dna_file, invocation)
            .await
            .decode()
            .unwrap();

    // Integration should have 12 ops in it
    let expected_count = 12 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            if !expected_invalid_entry(v, line!(), &invalid_header_hash, &invalid_entry_hash)
                && !expected_invalid_link(v, line!(), &invalid_link_hash)
                && !expected_invalid_remove_link(v, line!(), &invalid_remove_hash)
            {
                others(v, line!())
            }
        }
        assert_eq!(int.len(), expected_count);
    }
    expected_count
}

/// 1. Commits an entry with validate_create_entry_<EntryDefId> callback
/// 2. The callback rejects the entry proving that it actually ran.
/// 3. Reject only Post with "Banana" as the String to show it doesn't
///    affect other entries.
async fn run_test_entry_def_id(
    alice_cell_id: CellId,
    bob_cell_id: CellId,
    handle: ConductorHandle,
    dna_file: &DnaFile,
    expected_count: usize,
) {
    // Check if the correct number of ops are integrated
    // every 100 ms for a maximum of 10 seconds but early exit
    // if they are there.
    let num_attempts = 100;
    let delay_per_attempt = Duration::from_millis(100);

    let (invalid_header_hash, invalid_entry_hash) =
        commit_invalid_post(&bob_cell_id, &handle, dna_file).await;
    let invalid_entry_hash: AnyDgdHash = invalid_entry_hash.into();

    // Integration should have 3 ops in it
    // StoreEntry and StoreElement should be invalid.
    let expected_count = 3 + expected_count;
    let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();
    wait_for_integration(&alice_env, expected_count, num_attempts, delay_per_attempt).await;

    {
        let alice_env = handle.get_cell_env(&alice_cell_id).await.unwrap();

        let workspace = IncomingDgdOpsWorkspace::new(alice_env.clone().into()).unwrap();
        // Validation should be empty
        let val = inspect_val_limbo(&alice_env, &workspace);
        assert_eq!(val.len(), 0);
        let int = inspect_integrated(&alice_env, &workspace);
        for v in &int {
            match &v.1.op {
                // A Store entry that matches these hashes
                DgdOpLight::StoreEntry(hh, _, eh)
                    if *eh == invalid_entry_hash && *hh == invalid_header_hash =>
                {
                    assert_eq!(v.1.validation_status, ValidationStatus::Rejected, "{:?}", v)
                }
                // And the store element
                DgdOpLight::StoreElement(hh, _, _) if *hh == invalid_header_hash => {
                    assert_eq!(v.1.validation_status, ValidationStatus::Rejected, "{:?}", v);
                }
                _ => {}
            }
        }

        assert_eq!(int.len(), expected_count);
    }
}

// Need to "hack aingle" because otherwise the invalid
// commit is caught by the call zome workflow
async fn commit_invalid(
    bob_cell_id: &CellId,
    handle: &ConductorHandle,
    dna_file: &DnaFile,
) -> (HeaderHash, EntryHash) {
    let entry = ThisWasmEntry::NeverValidates;
    let entry_hash = EntryHash::with_data_sync(&Entry::try_from(entry.clone()).unwrap());
    let call_data = HostFnCaller::create(bob_cell_id, handle, dna_file).await;
    // 4
    let invalid_header_hash = call_data
        .commit_entry(entry.clone().try_into().unwrap(), INVALID_ID)
        .await;

    // Produce and publish these commits
    let mut triggers = handle.get_cell_triggers(&bob_cell_id).await.unwrap();
    triggers.produce_dgd_ops.trigger();
    (invalid_header_hash, entry_hash)
}

// Need to "hack aingle" because otherwise the invalid
// commit is caught by the call zome workflow
async fn commit_invalid_post(
    bob_cell_id: &CellId,
    handle: &ConductorHandle,
    dna_file: &DnaFile,
) -> (HeaderHash, EntryHash) {
    // Bananas are not allowed
    let entry = Post("Banana".into());
    let entry_hash = EntryHash::with_data_sync(&Entry::try_from(entry.clone()).unwrap());
    // Create call data for the 3rd zome Create
    let call_data = HostFnCaller::create_for_zome(bob_cell_id, handle, dna_file, 2).await;
    // 9
    let invalid_header_hash = call_data
        .commit_entry(entry.clone().try_into().unwrap(), POST_ID)
        .await;

    // Produce and publish these commits
    let mut triggers = handle.get_cell_triggers(&bob_cell_id).await.unwrap();
    triggers.produce_dgd_ops.trigger();
    (invalid_header_hash, entry_hash)
}

async fn call_zome_directly(
    bob_cell_id: &CellId,
    handle: &ConductorHandle,
    dna_file: &DnaFile,
    invocation: ZomeCallInvocation,
) -> ExternIO {
    let call_data = HostFnCaller::create(bob_cell_id, handle, dna_file).await;
    // 4
    let output = call_data.call_zome_direct(invocation).await;

    // Produce and publish these commits
    let mut triggers = handle.get_cell_triggers(&bob_cell_id).await.unwrap();
    triggers.produce_dgd_ops.trigger();
    output
}

#[instrument(skip(env, workspace))]
fn inspect_val_limbo(
    env: &EnvironmentWrite,
    workspace: &IncomingDgdOpsWorkspace,
) -> Vec<(DgdOpHash, ValidationLimboValue, Option<Element>)> {
    debug!("start");
    let element_buf = ElementBuf::pending(env.clone().into()).unwrap();
    fresh_reader_test!(env, |r| {
        workspace
            .validation_limbo
            .iter(&r)
            .unwrap()
            .map(|(k, i)| {
                let hash = DgdOpHash::from_raw_39_panicky(k.to_vec());
                let el = element_buf.get_element(&i.op.header_hash()).unwrap();
                debug!(?hash, ?i, op_in_val = ?el);
                Ok((hash, i, el))
            })
            .collect()
            .unwrap()
    })
}

#[instrument(skip(env, workspace))]
fn inspect_integrated(
    env: &EnvironmentWrite,
    workspace: &IncomingDgdOpsWorkspace,
) -> Vec<(DgdOpHash, IntegratedDgdOpsValue, Element)> {
    debug!("start");
    let element_buf = ElementBuf::vault(env.clone().into(), true).unwrap();
    let element_buf_reject = ElementBuf::rejected(env.clone().into()).unwrap();
    fresh_reader_test!(env, |r| {
        workspace
            .integrated_dgd_ops
            .iter(&r)
            .unwrap()
            .map(|(k, i)| {
                let hash = DgdOpHash::from_raw_39_panicky(k.to_vec());
                let el = element_buf
                    .get_element(&i.op.header_hash())
                    .unwrap()
                    .or_else(|| element_buf_reject.get_element(&i.op.header_hash()).unwrap())
                    .expect("missing element");
                debug!(?hash, ?i, op_in_int = ?el);
                Ok((hash, i, el))
            })
            .collect()
            .unwrap()
    })
}
