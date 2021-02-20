//! A Cell is an "instance" of AIngle DNA.
//!
//! It combines an AgentPubKey with a Dna to create a SourceChain, upon which
//! Elements can be added. A constructed Cell is guaranteed to have a valid
//! SourceChain which has already undergone Genesis.

use super::api::ZomeCall;
use super::interface::SignalBroadcaster;
use super::manager::ManagedTaskAdd;
use crate::conductor::api::error::ConductorApiError;
use crate::conductor::api::CellConductorApi;
use crate::conductor::api::CellConductorApiT;
use crate::conductor::cell::error::CellResult;
use crate::conductor::entry_def_store::get_entry_def_from_ids;
use crate::conductor::handle::ConductorHandle;
use crate::core::queue_consumer::spawn_queue_consumer_tasks;
use crate::core::queue_consumer::InitialQueueTriggers;
use crate::core::queue_consumer::QueueTriggers;
use crate::core::ribosome::guest_callback::init::InitResult;
use crate::core::ribosome::real_ribosome::RealRibosome;
use crate::core::ribosome::ZomeCallInvocation;
use crate::core::workflow::call_zome_workflow;
use crate::core::workflow::error::WorkflowError;
use crate::core::workflow::genesis_workflow::genesis_workflow;
<<<<<<< HEAD
use crate::core::workflow::incoming_dgd_ops_workflow::incoming_dgd_ops_workflow;
use crate::core::workflow::initialize_zomes_workflow;
use crate::core::workflow::produce_dgd_ops_workflow::dgd_op_light::light_to_op;
=======
use crate::core::workflow::incoming_dht_ops_workflow::incoming_dht_ops_workflow;
use crate::core::workflow::initialize_zomes_workflow;
use crate::core::workflow::produce_dht_ops_workflow::dht_op_light::light_to_op;
>>>>>>> master
use crate::core::workflow::CallZomeWorkflowArgs;
use crate::core::workflow::CallZomeWorkspace;
use crate::core::workflow::GenesisWorkflowArgs;
use crate::core::workflow::GenesisWorkspace;
use crate::core::workflow::InitializeZomesWorkflowArgs;
use crate::core::workflow::ZomeCallResult;
use call_zome_workflow::call_zome_workspace_lock::CallZomeWorkspaceLock;
use error::CellError;
use fallible_iterator::FallibleIterator;
use futures::future::FutureExt;
<<<<<<< HEAD
use hash_type::AnyDgd;
=======
use hash_type::AnyDht;
>>>>>>> master
use aingle_hash::*;
use aingle_cascade::authority;
use aingle_lmdb::db::GetDb;
use aingle_lmdb::env::EnvironmentRead;
use aingle_lmdb::env::EnvironmentWrite;
use aingle_lmdb::env::ReadManager;
use aingle_p2p::AIngleP2pCellT;
<<<<<<< HEAD
use aingle_middleware_bytes::SerializedBytes;
=======
use aingle_serialized_bytes::SerializedBytes;
>>>>>>> master
use aingle_state::prelude::*;
use aingle_types::prelude::*;
use observability::OpenSpanExt;
use std::hash::Hash;
use std::hash::Hasher;
use tokio::sync;
use tracing::*;
use tracing_futures::Instrument;
use validation_package::ValidationPackageDb;

mod validation_package;

#[allow(missing_docs)]
pub mod error;

#[cfg(test)]
mod gossip_test;

#[cfg(test)]
mod test;

impl Hash for Cell {
    fn hash<H>(&self, state: &mut H)
    where
        H: Hasher,
    {
        self.id.hash(state);
    }
}

impl PartialEq for Cell {
    fn eq(&self, other: &Self) -> bool {
        self.id == other.id
    }
}

/// A Cell is a grouping of the resources necessary to run workflows
/// on behalf of an agent. It does not have a lifetime of its own aside
/// from the lifetimes of the resources which it holds references to.
/// Any work it does is through running a workflow, passing references to
/// the resources needed to complete that workflow.
///
/// A Cell is guaranteed to contain a Source Chain which has undergone
/// Genesis.
///
/// The [Conductor] manages a collection of Cells, and will call functions
/// on the Cell when a Conductor API method is called (either a
/// [CellConductorApi] or an [AppInterfaceApi])
pub struct Cell<Api = CellConductorApi, P2pCell = aingle_p2p::AIngleP2pCell>
where
    Api: CellConductorApiT,
    P2pCell: aingle_p2p::AIngleP2pCellT,
{
    id: CellId,
    conductor_api: Api,
    env: EnvironmentWrite,
    aingle_p2p_cell: P2pCell,
    queue_triggers: QueueTriggers,
}

