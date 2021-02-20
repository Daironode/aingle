use ::fixt::prelude::*;
use anyhow::Result;
use assert_cmd::prelude::*;
use futures::future;
use futures::Future;
use hdk::prelude::RemoteSignal;
use aingle::test_utils::sweetest::SweetAgents;
use aingle::test_utils::sweetest::SweetConductorBatch;
use aingle::test_utils::sweetest::SweetDnaFile;
use aingle::{
    conductor::api::ZomeCall,
    conductor::{
        api::{AdminRequest, AdminResponse, AppRequest, AppResponse},
        config::*,
        error::ConductorError,
        Conductor,
    },
    fixt::*,
};
use aingle_types::{
    prelude::*,
    test_utils::{fake_agent_pubkey_1, fake_dna_zomes, write_fake_dna_file},
};
use aingle_wasm_test_utils::TestWasm;
use aingle_websocket::*;
use matches::assert_matches;
use observability;
use std::path::PathBuf;
use std::process::Stdio;
use std::sync::Arc;
use std::time::Duration;
use tempdir::TempDir;
use tokio::io::AsyncBufReadExt;
use tokio::io::BufReader;
use tokio::process::Child;
use tokio::process::Command;
use tokio::stream::StreamExt;
use tracing::*;
use url2::prelude::*;

use test_utils::*;

mod test_utils;

pub fn spawn_output(aingle: &mut Child) {
    let stdout = aingle.stdout.take();
    let stderr = aingle.stderr.take();
    tokio::task::spawn(async move {
        if let Some(stdout) = stdout {
            let mut reader = BufReader::new(stdout).lines();
            while let Ok(Some(line)) = reader.next_line().await {
                trace!("aingle bin stdout: {}", line);
            }
        }
    });
    tokio::task::spawn(async move {
        if let Some(stderr) = stderr {
            let mut reader = BufReader::new(stderr).lines();
            while let Ok(Some(line)) = reader.next_line().await {
                trace!("aingle bin stderr: {}", line);
            }
        }
    });
}

pub async fn check_started(aingle: &mut Child) {
    let started = tokio::time::timeout(std::time::Duration::from_secs(1), aingle).await;
    if let Ok(status) = started {
        panic!("AIngle failed to start. status: {:?}", status);
    }
}

fn create_config(port: u16, environment_path: PathBuf) -> ConductorConfig {
    ConductorConfig {
        admin_interfaces: Some(vec![AdminInterfaceConfig {
            driver: InterfaceDriver::Websocket { port },
        }]),
        environment_path: environment_path.into(),
        network: None,
        signing_service_uri: None,
        encryption_service_uri: None,
        decryption_service_uri: None,
        dpki: None,
        passphrase_service: Some(PassphraseServiceConfig::FromConfig {
            passphrase: "password".into(),
        }),
        keystore_path: None,
        use_dangerous_test_keystore: true,
    }
}

pub fn write_config(mut path: PathBuf, config: &ConductorConfig) -> PathBuf {
    path.push("conductor_config.yml");
    std::fs::write(path.clone(), serde_yaml::to_string(&config).unwrap()).unwrap();
    path
}

#[instrument(skip(aingle, response))]
async fn check_timeout<T>(
    aingle: &mut Child,
    response: impl Future<Output = Result<T, std::io::Error>>,
    timeout_millis: u64,
) -> T {
    match tokio::time::timeout(std::time::Duration::from_millis(timeout_millis), response).await {
        Ok(response) => response.unwrap(),
        Err(_) => {
            aingle.kill().unwrap();
            error!("Timeout");
            panic!("Timed out on request after {}", timeout_millis);
        }
    }
}

