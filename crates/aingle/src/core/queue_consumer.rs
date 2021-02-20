//! Manages the spawning of tasks which process the various work queues in
//! the system, as well as notifying subsequent queue processors to pick up the
//! work that was left off.
//!
//! The following table lays out the queues and the workflows that consume them,
//! as well as the follow-up workflows. A "source" queue is a database which
//! feeds data to the workflow, and a "destination" queue is a database which
//! said workflow writes to as part of its processing of its source queue.
//!
//! | workflow       | source queue     | dest. queue      | notifies       |
//! |----------------|------------------|------------------|----------------|
//! |                        **gossip path**                                |
//! | HandleGossip   | *n/a*            | ValidationQueue  | SysValidation  |
//! | SysValidation  | ValidationQueue  | ValidationQueue  | AppValidation  |
<<<<<<< HEAD
//! | AppValidation  | ValidationQueue  | ValidationQueue  | DgdOpIntegr.   |
//! |                       **authoring path**                              |
//! | CallZome       | *n/a*            | ChainSequence    | ProduceDgdOps  |
//! | ProduceDgdOps  | ChainSequence    | Auth'd + IntQ †  | DgdOpIntegr.   |
//! |                 **integration, common to both paths**                 |
//! | DgdOpIntegr.   | IntegrationLimbo | IntegratedDgdOps | Publish        |
//! | Publish        | AuthoredDgdOps   | *n/a*            | *n/a*          |
//!
//! († Auth'd + IntQ is short for: AuthoredDgdOps + IntegrationLimbo)
=======
//! | AppValidation  | ValidationQueue  | ValidationQueue  | DhtOpIntegr.   |
//! |                       **authoring path**                              |
//! | CallZome       | *n/a*            | ChainSequence    | ProduceDhtOps  |
//! | ProduceDhtOps  | ChainSequence    | Auth'd + IntQ †  | DhtOpIntegr.   |
//! |                 **integration, common to both paths**                 |
//! | DhtOpIntegr.   | IntegrationLimbo | IntegratedDhtOps | Publish        |
//! | Publish        | AuthoredDhtOps   | *n/a*            | *n/a*          |
//!
//! († Auth'd + IntQ is short for: AuthoredDhtOps + IntegrationLimbo)
>>>>>>> master
//!
//! Implicitly, every workflow also writes to its own source queue, i.e. to
//! remove the item it has just processed.

use derive_more::Constructor;
use derive_more::Display;
use derive_more::From;
use futures::future::Either;
use aingle_lmdb::env::EnvironmentWrite;
use aingle_lmdb::env::WriteManager;
use aingle_lmdb::prelude::Writer;
use tokio::sync;
use tokio::sync::mpsc;

// TODO: move these to workflow mod
<<<<<<< HEAD
mod integrate_dgd_ops_consumer;
use integrate_dgd_ops_consumer::*;
=======
mod integrate_dht_ops_consumer;
use integrate_dht_ops_consumer::*;
>>>>>>> master
mod sys_validation_consumer;
use sys_validation_consumer::*;
mod app_validation_consumer;
use app_validation_consumer::*;
<<<<<<< HEAD
mod produce_dgd_ops_consumer;
use produce_dgd_ops_consumer::*;
mod publish_dgd_ops_consumer;
=======
mod produce_dht_ops_consumer;
use produce_dht_ops_consumer::*;
mod publish_dht_ops_consumer;
>>>>>>> master
use crate::conductor::api::CellConductorApiT;
use crate::conductor::manager::ManagedTaskAdd;
use aingle_p2p::AIngleP2pCell;
use aingle_state::workspace::WorkspaceError;
<<<<<<< HEAD
use publish_dgd_ops_consumer::*;
=======
use publish_dht_ops_consumer::*;
>>>>>>> master