impl Cell {
    /// Constructor for a Cell. The SourceChain will be created, and genesis
    /// will be run if necessary. A Cell will not be created if the SourceChain
    /// is not ready to be used.
    pub async fn create(
        id: CellId,
        conductor_handle: ConductorHandle,
        env: EnvironmentWrite,
        aingle_p2p_cell: aingle_p2p::AIngleP2pCell,
        managed_task_add_sender: sync::mpsc::Sender<ManagedTaskAdd>,
        managed_task_stop_broadcaster: sync::broadcast::Sender<()>,
    ) -> CellResult<(Self, InitialQueueTriggers)> {
        let conductor_api = CellConductorApi::new(conductor_handle.clone(), id.clone());

        // check if genesis has been run
        let has_genesis = {
            // check if genesis ran on source chain buf
            SourceChainBuf::new(env.clone().into())?.has_genesis()
        };

        if has_genesis {
            tokio::spawn({
                let mut network = aingle_p2p_cell.clone();
                async move { network.join().await }
            });
            let (queue_triggers, initial_queue_triggers) = spawn_queue_consumer_tasks(
                &env,
                aingle_p2p_cell.clone(),
                conductor_api.clone(),
                managed_task_add_sender,
                managed_task_stop_broadcaster,
            )
            .await;

            Ok((
                Self {
                    id,
                    conductor_api,
                    env,
                    aingle_p2p_cell,
                    queue_triggers,
                },
                initial_queue_triggers,
            ))
        } else {
            Err(CellError::CellWithoutGenesis(id))
        }
    }

    /// Performs the Genesis workflow the Cell, ensuring that its initial
    /// elements are committed. This is a prerequisite for any other interaction
    /// with the SourceChain
    pub async fn genesis(
        id: CellId,
        conductor_handle: ConductorHandle,
        cell_env: EnvironmentWrite,
        membrane_proof: Option<SerializedBytes>,
    ) -> CellResult<()> {
        // get the dna
        let dna_file = conductor_handle
            .get_dna(id.dna_hash())
            .await
            .ok_or(CellError::DnaMissing)?;

        let conductor_api = CellConductorApi::new(conductor_handle, id.clone());

        // run genesis
        let workspace = GenesisWorkspace::new(cell_env.clone().into())
            .await
            .map_err(ConductorApiError::from)
            .map_err(Box::new)?;
        let args = GenesisWorkflowArgs::new(dna_file, id.agent_pubkey().clone(), membrane_proof);

        genesis_workflow(workspace, cell_env.clone().into(), conductor_api, args)
            .await
            .map_err(Box::new)
            .map_err(ConductorApiError::from)
            .map_err(Box::new)?;
        Ok(())
    }

    fn dna_hash(&self) -> &DnaHash {
        &self.id.dna_hash()
    }

    #[allow(unused)]
    fn agent_pubkey(&self) -> &AgentPubKey {
        &self.id.agent_pubkey()
    }

    /// Accessor
    pub fn id(&self) -> &CellId {
        &self.id
    }

    /// Access a network sender that is partially applied to this cell's DnaHash/AgentPubKey
    pub fn aingle_p2p_cell(&self) -> &aingle_p2p::AIngleP2pCell {
        &self.aingle_p2p_cell
    }

    async fn signal_broadcaster(&self) -> SignalBroadcaster {
        self.conductor_api.signal_broadcaster().await
    }