#[tokio::test(threaded_scheduler)]
#[cfg(feature = "slow_tests")]
async fn call_admin() {
    observability::test_run().ok();
    // NOTE: This is a full integration test that
    // actually runs the aingle binary

    // TODO: B-01453: can we make this port 0 and find out the dynamic port later?
    let port = 9909;

    let tmp_dir = TempDir::new("conductor_cfg").unwrap();
    let path = tmp_dir.path().to_path_buf();
    let environment_path = path.clone();
    let config = create_config(port, environment_path);
    let config_path = write_config(path, &config);

    let uuid = uuid::Uuid::new_v4();
    let dna = fake_dna_zomes(
        &uuid.to_string(),
        vec![(TestWasm::Foo.into(), TestWasm::Foo.into())],
    );

    let cmd = std::process::Command::cargo_bin("aingle").unwrap();
    let mut cmd = Command::from(cmd);
    cmd.arg("--structured")
        .arg("--config-path")
        .arg(config_path)
        .env("RUST_LOG", "debug")
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .kill_on_drop(true);
    let mut aingle = cmd.spawn().expect("Failed to spawn aingle");
    spawn_output(&mut aingle);
    check_started(&mut aingle).await;

    let (mut client, _) = websocket_client_by_port(port).await.unwrap();

    let original_dna_hash = dna.dna_hash().clone();

    // Make properties
    let properties: aingle_types::dna::JsonProperties = serde_json::json!({
        "test": "example",
        "how_many": 42,
    })
    .into();

    // Install Dna
    let (fake_dna_path, _tmpdir) = write_fake_dna_file(dna.clone()).await.unwrap();
    let dna_payload = InstallAppDnaPayload {
        path: Some(fake_dna_path),
        hash: None,
        nick: "nick".into(),
        properties: Some(properties.clone()),
        membrane_proof: None,
    };
    let agent_key = fake_agent_pubkey_1();
    let payload = InstallAppPayload {
        dnas: vec![dna_payload],
        installed_app_id: "test".to_string(),
        agent_key,
    };
    let request = AdminRequest::InstallApp(Box::new(payload));
    let response = client.request(request);
    let response = check_timeout(&mut aingle, response, 3000).await;
    assert_matches!(response, AdminResponse::AppInstalled(_));

    // List Dnas
    let request = AdminRequest::ListDnas;
    let response = client.request(request);
    let response = check_timeout(&mut aingle, response, 3000).await;

    let tmp_wasm = dna.code().values().cloned().collect::<Vec<_>>();
    let mut tmp_dna = dna.dna_def().clone();
    tmp_dna.properties = properties.try_into().unwrap();
    let dna = aingle_types::dna::DnaFile::new(tmp_dna, tmp_wasm)
        .await
        .unwrap();

    assert_ne!(&original_dna_hash, dna.dna_hash());

    let expects = vec![dna.dna_hash().clone()];
    assert_matches!(response, AdminResponse::DnasListed(a) if a == expects);

    aingle.kill().expect("Failed to kill aingle");
}

pub async fn start_aingle(config_path: PathBuf) -> Child {
    tracing::info!("\n\n----\nstarting aingle\n----\n\n");
    let cmd = std::process::Command::cargo_bin("aingle").unwrap();
    let mut cmd = Command::from(cmd);
    cmd.arg("--structured")
        .arg("--config-path")
        .arg(config_path)
        .env("RUST_LOG", "trace")
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .kill_on_drop(true);
    let mut aingle = cmd.spawn().expect("Failed to spawn aingle");
    spawn_output(&mut aingle);
    check_started(&mut aingle).await;
    aingle
}

pub async fn call_foo_fn(app_port: u16, original_dna_hash: DnaHash, aingle: &mut Child) {
    // Connect to App Interface
    let (mut app_tx, _) = websocket_client_by_port(app_port).await.unwrap();
    let cell_id = CellId::from((original_dna_hash, fake_agent_pubkey_1()));
    call_zome_fn(
        aingle,
        &mut app_tx,
        cell_id,
        TestWasm::Foo,
        "foo".into(),
        (),
    )
    .await;
    app_tx.close(1000, "Shutting down".into()).await.unwrap();
}

pub async fn call_zome_fn<S>(
    aingle: &mut Child,
    app_tx: &mut WebsocketSender,
    cell_id: CellId,
    wasm: TestWasm,
    fn_name: String,
    input: S,
) where
    S: Serialize + std::fmt::Debug,
{
    let call: ZomeCall = ZomeCallInvocationFixturator::new(NamedInvocation(
        cell_id,
        wasm,
        fn_name,
        ExternIO::encode(input).unwrap(),
    ))
    .next()
    .unwrap()
    .into();
    let request = AppRequest::ZomeCallInvocation(Box::new(call));
    let response = app_tx.request(request);
    let call_response = check_timeout(aingle, response, 3000).await;
    trace!(?call_response);
    assert_matches!(call_response, AppResponse::ZomeCallInvocation(_));
}

pub async fn attach_app_interface(
    client: &mut WebsocketSender,
    aingle: &mut Child,
    port: Option<u16>,
) -> u16 {
    let request = AdminRequest::AttachAppInterface { port };
    let response = client.request(request);
    let response = check_timeout(aingle, response, 1000).await;
    match response {
        AdminResponse::AppInterfaceAttached { port } => port,
        _ => panic!("Attach app interface failed: {:?}", response),
    }
}

