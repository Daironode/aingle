//! The workflow and queue consumer for DgdOp production

use super::*;
use crate::conductor::manager::ManagedTaskResult;
use crate::core::workflow::produce_dgd_ops_workflow::produce_dgd_ops_workflow;
use crate::core::workflow::produce_dgd_ops_workflow::ProduceDgdOpsWorkspace;
use aingle_lmdb::env::EnvironmentWrite;

use tokio::task::JoinHandle;
use tracing::*;

/// Spawn the QueueConsumer for Produce_dgd_ops workflow
#[instrument(skip(env, stop, trigger_publish))]
pub fn spawn_produce_dgd_ops_consumer(
    env: EnvironmentWrite,
    mut stop: sync::broadcast::Receiver<()>,
    mut trigger_publish: TriggerSender,
) -> (TriggerSender, JoinHandle<ManagedTaskResult>) {
    let (tx, mut rx) = TriggerSender::new();
    let mut trigger_self = tx.clone();
    let handle = tokio::spawn(async move {
        loop {
            if let Job::Shutdown = next_job_or_exit(&mut rx, &mut stop).await {
                tracing::warn!(
                    "Cell is shutting down: stopping produce_dgd_ops_workflow queue consumer."
                );
                break;
            }

            let workspace = ProduceDgdOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                produce_dgd_ops_workflow(workspace, env.clone().into(), &mut trigger_publish)
                    .await
                    .expect("Error running Workflow")
            {
                trigger_self.trigger()
            };
        }
        Ok(())
    });
    (tx, handle)
}