    #[instrument(skip(self, evt))]
    /// Entry point for incoming messages from the network that need to be handled
    pub async fn handle_aingle_p2p_event(
        &self,
        evt: aingle_p2p::event::AIngleP2pEvent,
    ) -> CellResult<()> {
        use aingle_p2p::event::AIngleP2pEvent::*;
        match evt {
            PutAgentInfoSigned { .. } | GetAgentInfoSigned { .. } | QueryAgentInfoSigned { .. } => {
                // PutAgentInfoSigned needs to be handled at the conductor level where the p2p
                // store lives.
                unreachable!()
            }
            CallRemote {
                span_context: _,
                from_agent,
                zome_name,
                fn_name,
                cap,
                respond,
                payload,
                ..
            } => {
                async {
                    let res = self
                        .handle_call_remote(from_agent, zome_name, fn_name, cap, payload)
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("call_remote"))
                .await;
            }
            Publish {
                span_context,
                respond,
                from_agent,
                request_validation_receipt,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                ops,
                ..
            } => {
                async {
                    tracing::Span::set_current_context(span_context);
                    let res = self
<<<<<<< HEAD
                        .handle_publish(from_agent, request_validation_receipt, dgd_hash, ops)
=======
                        .handle_publish(from_agent, request_validation_receipt, dht_hash, ops)
>>>>>>> master
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_publish"))
                .await;
            }
            GetValidationPackage {
                span_context: _,
                respond,
                header_hash,
                ..
            } => {
                async {
                    let res = self
                        .handle_get_validation_package(header_hash)
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_get_validation_package"))
                .await;
            }
            Get {
                span_context: _,
                respond,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                options,
                ..
            } => {
                async {
                    let res = self
<<<<<<< HEAD
                        .handle_get(dgd_hash, options)
=======
                        .handle_get(dht_hash, options)
>>>>>>> master
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_get"))
                .await;
            }
            GetMeta {
                span_context: _,
                respond,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                options,
                ..
            } => {
                async {
                    let res = self
<<<<<<< HEAD
                        .handle_get_meta(dgd_hash, options)
=======
                        .handle_get_meta(dht_hash, options)
>>>>>>> master
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_get_meta"))
                .await;
            }
            GetLinks {
                span_context: _,
                respond,
                link_key,
                options,
                ..
            } => {
                async {
                    let res = self
                        .handle_get_links(link_key, options)
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_get_links"))
                .await;
            }
            GetAgentActivity {
                span_context: _,
                respond,
                agent,
                query,
                options,
                ..
            } => {
                async {
                    let res = self
                        .handle_get_agent_activity(agent, query, options)
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_get_agent_activity"))
                .await;
            }
            ValidationReceiptReceived {
                span_context: _,
                respond,
                receipt,
                ..
            } => {
                async {
                    let res = self
                        .handle_validation_receipt(receipt)
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_validation_receipt_received"))
                .await;
            }
            FetchOpHashesForConstraints {
                span_context: _,
                respond,
<<<<<<< HEAD
                dgd_arc,
=======
                dht_arc,
>>>>>>> master
                since,
                until,
                ..
            } => {
                async {
                    let res = self
<<<<<<< HEAD
                        .handle_fetch_op_hashes_for_constraints(dgd_arc, since, until)
=======
                        .handle_fetch_op_hashes_for_constraints(dht_arc, since, until)
>>>>>>> master
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_fetch_op_hashes_for_constraints"))
                .await;
            }
            FetchOpHashData {
                span_context: _,
                respond,
                op_hashes,
                ..
            } => {
                async {
                    let res = self
                        .handle_fetch_op_hash_data(op_hashes)
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_fetch_op_hash_data"))
                .await;
            }
            SignNetworkData {
                span_context: _,
                respond,
                ..
            } => {
                async {
                    let res = self
                        .handle_sign_network_data()
                        .await
                        .map_err(aingle_p2p::AIngleP2pError::other);
                    respond.respond(Ok(async move { res }.boxed().into()));
                }
                .instrument(debug_span!("cell_handle_sign_network_data"))
                .await;
            }
        }
        Ok(())
    }

<<<<<<< HEAD
    #[instrument(skip(self, _request_validation_receipt, _dgd_hash, ops))]
=======
    #[instrument(skip(self, _request_validation_receipt, _dht_hash, ops))]
>>>>>>> master
    /// we are receiving a "publish" event from the network
    async fn handle_publish(
        &self,
        from_agent: AgentPubKey,
        _request_validation_receipt: bool,
<<<<<<< HEAD
        _dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
    ) -> CellResult<()> {
        incoming_dgd_ops_workflow(
=======
        _dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
    ) -> CellResult<()> {
        incoming_dht_ops_workflow(
>>>>>>> master
            &self.env,
            self.queue_triggers.sys_validation.clone(),
            ops,
            Some(from_agent),
        )
        .await
        .map_err(Box::new)
        .map_err(ConductorApiError::from)
        .map_err(Box::new)?;
        Ok(())
    }

    #[instrument(skip(self))]
    /// a remote node is attempting to retrieve a validation package
    #[tracing::instrument(skip(self), level = "trace")]
    async fn handle_get_validation_package(
        &self,
        header_hash: HeaderHash,
    ) -> CellResult<ValidationPackageResponse> {
        let env: EnvironmentRead = self.env.clone().into();

        // Get the header
        let databases = ValidationPackageDb::create(env.clone())?;
        let mut cascade = databases.cascade();
        let header = match cascade
            .retrieve_header(header_hash, Default::default())
            .await?
        {
            Some(shh) => shh.into_header_and_signature().0,
            None => return Ok(None.into()),
        };

        let ribosome = self.get_ribosome().await?;

        // This agent is the author so get the validation package from the source chain
        if header.author() == self.id.agent_pubkey() {
            validation_package::get_as_author(
                header,
                env,
                &ribosome,
                &self.conductor_api,
                &self.aingle_p2p_cell,
            )
            .await
        } else {
            validation_package::get_as_authority(
                header,
                env,
                &ribosome.dna_file,
                &self.conductor_api,
            )
            .await
        }
    }

    #[instrument(skip(self, options))]
    /// a remote node is asking us for entry data
    async fn handle_get(
        &self,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: aingle_p2p::event::GetOptions,
    ) -> CellResult<GetElementResponse> {
        debug!("handling get");
        // TODO: Later we will need more get types but for now
        // we can just have these defaults depending on whether or not
        // the hash is an entry or header.
        // In the future we should use GetOptions to choose which get to run.
<<<<<<< HEAD
        let r = match *dgd_hash.hash_type() {
            AnyDgd::Entry => self.handle_get_entry(dgd_hash.into(), options).await,
            AnyDgd::Header => self.handle_get_element(dgd_hash.into()).await,
=======
        let r = match *dht_hash.hash_type() {
            AnyDht::Entry => self.handle_get_entry(dht_hash.into(), options).await,
            AnyDht::Header => self.handle_get_element(dht_hash.into()).await,
>>>>>>> master
        };
        if let Err(e) = &r {
            error!(msg = "Error handling a get", ?e, agent = ?self.id.agent_pubkey());
        }
        r
    }

    #[instrument(skip(self, options))]
    async fn handle_get_entry(
        &self,
        hash: EntryHash,
        options: aingle_p2p::event::GetOptions,
    ) -> CellResult<GetElementResponse> {
        let env = self.env.clone();
        authority::handle_get_entry(env, hash, options).map_err(Into::into)
    }

    #[tracing::instrument(skip(self))]
    async fn handle_get_element(&self, hash: HeaderHash) -> CellResult<GetElementResponse> {
        let env = self.env.clone();
        authority::handle_get_element(env, hash).map_err(Into::into)
    }

<<<<<<< HEAD
    #[instrument(skip(self, _dgd_hash, _options))]
    /// a remote node is asking us for metadata
    async fn handle_get_meta(
        &self,
        _dgd_hash: aingle_hash::AnyDgdHash,
=======
    #[instrument(skip(self, _dht_hash, _options))]
    /// a remote node is asking us for metadata
    async fn handle_get_meta(
        &self,
        _dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        _options: aingle_p2p::event::GetMetaOptions,
    ) -> CellResult<MetadataSet> {
        unimplemented!()
    }

    #[instrument(skip(self, options))]
    /// a remote node is asking us for links
    // TODO: Right now we are returning all the full headers
    // We could probably send some smaller types instead of the full headers
    // if we are careful.
    fn handle_get_links(
        &self,
        link_key: WireLinkMetaKey,
        options: aingle_p2p::event::GetLinksOptions,
    ) -> CellResult<GetLinksResponse> {
        debug!(id = ?self.id());
        let env = self.env.clone();
        authority::handle_get_links(env.into(), link_key, options).map_err(Into::into)
    }

    #[instrument(skip(self, options))]
    fn handle_get_agent_activity(
        &self,
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: aingle_p2p::event::GetActivityOptions,
    ) -> CellResult<AgentActivityResponse> {
        let env = self.env.clone();
        authority::handle_get_agent_activity(env.into(), agent, query, options).map_err(Into::into)
    }

    /// a remote agent is sending us a validation receipt.
    #[tracing::instrument(skip(self))]
    async fn handle_validation_receipt(&self, _receipt: SerializedBytes) -> CellResult<()> {
        unimplemented!()
    }

<<<<<<< HEAD
    #[instrument(skip(self, dgd_arc, since, until))]
    /// the network module is requesting a list of dgd op hashes
    fn handle_fetch_op_hashes_for_constraints(
        &self,
        dgd_arc: aingle_p2p::dgd_arc::DgdArc,
        since: Timestamp,
        until: Timestamp,
    ) -> CellResult<Vec<DgdOpHash>> {
        let env_ref = self.env.guard();
        let reader = env_ref.reader()?;
        let integrated_dgd_ops = IntegratedDgdOpsBuf::new(self.env().clone().into())?;
        let result: Vec<DgdOpHash> = integrated_dgd_ops
            .query(&reader, Some(since), Some(until), Some(dgd_arc))?
=======
    #[instrument(skip(self, dht_arc, since, until))]
    /// the network module is requesting a list of dht op hashes
    fn handle_fetch_op_hashes_for_constraints(
        &self,
        dht_arc: aingle_p2p::dht_arc::DhtArc,
        since: Timestamp,
        until: Timestamp,
    ) -> CellResult<Vec<DhtOpHash>> {
        let env_ref = self.env.guard();
        let reader = env_ref.reader()?;
        let integrated_dht_ops = IntegratedDhtOpsBuf::new(self.env().clone().into())?;
        let result: Vec<DhtOpHash> = integrated_dht_ops
            .query(&reader, Some(since), Some(until), Some(dht_arc))?
>>>>>>> master
            .map(|(k, _)| Ok(k))
            .collect()?;
        Ok(result)
    }

    #[instrument(skip(self, op_hashes))]
<<<<<<< HEAD
    /// The network module is requesting the content for dgd ops
    async fn handle_fetch_op_hash_data(
        &self,
        op_hashes: Vec<aingle_hash::DgdOpHash>,
    ) -> CellResult<
        Vec<(
            aingle_hash::AnyDgdHash,
            aingle_hash::DgdOpHash,
            aingle_types::dgd_op::DgdOp,
        )>,
    > {
        let integrated_dgd_ops = IntegratedDgdOpsBuf::new(self.env().clone().into())?;
        let mut out = vec![];
        for op_hash in op_hashes {
            let val = integrated_dgd_ops.get(&op_hash)?;
=======
    /// The network module is requesting the content for dht ops
    async fn handle_fetch_op_hash_data(
        &self,
        op_hashes: Vec<aingle_hash::DhtOpHash>,
    ) -> CellResult<
        Vec<(
            aingle_hash::AnyDhtHash,
            aingle_hash::DhtOpHash,
            aingle_types::dht_op::DhtOp,
        )>,
    > {
        let integrated_dht_ops = IntegratedDhtOpsBuf::new(self.env().clone().into())?;
        let mut out = vec![];
        for op_hash in op_hashes {
            let val = integrated_dht_ops.get(&op_hash)?;
>>>>>>> master
            if let Some(val) = val {
                let full_op = match &val.validation_status {
                    ValidationStatus::Valid => {
                        let cas = ElementBuf::vault(self.env.clone().into(), false)?;
                        light_to_op(val.op, &cas)?
                    }
                    ValidationStatus::Rejected => {
                        let cas = ElementBuf::rejected(self.env.clone().into())?;
                        light_to_op(val.op, &cas)?
                    }
                    ValidationStatus::Abandoned => todo!("Add when abandoned store is added"),
                };
<<<<<<< HEAD
                let basis = full_op.dgd_basis();
=======
                let basis = full_op.dht_basis();
>>>>>>> master
                out.push((basis, op_hash, full_op));
            }
        }
        Ok(out)
    }

    /// the network module would like this cell/agent to sign some data
    #[tracing::instrument(skip(self))]
    async fn handle_sign_network_data(&self) -> CellResult<Signature> {
        Ok(vec![0; 64].into())
    }

    /// When the Conductor determines that it's time to execute some [AutonomicProcess],
    /// whether scheduled or through an [AutonomicCue], this function gets called
    #[tracing::instrument(skip(self, process))]
    pub async fn handle_autonomic_process(&self, process: AutonomicProcess) -> CellResult<()> {
        match process {
            AutonomicProcess::SlowHeal => unimplemented!(),
            AutonomicProcess::HealthCheck => unimplemented!(),
        }
    }

    #[instrument(skip(self, from_agent, fn_name, cap, payload))]
    /// a remote agent is attempting a "call_remote" on this cell.
    async fn handle_call_remote(
        &self,
        from_agent: AgentPubKey,
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        payload: ExternIO,
    ) -> CellResult<SerializedBytes> {
        let invocation = ZomeCall {
            cell_id: self.id.clone(),
            zome_name,
            cap,
            payload,
            provenance: from_agent,
            fn_name,
        };
        // double ? because
        // - ConductorApiResult
        // - ZomeCallResult
        Ok(self.call_zome(invocation, None).await??.try_into()?)
    }

    /// Function called by the Conductor
    #[instrument(skip(self, call, workspace_lock))]
    pub async fn call_zome(
        &self,
        call: ZomeCall,
        workspace_lock: Option<CallZomeWorkspaceLock>,
    ) -> CellResult<ZomeCallResult> {
        // Check if init has run if not run it
        self.check_or_run_zome_init().await?;

        let arc = self.env();
        let keystore = arc.keystore().clone();

        // If there is no existing zome call then this is the root zome call
        let is_root_zome_call = workspace_lock.is_none();
        let workspace_lock = match workspace_lock {
            Some(l) => l,
            None => CallZomeWorkspaceLock::new(CallZomeWorkspace::new(arc.clone().into())?),
        };

        let conductor_api = self.conductor_api.clone();
        let signal_tx = self.signal_broadcaster().await;
        let ribosome = self.get_ribosome().await?;
        let invocation = ZomeCallInvocation::from_interface_call(conductor_api.clone(), call).await;

        let args = CallZomeWorkflowArgs {
            ribosome,
            invocation,
            conductor_api,
            signal_tx,
            is_root_zome_call,
        };
        Ok(call_zome_workflow(
            workspace_lock,
            self.aingle_p2p_cell.clone(),
            keystore,
            arc.clone().into(),
            args,
<<<<<<< HEAD
            self.queue_triggers.produce_dgd_ops.clone(),
=======
            self.queue_triggers.produce_dht_ops.clone(),
>>>>>>> master
        )
        .await
        .map_err(Box::new)?)
    }

    /// Check if each Zome's init callback has been run, and if not, run it.
    #[tracing::instrument(skip(self))]
    async fn check_or_run_zome_init(&self) -> CellResult<()> {
        // If not run it
        let env = self.env.clone();
        let keystore = env.keystore().clone();
        let id = self.id.clone();
        let conductor_api = self.conductor_api.clone();
        // Create the workspace
        let workspace = CallZomeWorkspace::new(self.env().clone().into())
            .map_err(WorkflowError::from)
            .map_err(Box::new)?;

        // Check if initialization has run
        if workspace.source_chain.has_initialized() {
            return Ok(());
        }
        trace!("running init");

        // get the dna
        let dna_file = conductor_api
            .get_dna(id.dna_hash())
            .await
            .ok_or(CellError::DnaMissing)?;
        let dna_def = dna_file.dna_def().clone();

        // Get the ribosome
        let ribosome = RealRibosome::new(dna_file);

        // Run the workflow
        let args = InitializeZomesWorkflowArgs { dna_def, ribosome };
        let init_result = initialize_zomes_workflow(
            workspace,
            self.aingle_p2p_cell.clone(),
            keystore,
            env.clone().into(),
            args,
        )
        .await
        .map_err(Box::new)?;
        trace!(?init_result);
        match init_result {
            InitResult::Pass => {}
            r => return Err(CellError::InitFailed(r)),
        }
        Ok(())
    }

    /// Delete all data associated with this Cell by deleting the associated
    /// LMDB environment. Completely reverses Cell creation.
    #[tracing::instrument(skip(self))]
    pub async fn destroy(self) -> CellResult<()> {
        let path = self.env.path().clone();
        // Remove db from global map
        // Delete directory
        self.env
            .remove()
            .await
            .map_err(|e| CellError::Cleanup(e.to_string(), path))?;
        Ok(())
    }

    /// Instantiate a Ribosome for use by this Cell's workflows
    // TODO: reevaluate once Workflows are fully implemented (after B-01567)
    pub(crate) async fn get_ribosome(&self) -> CellResult<RealRibosome> {
        match self.conductor_api.get_dna(self.dna_hash()).await {
            Some(dna) => Ok(RealRibosome::new(dna)),
            None => Err(CellError::DnaMissing),
        }
    }

    /// Accessor for the LMDB environment backing this Cell
    // TODO: reevaluate once Workflows are fully implemented (after B-01567)
    pub(crate) fn env(&self) -> &EnvironmentWrite {
        &self.env
    }

    #[cfg(any(test, feature = "test_utils"))]
    /// Get the triggers for the cell
    /// Useful for testing when you want to
    /// Cause workflows to trigger
    pub(crate) fn triggers(&self) -> &QueueTriggers {
        &self.queue_triggers
    }
}