pub async fn retry_admin_interface(
    port: u16,
    mut attempts: usize,
    delay: Duration,
) -> WebsocketSender {
    loop {
        match websocket_client_by_port(port).await {
            Ok(c) => return c.0,
            Err(e) => {
                attempts -= 1;
                if attempts == 0 {
                    panic!("Failed to join admin interface");
                }
                warn!(
                    "Failed with {:?} to open admin interface, trying {} more times",
                    e, attempts
                );
                tokio::time::delay_for(delay).await;
            }
        }
    }
}

#[tokio::test(threaded_scheduler)]
#[cfg(feature = "slow_tests")]
async fn call_zome() {
    observability::test_run().ok();
    // NOTE: This is a full integration test that
    // actually runs the aingle binary

    // TODO: B-01453: can we make this port 0 and find out the dynamic port later?
    let admin_port = 9910;
    let app_port = 9913;

    let tmp_dir = TempDir::new("conductor_cfg_2").unwrap();
    let path = tmp_dir.path().to_path_buf();
    let environment_path = path.clone();
    let config = create_config(admin_port, environment_path);
    let config_path = write_config(path, &config);

    let mut aingle = start_aingle(config_path.clone()).await;

    let (mut client, _) = websocket_client_by_port(admin_port).await.unwrap();
    let (_, receiver2) = websocket_client_by_port(admin_port).await.unwrap();

    let uuid = uuid::Uuid::new_v4();
    let dna = fake_dna_zomes(
        &uuid.to_string(),
        vec![(TestWasm::Foo.into(), TestWasm::Foo.into())],
    );
    let original_dna_hash = dna.dna_hash().clone();

    // Install Dna
    let (fake_dna_path, _tmpdir) = write_fake_dna_file(dna.clone()).await.unwrap();
    let dna_payload = InstallAppDnaPayload::path_only(fake_dna_path, "".to_string());
    let agent_key = fake_agent_pubkey_1();
    let payload = InstallAppPayload {
        dnas: vec![dna_payload],
        installed_app_id: "test".to_string(),
        agent_key,
    };
    let request = AdminRequest::InstallApp(Box::new(payload));
    let response = client.request(request);
    let response = check_timeout(&mut aingle, response, 3000).await;
    assert_matches!(response, AdminResponse::AppInstalled(_));

    // List Dnas
    let request = AdminRequest::ListDnas;
    let response = client.request(request);
    let response = check_timeout(&mut aingle, response, 1000).await;

    let expects = vec![original_dna_hash.clone()];
    assert_matches!(response, AdminResponse::DnasListed(a) if a == expects);

    // Activate cells
    let request = AdminRequest::ActivateApp {
        installed_app_id: "test".to_string(),
    };
    let response = client.request(request);
    let response = check_timeout(&mut aingle, response, 1000).await;
    assert_matches!(response, AdminResponse::AppActivated);

    // Attach App Interface
    let app_port_rcvd = attach_app_interface(&mut client, &mut aingle, Some(app_port)).await;
    assert_eq!(app_port, app_port_rcvd);

    // Call Zome
    call_foo_fn(app_port, original_dna_hash.clone(), &mut aingle).await;

    // Ensure that the other client does not receive any messages, i.e. that
    // responses are not broadcast to all connected clients, only the one
    // that made the request.
    assert!(
        receiver2
            .timeout(Duration::from_millis(500))
            .next()
            .await
            .unwrap()
            .is_err() // Err means the timeout elapsed
    );

    client.close(1000, "Shutting down".into()).await.unwrap();
    // Shutdown aingle
    aingle.kill().expect("Failed to kill aingle");
    std::mem::drop(client);

    // Call zome after resart
    let mut aingle = start_aingle(config_path).await;

    tokio::time::delay_for(std::time::Duration::from_millis(1000)).await;

    // Call Zome again on the existing app interface port
    call_foo_fn(app_port, original_dna_hash, &mut aingle).await;

    // Shutdown aingle
    aingle.kill().expect("Failed to kill aingle");
}

