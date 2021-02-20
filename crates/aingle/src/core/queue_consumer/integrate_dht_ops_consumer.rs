//! The workflow and queue consumer for DgdOp integration

use super::*;

use crate::conductor::manager::ManagedTaskResult;
use crate::core::workflow::integrate_dgd_ops_workflow::integrate_dgd_ops_workflow;
use crate::core::workflow::integrate_dgd_ops_workflow::IntegrateDgdOpsWorkspace;
use aingle_lmdb::env::EnvironmentWrite;

use tokio::task::JoinHandle;
use tracing::*;

/// Spawn the QueueConsumer for DgdOpIntegration workflow
#[instrument(skip(env, stop, trigger_sys))]
pub fn spawn_integrate_dgd_ops_consumer(
    env: EnvironmentWrite,
    mut stop: sync::broadcast::Receiver<()>,
    trigger_sys: sync::oneshot::Receiver<TriggerSender>,
) -> (TriggerSender, JoinHandle<ManagedTaskResult>) {
    let (tx, mut rx) = TriggerSender::new();
    let mut trigger_self = tx.clone();
    let handle = tokio::spawn(async move {
        let mut trigger_sys = trigger_sys.await.expect("failed to get tx sys");
        loop {
            // Wait for next job
            if let Job::Shutdown = next_job_or_exit(&mut rx, &mut stop).await {
                tracing::warn!(
                    "Cell is shutting down: stopping integrate_dgd_ops_workflow queue consumer."
                );
                break;
            }

            // Run the workflow
            let workspace = IntegrateDgdOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                integrate_dgd_ops_workflow(workspace, env.clone().into(), &mut trigger_sys)
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
