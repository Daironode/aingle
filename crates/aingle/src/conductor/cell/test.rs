use crate::conductor::manager::spawn_task_manager;
use crate::core::ribosome::guest_callback::genesis_self_check::GenesisSelfCheckResult;
use crate::core::ribosome::MockRibosomeT;
use crate::core::workflow::incoming_sgd_ops_workflow::op_exists;
use crate::fixt::SafFileFixturator;
use crate::fixt::SignatureFixturator;
use crate::test_utils::test_network;
use ::fixt::prelude::*;
use ai_hash::HasHash;
use aingle_state::prelude::*;
use aingle_types::prelude::*;
use aingle_zome_types::header;
use aingle_zome_types::HeaderHashed;
use std::sync::Arc;
use tokio::sync;

#[tokio::test(flavor = "multi_thread")]
async fn test_cell_handle_publish() {
    let cell_env = test_cell_env();
    let cache_env = test_cache_env();
    let env = cell_env.env();
    let cache = cache_env.env();

    let cell_id = fake_cell_id(1);
    let saf = cell_id.saf_hash().clone();
    let agent = cell_id.agent_pubkey().clone();

    let test_network = test_network(Some(saf.clone()), Some(agent.clone())).await;
    let aingle_p2p_cell = test_network.cell_network();

    let mut mock_handle = crate::conductor::handle::MockConductorHandleT::new();
    mock_handle
        .expect_get_saf()
        .returning(|_| Some(fixt!(SafFile)));

    let mock_handle: crate::conductor::handle::ConductorHandle = Arc::new(mock_handle);
    let mut mock_ribosome = MockRibosomeT::new();
    mock_ribosome
        .expect_run_genesis_self_check()
        .returning(|_, _| Ok(GenesisSelfCheckResult::Valid));

    super::Cell::genesis(
        cell_id.clone(),
        mock_handle.clone(),
        env.clone(),
        mock_ribosome,
        None,
    )
    .await
    .unwrap();

    let (add_task_sender, shutdown) = spawn_task_manager(mock_handle.clone());
    let (stop_tx, _) = sync::broadcast::channel(1);

    let (cell, _) = super::Cell::create(
        cell_id,
        mock_handle,
        env.clone(),
        cache.clone(),
        aingle_p2p_cell,
        add_task_sender,
        stop_tx.clone(),
    )
    .await
    .unwrap();

    let sig = fixt!(Signature);
    let header = header::Header::Saf(header::Saf {
        author: agent.clone(),
        timestamp: timestamp::now().into(),
        hash: saf.clone(),
    });
    let op = SgdOp::StoreElement(sig, header.clone(), None);
    let op_hash = SgdOpHashed::from_content_sync(op.clone()).into_hash();
    let header_hash = HeaderHashed::from_content_sync(header.clone()).into_hash();

    cell.handle_publish(
        fake_agent_pubkey_2(),
        true,
        header_hash.clone().into(),
        vec![(op_hash.clone(), op.clone())],
    )
    .await
    .unwrap();

    op_exists(&cell.env, &op_hash).unwrap();

    stop_tx.send(()).unwrap();
    shutdown.await.unwrap().unwrap();
}