#[tokio::test(threaded_scheduler)]
#[cfg(feature = "slow_tests")]
async fn remote_signals() -> anyhow::Result<()> {
    observability::test_run().ok();
    const NUM_CONDUCTORS: usize = 5;

    let mut conductors = SweetConductorBatch::from_standard_config(NUM_CONDUCTORS).await;

    // TODO: write helper for agents across conductors
    let all_agents: Vec<AIngleHash<hash_type::Agent>> =
        future::join_all(conductors.iter().map(|c| SweetAgents::one(c.keystore()))).await;

    let dna_file = SweetDnaFile::unique_from_test_wasms(vec![TestWasm::EmitSignal])
        .await
        .unwrap()
        .0;

    let apps = conductors
        .setup_app_for_zipped_agents("app", &all_agents, &[dna_file])
        .await;

    conductors.exchange_peer_info().await;

    let cells = apps.cells_flattened();

    let mut rxs = Vec::new();
    for h in conductors.iter().map(|c| c) {
        rxs.push(h.signal_broadcaster().await.subscribe_separately())
    }
    let rxs = rxs.into_iter().flatten().collect::<Vec<_>>();

    let signal = fixt!(ExternIo);

    let _: () = conductors[0]
        .call(
            &cells[0].zome(TestWasm::EmitSignal),
            "signal_others",
            RemoteSignal {
                signal: signal.clone(),
                agents: all_agents,
            },
        )
        .await;

    tokio::time::delay_for(std::time::Duration::from_millis(2000)).await;

    let signal = AppSignal::new(signal);
    for mut rx in rxs {
        let r = rx.try_recv();
        // Each handle should recv a signal
        assert_matches!(r, Ok(Signal::App(_, a)) if a == signal);
    }

    Ok(())
}

#[tokio::test(threaded_scheduler)]
#[cfg(feature = "slow_tests")]
async fn emit_signals() {
    observability::test_run().ok();
    // NOTE: This is a full integration test that
    // actually runs the aingle binary

    // TODO: B-01453: can we make this port 0 and find out the dynamic port later?
    let port = 9911;

    let tmp_dir = TempDir::new("conductor_cfg_3").unwrap();
    let path = tmp_dir.path().to_path_buf();
    let environment_path = path.clone();
    let config = create_config(port, environment_path);
    let config_path = write_config(path, &config);

    let mut aingle = start_aingle(config_path.clone()).await;

    let (mut admin_tx, _) = websocket_client_by_port(port).await.unwrap();

    let uuid = uuid::Uuid::new_v4();
    let dna = fake_dna_zomes(
        &uuid.to_string(),
        vec![(TestWasm::EmitSignal.into(), TestWasm::EmitSignal.into())],
    );
    let dna_hash = dna.dna_hash().clone();

    // Install Dna
    let (fake_dna_path, _tmpdir) = write_fake_dna_file(dna.clone()).await.unwrap();
    let dna_payload = InstallAppDnaPayload::path_only(fake_dna_path, "".to_string());
    let agent_key = fake_agent_pubkey_1();
    let cell_id = CellId::new(dna_hash.clone(), agent_key.clone());
    let payload = InstallAppPayload {
        dnas: vec![dna_payload],
        installed_app_id: "test".to_string(),
        agent_key: agent_key.clone(),
    };
    let request = AdminRequest::InstallApp(Box::new(payload));
    let response = admin_tx.request(request);
    let response = check_timeout(&mut aingle, response, 3000).await;
    assert_matches!(response, AdminResponse::AppInstalled(_));

    // Activate cells
    let request = AdminRequest::ActivateApp {
        installed_app_id: "test".to_string(),
    };
    let response = admin_tx.request(request);
    let response = check_timeout(&mut aingle, response, 1000).await;
    assert_matches!(response, AdminResponse::AppActivated);

    // Attach App Interface
    let app_port = attach_app_interface(&mut admin_tx, &mut aingle, None).await;

    ///////////////////////////////////////////////////////
    // Emit signals (the real test!)

    let (mut app_tx_1, app_rx_1) = websocket_client_by_port(app_port).await.unwrap();
    let (_, app_rx_2) = websocket_client_by_port(app_port).await.unwrap();

    call_zome_fn(
        &mut aingle,
        &mut app_tx_1,
        cell_id.clone(),
        TestWasm::EmitSignal,
        "emit".into(),
        (),
    )
    .await;

    let msg1 = app_rx_1
        .timeout(Duration::from_secs(1))
        .next()
        .await
        .unwrap()
        .unwrap();
    let sig1: SerializedBytes = unwrap_to::unwrap_to!(msg1 => WebsocketMessage::Signal).clone();

    let msg2 = app_rx_2
        .timeout(Duration::from_secs(1))
        .next()
        .await
        .unwrap()
        .unwrap();
    let sig2: SerializedBytes = unwrap_to::unwrap_to!(msg2 => WebsocketMessage::Signal).clone();

    assert_eq!(
        Signal::App(cell_id, AppSignal::new(ExternIO::encode(()).unwrap())),
        Signal::try_from(sig1.clone()).unwrap(),
    );
    assert_eq!(sig1, sig2);

    ///////////////////////////////////////////////////////

    admin_tx.close(1000, "Shutting down".into()).await.unwrap();
    // Shutdown aingle
    aingle.kill().expect("Failed to kill aingle");
}