/// Spawns several long-running tasks which are responsible for processing work
/// which shows up on various databases.
///
/// Waits for the initial loop to complete before returning, to prevent causing
/// a race condition by trying to run a workflow too soon after cell creation.
pub async fn spawn_queue_consumer_tasks(
    env: &EnvironmentWrite,
    cell_network: AIngleP2pCell,
    conductor_api: impl CellConductorApiT + 'static,
    mut task_sender: sync::mpsc::Sender<ManagedTaskAdd>,
    stop: sync::broadcast::Sender<()>,
) -> (QueueTriggers, InitialQueueTriggers) {
    // Publish
    let (tx_publish, handle) =
<<<<<<< HEAD
        spawn_publish_dgd_ops_consumer(env.clone(), stop.subscribe(), cell_network.clone());
=======
        spawn_publish_dht_ops_consumer(env.clone(), stop.subscribe(), cell_network.clone());
>>>>>>> master
    task_sender
        .send(ManagedTaskAdd::dont_handle(handle))
        .await
        .expect("Failed to manage workflow handle");

    let (create_tx_sys, get_tx_sys) = tokio::sync::oneshot::channel();

    // Integration
    let (tx_integration, handle) =
<<<<<<< HEAD
        spawn_integrate_dgd_ops_consumer(env.clone(), stop.subscribe(), get_tx_sys);
=======
        spawn_integrate_dht_ops_consumer(env.clone(), stop.subscribe(), get_tx_sys);
>>>>>>> master
    task_sender
        .send(ManagedTaskAdd::dont_handle(handle))
        .await
        .expect("Failed to manage workflow handle");

    // App validation
    let (tx_app, handle) = spawn_app_validation_consumer(
        env.clone(),
        stop.subscribe(),
        tx_integration.clone(),
        conductor_api.clone(),
        cell_network.clone(),
    );
    task_sender
        .send(ManagedTaskAdd::dont_handle(handle))
        .await
        .expect("Failed to manage workflow handle");

    // Sys validation
    let (tx_sys, handle) = spawn_sys_validation_consumer(
        env.clone(),
        stop.subscribe(),
        tx_app.clone(),
        cell_network,
        conductor_api,
    );
    task_sender
        .send(ManagedTaskAdd::dont_handle(handle))
        .await
        .expect("Failed to manage workflow handle");
    if create_tx_sys.send(tx_sys.clone()).is_err() {
        panic!("Failed to send tx_sys");
    }

    // Produce
    let (tx_produce, handle) =
<<<<<<< HEAD
        spawn_produce_dgd_ops_consumer(env.clone(), stop.subscribe(), tx_publish.clone());
=======
        spawn_produce_dht_ops_consumer(env.clone(), stop.subscribe(), tx_publish.clone());
>>>>>>> master
    task_sender
        .send(ManagedTaskAdd::dont_handle(handle))
        .await
        .expect("Failed to manage workflow handle");

    (
        QueueTriggers::new(tx_sys.clone(), tx_produce.clone()),
        InitialQueueTriggers::new(tx_sys, tx_produce, tx_publish, tx_app, tx_integration),
    )
}

#[derive(Clone)]
/// The entry points for kicking off a chain reaction of queue activity
pub struct QueueTriggers {
    /// Notify the SysValidation workflow to run, i.e. after handling gossip
    pub sys_validation: TriggerSender,
<<<<<<< HEAD
    /// Notify the ProduceDgdOps workflow to run, i.e. after InvokeCallZome
    pub produce_dgd_ops: TriggerSender,
=======
    /// Notify the ProduceDhtOps workflow to run, i.e. after InvokeCallZome
    pub produce_dht_ops: TriggerSender,
>>>>>>> master
}

/// The triggers to run once at the start of a cell
pub struct InitialQueueTriggers {
    /// These triggers can only be run once
    /// so they are private
    sys_validation: TriggerSender,
<<<<<<< HEAD
    produce_dgd_ops: TriggerSender,
    publish_dgd_ops: TriggerSender,
    app_validation: TriggerSender,
    integrate_dgd_ops: TriggerSender,
=======
    produce_dht_ops: TriggerSender,
    publish_dht_ops: TriggerSender,
    app_validation: TriggerSender,
    integrate_dht_ops: TriggerSender,
>>>>>>> master
}

