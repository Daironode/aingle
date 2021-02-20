<<<<<<< HEAD
//! The workflow and queue consumer for DgdOp integration
=======
//! The workflow and queue consumer for DhtOp integration
>>>>>>> master

use super::*;

use crate::conductor::manager::ManagedTaskResult;
<<<<<<< HEAD
use crate::core::workflow::integrate_dgd_ops_workflow::integrate_dgd_ops_workflow;
use crate::core::workflow::integrate_dgd_ops_workflow::IntegrateDgdOpsWorkspace;
=======
use crate::core::workflow::integrate_dht_ops_workflow::integrate_dht_ops_workflow;
use crate::core::workflow::integrate_dht_ops_workflow::IntegrateDhtOpsWorkspace;
>>>>>>> master
use aingle_lmdb::env::EnvironmentWrite;

use tokio::task::JoinHandle;
use tracing::*;

<<<<<<< HEAD
/// Spawn the QueueConsumer for DgdOpIntegration workflow
#[instrument(skip(env, stop, trigger_sys))]
pub fn spawn_integrate_dgd_ops_consumer(
=======
/// Spawn the QueueConsumer for DhtOpIntegration workflow
#[instrument(skip(env, stop, trigger_sys))]
pub fn spawn_integrate_dht_ops_consumer(
>>>>>>> master
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
<<<<<<< HEAD
                    "Cell is shutting down: stopping integrate_dgd_ops_workflow queue consumer."
=======
                    "Cell is shutting down: stopping integrate_dht_ops_workflow queue consumer."
>>>>>>> master
                );
                break;
            }

            // Run the workflow
<<<<<<< HEAD
            let workspace = IntegrateDgdOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                integrate_dgd_ops_workflow(workspace, env.clone().into(), &mut trigger_sys)
=======
            let workspace = IntegrateDhtOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                integrate_dht_ops_workflow(workspace, env.clone().into(), &mut trigger_sys)
>>>>>>> master
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
