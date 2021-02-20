//! The workflow and queue consumer for sys validation

use super::*;

use crate::conductor::manager::ManagedTaskResult;
<<<<<<< HEAD
use crate::core::workflow::publish_dgd_ops_workflow::publish_dgd_ops_workflow;
use crate::core::workflow::publish_dgd_ops_workflow::PublishDgdOpsWorkspace;
=======
use crate::core::workflow::publish_dht_ops_workflow::publish_dht_ops_workflow;
use crate::core::workflow::publish_dht_ops_workflow::PublishDhtOpsWorkspace;
>>>>>>> master
use aingle_lmdb::env::EnvironmentWrite;

use tokio::task::JoinHandle;
use tracing::*;

/// Spawn the QueueConsumer for Publish workflow
#[instrument(skip(env, stop, cell_network))]
<<<<<<< HEAD
pub fn spawn_publish_dgd_ops_consumer(
=======
pub fn spawn_publish_dht_ops_consumer(
>>>>>>> master
    env: EnvironmentWrite,
    mut stop: sync::broadcast::Receiver<()>,
    mut cell_network: AIngleP2pCell,
) -> (TriggerSender, JoinHandle<ManagedTaskResult>) {
    let (tx, mut rx) = TriggerSender::new();
    let mut trigger_self = tx.clone();
    let handle = tokio::spawn(async move {
        loop {
            // Wait for next job
            if let Job::Shutdown = next_job_or_exit(&mut rx, &mut stop).await {
                tracing::warn!(
<<<<<<< HEAD
                    "Cell is shutting down: stopping publish_dgd_ops_workflow queue consumer."
=======
                    "Cell is shutting down: stopping publish_dht_ops_workflow queue consumer."
>>>>>>> master
                );
                break;
            }

            // Run the workflow
<<<<<<< HEAD
            let workspace = PublishDgdOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                publish_dgd_ops_workflow(workspace, env.clone().into(), &mut cell_network)
=======
            let workspace = PublishDhtOpsWorkspace::new(env.clone().into())
                .expect("Could not create Workspace");
            if let WorkComplete::Incomplete =
                publish_dht_ops_workflow(workspace, env.clone().into(), &mut cell_network)
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