impl QueueTriggers {
    /// Create a new queue trigger
<<<<<<< HEAD
    pub fn new(sys_validation: TriggerSender, produce_dgd_ops: TriggerSender) -> Self {
        Self {
            sys_validation,
            produce_dgd_ops,
=======
    pub fn new(sys_validation: TriggerSender, produce_dht_ops: TriggerSender) -> Self {
        Self {
            sys_validation,
            produce_dht_ops,
>>>>>>> master
        }
    }
}

impl InitialQueueTriggers {
    fn new(
        sys_validation: TriggerSender,
<<<<<<< HEAD
        produce_dgd_ops: TriggerSender,
        publish_dgd_ops: TriggerSender,
        app_validation: TriggerSender,
        integrate_dgd_ops: TriggerSender,
    ) -> Self {
        Self {
            sys_validation,
            produce_dgd_ops,
            publish_dgd_ops,
            app_validation,
            integrate_dgd_ops,
=======
        produce_dht_ops: TriggerSender,
        publish_dht_ops: TriggerSender,
        app_validation: TriggerSender,
        integrate_dht_ops: TriggerSender,
    ) -> Self {
        Self {
            sys_validation,
            produce_dht_ops,
            publish_dht_ops,
            app_validation,
            integrate_dht_ops,
>>>>>>> master
        }
    }

    /// Initialize all the workflows once.
    pub fn initialize_workflows(mut self) {
        self.sys_validation.trigger();
        self.app_validation.trigger();
<<<<<<< HEAD
        self.publish_dgd_ops.trigger();
        self.integrate_dgd_ops.trigger();
        self.produce_dgd_ops.trigger();
=======
        self.publish_dht_ops.trigger();
        self.integrate_dht_ops.trigger();
        self.produce_dht_ops.trigger();
>>>>>>> master
    }
}
/// The means of nudging a queue consumer to tell it to look for more work
#[derive(Clone)]
pub struct TriggerSender(mpsc::Sender<()>);

/// The receiving end of a queue trigger channel
pub struct TriggerReceiver(mpsc::Receiver<()>);

impl TriggerSender {
    /// Create a new channel for waking a consumer
    ///
    /// The channel buffer is set to num_cpus to deal with the potential
    /// inconsistency from the perspective of any particular CPU thread
    pub fn new() -> (TriggerSender, TriggerReceiver) {
        let (tx, rx) = mpsc::channel(num_cpus::get());
        (TriggerSender(tx), TriggerReceiver(rx))
    }

    /// Lazily nudge the consumer task, ignoring the case where the consumer
    /// already has a pending trigger signal
    pub fn trigger(&mut self) {
        match self.0.try_send(()) {
            Err(mpsc::error::TrySendError::Closed(_)) => {
                tracing::warn!(
                    "Queue consumer trigger was sent while Cell is shutting down: ignoring."
                );
            }
            Err(mpsc::error::TrySendError::Full(_)) => {}
            Ok(()) => {}
        };
    }
}

impl TriggerReceiver {
    /// Listen for one or more items to come through, draining the channel
    /// each time. Bubble up errors on empty channel.
    pub async fn listen(&mut self) -> Result<(), QueueTriggerClosedError> {
        use tokio::sync::mpsc::error::TryRecvError;

        // wait for next item
        if self.0.recv().await.is_some() {
            // drain the channel
            loop {
                match self.0.try_recv() {
                    Err(TryRecvError::Closed) => return Err(QueueTriggerClosedError),
                    Err(TryRecvError::Empty) => return Ok(()),
                    Ok(()) => {}
                }
            }
        } else {
            Err(QueueTriggerClosedError)
        }
    }
}

/// A lazy Writer factory which can only be used once.
///
/// This is a way of encapsulating an EnvironmentWrite so that it can only be
/// used to create a single Writer before being consumed.
#[derive(Constructor, From)]
pub struct OneshotWriter(EnvironmentWrite);

impl OneshotWriter {
    /// Create the writer and pass it into a closure.
    pub fn with_writer<F>(self, f: F) -> Result<(), WorkspaceError>
    where
        F: FnOnce(&mut Writer) -> Result<(), WorkspaceError> + Send,
    {
        let env_ref = self.0.guard();
        env_ref.with_commit::<WorkspaceError, (), _>(|w| {
            f(w)?;
            Ok(())
        })?;
        Ok(())
    }
}

/// Declares whether a workflow has exhausted the queue or not
#[derive(Clone, Debug, PartialEq)]
pub enum WorkComplete {
    /// The queue has been exhausted
    Complete,
    /// Items still remain on the queue
    Incomplete,
}

/// The only error possible when attempting to trigger: the channel is closed
#[derive(Debug, Display, thiserror::Error)]
pub struct QueueTriggerClosedError;

/// Inform a workflow to run a job or shutdown
enum Job {
    Run,
    Shutdown,
}

/// Wait for the next job or exit command
async fn next_job_or_exit(
    rx: &mut TriggerReceiver,
    stop: &mut sync::broadcast::Receiver<()>,
) -> Job {
    if stop.try_recv().is_ok() {
        return Job::Shutdown;
    }
    // Check for shutdown or next job
    let next_job = rx.listen();
    let kill = stop.recv();
    tokio::pin!(next_job);
    tokio::pin!(kill);

    if let Either::Left((Err(_), _)) | Either::Right((_, _)) =
        futures::future::select(next_job, kill).await
    {
        Job::Shutdown
    } else {
        Job::Run
    }
}