#[tokio::test(threaded_scheduler)]
async fn conductor_admin_interface_runs_from_config() -> Result<()> {
    observability::test_run().ok();
    let tmp_dir = TempDir::new("conductor_cfg").unwrap();
    let environment_path = tmp_dir.path().to_path_buf();
    let config = create_config(0, environment_path);
    let conductor_handle = Conductor::builder().config(config).build().await?;
    let (mut client, _) = websocket_client(&conductor_handle).await?;

    let dna = fake_dna_zomes(
        "".into(),
        vec![(TestWasm::Foo.into(), TestWasm::Foo.into())],
    );
    let (fake_dna_path, _tmpdir) = write_fake_dna_file(dna).await.unwrap();
    let dna_payload = InstallAppDnaPayload::path_only(fake_dna_path, "".to_string());
    let agent_key = fake_agent_pubkey_1();
    let payload = InstallAppPayload {
        dnas: vec![dna_payload],
        installed_app_id: "test".to_string(),
        agent_key,
    };
    let request = AdminRequest::InstallApp(Box::new(payload));
    let response = client.request(request).await;
    assert_matches!(response, Ok(AdminResponse::AppInstalled(_)));
    conductor_handle.shutdown().await;

    Ok(())
}

#[tokio::test(threaded_scheduler)]
async fn conductor_admin_interface_ends_with_shutdown() -> Result<()> {
    if let Err(e) = conductor_admin_interface_ends_with_shutdown_inner().await {
        panic!("{:#?}", e);
    }
    Ok(())
}

async fn conductor_admin_interface_ends_with_shutdown_inner() -> Result<()> {
    observability::test_run().ok();

    info!("creating config");
    let tmp_dir = TempDir::new("conductor_cfg").unwrap();
    let environment_path = tmp_dir.path().to_path_buf();
    let config = create_config(0, environment_path);
    let conductor_handle = Conductor::builder().config(config).build().await?;
    let port = admin_port(&conductor_handle).await;
    info!("building conductor");
    let (mut client, rx): (WebsocketSender, WebsocketReceiver) = websocket_connect(
        url2!("ws://127.0.0.1:{}", port),
        Arc::new(WebsocketConfig {
            default_request_timeout_s: 1,
            ..Default::default()
        }),
    )
    .await?;

    info!("client connect");

    conductor_handle.shutdown().await;

    info!("shutdown");

    assert_matches!(
        conductor_handle.check_running().await,
        Err(ConductorError::ShuttingDown)
    );

    let incoming: Vec<_> = rx.collect().await;
    assert_eq!(incoming.len(), 1);
    assert_matches!(incoming[0], WebsocketMessage::Close(_));

    info!("About to make failing request");

    let dna = fake_dna_zomes(
        "".into(),
        vec![(TestWasm::Foo.into(), TestWasm::Foo.into())],
    );
    let (fake_dna_path, _tmpdir) = write_fake_dna_file(dna).await.unwrap();
    let dna_payload = InstallAppDnaPayload::path_only(fake_dna_path, "".to_string());
    let agent_key = fake_agent_pubkey_1();
    let payload = InstallAppPayload {
        dnas: vec![dna_payload],
        installed_app_id: "test".to_string(),
        agent_key,
    };
    let request = AdminRequest::InstallApp(Box::new(payload));

    // send a request after the conductor has shutdown
    let response: Result<Result<AdminResponse, _>, tokio::time::Elapsed> =
        tokio::time::timeout(Duration::from_secs(1), client.request(request)).await;

    // request should have encountered an error since the conductor shut down,
    // but should not have timed out (which would be an `Err(Err(_))`)
    assert_matches!(response, Ok(Err(_)));

    Ok(())
}

#[tokio::test(threaded_scheduler)]
async fn too_many_open() {
    observability::test_run().ok();

    info!("creating config");
    let tmp_dir = TempDir::new("conductor_cfg").unwrap();
    let environment_path = tmp_dir.path().to_path_buf();
    let config = create_config(0, environment_path);
    let conductor_handle = Conductor::builder().config(config).build().await.unwrap();
    let port = admin_port(&conductor_handle).await;
    info!("building conductor");
    for i in 0..1000 {
        dbg!(i);
        let (_client, _rx): (WebsocketSender, WebsocketReceiver) = websocket_connect(
            url2!("ws://127.0.0.1:{}", port),
            Arc::new(WebsocketConfig {
                default_request_timeout_s: 1,
                ..Default::default()
            }),
        )
        .await
        .unwrap();
    }
    conductor_handle.shutdown().await;
}