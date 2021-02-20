use super::*;
use ::fixt::prelude::*;
use aingle_keystore::AgentPubKeyExt;

#[tokio::test(threaded_scheduler)]
async fn incoming_ops_to_limbo() {
    let test_env = aingle_lmdb::test_utils::test_cell_env();
    let env = test_env.env();
    let keystore = aingle_lmdb::test_utils::test_keystore();
    let (sys_validation_trigger, mut rx) = TriggerSender::new();

    let author = fake_agent_pubkey_1();
    let mut header = fixt!(CreateLink);
    header.author = author.clone();
    let header = Header::CreateLink(header);
    let signature = author.sign(&keystore, &header).await.unwrap();

    let op = DgdOp::RegisterAgentActivity(signature, header);
    let op_light = op.to_light();
    let hash = DgdOpHash::with_data_sync(&op);
    let ops = vec![(hash.clone(), op.clone())];

    incoming_dgd_ops_workflow(&env, sys_validation_trigger.clone(), ops, None)
        .await
        .unwrap();
    rx.listen().await.unwrap();

    let workspace = IncomingDgdOpsWorkspace::new(env.clone().into()).unwrap();
    let r = workspace.validation_limbo.get(&hash).unwrap().unwrap();
    assert_eq!(r.op, op_light);
}
