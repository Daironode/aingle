use crate::actor::*;
use crate::event::*;
use crate::*;

use futures::future::FutureExt;

use crate::types::AgentPubKeyExt;

use ghost_actor::dependencies::tracing;
use ghost_actor::dependencies::tracing_futures::Instrument;

use aingle_zome_types::zome::FunctionName;
use kitsune_p2p::actor::KitsuneP2pSender;
use kitsune_p2p::agent_store::AgentInfoSigned;

pub(crate) struct AIngleP2pActor {
    evt_sender: futures::channel::mpsc::Sender<AIngleP2pEvent>,
    kitsune_p2p: ghost_actor::GhostSender<kitsune_p2p::actor::KitsuneP2p>,
}

impl ghost_actor::GhostControlHandler for AIngleP2pActor {}

impl AIngleP2pActor {
    /// constructor
    pub async fn new(
        config: kitsune_p2p::KitsuneP2pConfig,
        tls_config: kitsune_p2p::dependencies::kitsune_p2p_proxy::TlsConfig,
        channel_factory: ghost_actor::actor_builder::GhostActorChannelFactory<Self>,
        evt_sender: futures::channel::mpsc::Sender<AIngleP2pEvent>,
    ) -> AIngleP2pResult<Self> {
        let (kitsune_p2p, kitsune_p2p_events) =
            kitsune_p2p::spawn_kitsune_p2p(config, tls_config).await?;

        channel_factory.attach_receiver(kitsune_p2p_events).await?;

        Ok(Self {
            evt_sender,
            kitsune_p2p,
        })
    }

    /// receiving an incoming request from a remote node
    #[allow(clippy::too_many_arguments)]
    fn handle_incoming_call_remote(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
        from_agent: AgentPubKey,
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        data: Vec<u8>,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let res = evt_sender
                .call_remote(
                    dna_hash,
                    to_agent,
                    from_agent,
                    zome_name,
                    fn_name,
                    cap,
                    ExternIO::from(data),
                )
                .await;
            res.map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .boxed()
        .into())
    }

    /// receiving an incoming get request from a remote node
<<<<<<< HEAD
    #[tracing::instrument(skip(self, dna_hash, to_agent, dgd_hash, options), level = "trace")]
=======
    #[tracing::instrument(skip(self, dna_hash, to_agent, dht_hash, options), level = "trace")]
>>>>>>> master
    fn handle_incoming_get(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: event::GetOptions,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
<<<<<<< HEAD
            let res = evt_sender.get(dna_hash, to_agent, dgd_hash, options).await;
=======
            let res = evt_sender.get(dna_hash, to_agent, dht_hash, options).await;
>>>>>>> master
            res.and_then(|r| Ok(SerializedBytes::try_from(r)?))
                .map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .instrument(tracing::debug_span!("incoming_get_task"))
        .boxed()
        .into())
    }

    /// receiving an incoming get_meta request from a remote node
    fn handle_incoming_get_meta(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: event::GetMetaOptions,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let res = evt_sender
<<<<<<< HEAD
                .get_meta(dna_hash, to_agent, dgd_hash, options)
=======
                .get_meta(dna_hash, to_agent, dht_hash, options)
>>>>>>> master
                .await;
            res.and_then(|r| Ok(SerializedBytes::try_from(r)?))
                .map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .boxed()
        .into())
    }

    /// receiving an incoming get_links request from a remote node
    fn handle_incoming_get_links(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
        link_key: WireLinkMetaKey,
        options: event::GetLinksOptions,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let res = evt_sender
                .get_links(dna_hash, to_agent, link_key, options)
                .await;
            res.and_then(|r| Ok(SerializedBytes::try_from(r)?))
                .map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .boxed()
        .into())
    }

    /// receiving an incoming get_links request from a remote node
    fn handle_incoming_get_agent_activity(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: event::GetActivityOptions,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let res = evt_sender
                .get_agent_activity(dna_hash, to_agent, agent, query, options)
                .await;
            res.and_then(|r| Ok(SerializedBytes::try_from(r)?))
                .map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .boxed()
        .into())
    }

    /// receiving an incoming publish from a remote node
    fn handle_incoming_publish(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
        from_agent: AgentPubKey,
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
>>>>>>> master
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<()> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            evt_sender
                .publish(
                    dna_hash,
                    to_agent,
                    from_agent,
                    request_validation_receipt,
<<<<<<< HEAD
                    dgd_hash,
=======
                    dht_hash,
>>>>>>> master
                    ops,
                )
                .await?;
            Ok(())
        }
        .boxed()
        .into())
    }

    /// receiving an incoming validation receipt from a remote node
    fn handle_incoming_validation_receipt(
        &mut self,
        dna_hash: DnaHash,
        agent_pub_key: AgentPubKey,
        receipt: Vec<u8>,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let receipt: SerializedBytes = UnsafeBytes::from(receipt).into();
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            evt_sender
                .validation_receipt_received(dna_hash, agent_pub_key, receipt)
                .await?;

            // validation receipts don't need a response
            // send back an empty vec for now
            Ok(Vec::with_capacity(0))
        }
        .boxed()
        .into())
    }

    /// Receiving an incoming validation package request
    fn handle_incoming_get_validation_package(
        &mut self,
        dna_hash: DnaHash,
        agent_pub_key: AgentPubKey,
        header_hash: HeaderHash,
    ) -> kitsune_p2p::actor::KitsuneP2pHandlerResult<Vec<u8>> {
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let res = evt_sender
                .get_validation_package(dna_hash, agent_pub_key, header_hash)
                .await;

            res.and_then(|r| Ok(SerializedBytes::try_from(r)?))
                .map_err(kitsune_p2p::KitsuneP2pError::from)
                .map(|res| UnsafeBytes::from(res).into())
        }
        .boxed()
        .into())
    }
}

impl ghost_actor::GhostHandler<kitsune_p2p::event::KitsuneP2pEvent> for AIngleP2pActor {}

impl kitsune_p2p::event::KitsuneP2pEventHandler for AIngleP2pActor {
    /// We need to store signed agent info.
    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_put_agent_info_signed(
        &mut self,
        input: kitsune_p2p::event::PutAgentInfoSignedEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<()> {
        let kitsune_p2p::event::PutAgentInfoSignedEvt {
            space,
            agent,
            agent_info_signed,
        } = input;
        let space = DnaHash::from_kitsune(&space);
        let agent = AgentPubKey::from_kitsune(&agent);
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            Ok(evt_sender
                .put_agent_info_signed(space, agent, agent_info_signed)
                .await?)
        }
        .boxed()
        .into())
    }

    /// We need to get previously stored agent info.
    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_get_agent_info_signed(
        &mut self,
        input: kitsune_p2p::event::GetAgentInfoSignedEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<Option<AgentInfoSigned>> {
        let kitsune_p2p::event::GetAgentInfoSignedEvt { space, agent } = input;
        let h_space = DnaHash::from_kitsune(&space);
        let h_agent = AgentPubKey::from_kitsune(&agent);
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            Ok(evt_sender
                .get_agent_info_signed(h_space, h_agent, space, agent)
                .await?)
        }
        .boxed()
        .into())
    }

    /// We need to get previously stored agent info.
    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_query_agent_info_signed(
        &mut self,
        input: kitsune_p2p::event::QueryAgentInfoSignedEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<Vec<AgentInfoSigned>> {
        let kitsune_p2p::event::QueryAgentInfoSignedEvt { space, agent } = input;
        let h_space = DnaHash::from_kitsune(&space);
        let h_agent = AgentPubKey::from_kitsune(&agent);
        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            Ok(evt_sender
                .query_agent_info_signed(h_space, h_agent, space, agent)
                .await?)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self, space, to_agent, from_agent, payload), level = "trace")]
    fn handle_call(
        &mut self,
        space: Arc<kitsune_p2p::KitsuneSpace>,
        to_agent: Arc<kitsune_p2p::KitsuneAgent>,
        from_agent: Arc<kitsune_p2p::KitsuneAgent>,
        payload: Vec<u8>,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<Vec<u8>> {
        let space = DnaHash::from_kitsune(&space);
        let to_agent = AgentPubKey::from_kitsune(&to_agent);
        let from_agent = AgentPubKey::from_kitsune(&from_agent);

        let request =
            crate::wire::WireMessage::decode(payload.as_ref()).map_err(AIngleP2pError::from)?;

        match request {
            crate::wire::WireMessage::CallRemote {
                zome_name,
                fn_name,
                cap,
                data,
            } => self.handle_incoming_call_remote(
                space, to_agent, from_agent, zome_name, fn_name, cap, data,
            ),
<<<<<<< HEAD
            crate::wire::WireMessage::Get { dgd_hash, options } => {
                self.handle_incoming_get(space, to_agent, dgd_hash, options)
            }
            crate::wire::WireMessage::GetMeta { dgd_hash, options } => {
                self.handle_incoming_get_meta(space, to_agent, dgd_hash, options)
=======
            crate::wire::WireMessage::Get { dht_hash, options } => {
                self.handle_incoming_get(space, to_agent, dht_hash, options)
            }
            crate::wire::WireMessage::GetMeta { dht_hash, options } => {
                self.handle_incoming_get_meta(space, to_agent, dht_hash, options)
>>>>>>> master
            }
            crate::wire::WireMessage::GetLinks { link_key, options } => {
                self.handle_incoming_get_links(space, to_agent, link_key, options)
            }
            crate::wire::WireMessage::GetAgentActivity {
                agent,
                query,
                options,
            } => self.handle_incoming_get_agent_activity(space, to_agent, agent, query, options),
            // aingle_p2p never publishes via request
            // these only occur on broadcasts
            crate::wire::WireMessage::Publish { .. } => {
                Err(AIngleP2pError::invalid_p2p_message(
                    "invalid: publish is a broadcast type, not a request".to_string(),
                )
                .into())
            }
            crate::wire::WireMessage::ValidationReceipt { receipt } => {
                self.handle_incoming_validation_receipt(space, to_agent, receipt)
            }
            crate::wire::WireMessage::GetValidationPackage { header_hash } => {
                self.handle_incoming_get_validation_package(space, to_agent, header_hash)
            }
        }
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_notify(
        &mut self,
        space: Arc<kitsune_p2p::KitsuneSpace>,
        to_agent: Arc<kitsune_p2p::KitsuneAgent>,
        from_agent: Arc<kitsune_p2p::KitsuneAgent>,
        payload: Vec<u8>,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<()> {
        let space = DnaHash::from_kitsune(&space);
        let to_agent = AgentPubKey::from_kitsune(&to_agent);
        let from_agent = AgentPubKey::from_kitsune(&from_agent);

        let request =
            crate::wire::WireMessage::decode(payload.as_ref()).map_err(AIngleP2pError::from)?;

        match request {
            // error on these call type messages
            crate::wire::WireMessage::CallRemote { .. }
            | crate::wire::WireMessage::Get { .. }
            | crate::wire::WireMessage::GetMeta { .. }
            | crate::wire::WireMessage::GetLinks { .. }
            | crate::wire::WireMessage::GetAgentActivity { .. }
            | crate::wire::WireMessage::GetValidationPackage { .. }
            | crate::wire::WireMessage::ValidationReceipt { .. } => {
                Err(AIngleP2pError::invalid_p2p_message(
                    "invalid call type message in a notify".to_string(),
                )
                .into())
            }
            crate::wire::WireMessage::Publish {
                request_validation_receipt,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                ops,
            } => self.handle_incoming_publish(
                space,
                to_agent,
                from_agent,
                request_validation_receipt,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                ops,
            ),
        }
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_gossip(
        &mut self,
        space: Arc<kitsune_p2p::KitsuneSpace>,
        to_agent: Arc<kitsune_p2p::KitsuneAgent>,
        from_agent: Arc<kitsune_p2p::KitsuneAgent>,
        op_hash: Arc<kitsune_p2p::KitsuneOpHash>,
        op_data: Vec<u8>,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<()> {
        let space = DnaHash::from_kitsune(&space);
        let to_agent = AgentPubKey::from_kitsune(&to_agent);
        let _from_agent = AgentPubKey::from_kitsune(&from_agent);
<<<<<<< HEAD
        let op_hash = DgdOpHash::from_kitsune(&op_hash);
        let op_data =
            crate::wire::WireDgdOpData::decode(op_data).map_err(AIngleP2pError::from)?;
=======
        let op_hash = DhtOpHash::from_kitsune(&op_hash);
        let op_data =
            crate::wire::WireDhtOpData::decode(op_data).map_err(AIngleP2pError::from)?;
>>>>>>> master
        self.handle_incoming_publish(
            space,
            to_agent,
            op_data.from_agent,
            false,
<<<<<<< HEAD
            op_data.dgd_hash,
=======
            op_data.dht_hash,
>>>>>>> master
            vec![(op_hash, op_data.op_data)],
        )
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_fetch_op_hashes_for_constraints(
        &mut self,
        input: kitsune_p2p::event::FetchOpHashesForConstraintsEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<Vec<Arc<kitsune_p2p::KitsuneOpHash>>>
    {
        let kitsune_p2p::event::FetchOpHashesForConstraintsEvt {
            space,
            agent,
<<<<<<< HEAD
            dgd_arc,
=======
            dht_arc,
>>>>>>> master
            since_utc_epoch_s,
            until_utc_epoch_s,
        } = input;
        let space = DnaHash::from_kitsune(&space);
        let agent = AgentPubKey::from_kitsune(&agent);
        let since = Timestamp(since_utc_epoch_s, 0);
        let until = Timestamp(until_utc_epoch_s, 0);

        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            Ok(evt_sender
<<<<<<< HEAD
                .fetch_op_hashes_for_constraints(space, agent, dgd_arc, since, until)
=======
                .fetch_op_hashes_for_constraints(space, agent, dht_arc, since, until)
>>>>>>> master
                .await?
                .into_iter()
                .map(|h| h.into_kitsune())
                .collect())
        }
        .boxed()
        .into())
    }

    #[allow(clippy::needless_collect)]
    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_fetch_op_hash_data(
        &mut self,
        input: kitsune_p2p::event::FetchOpHashDataEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<
        Vec<(Arc<kitsune_p2p::KitsuneOpHash>, Vec<u8>)>,
    > {
        let kitsune_p2p::event::FetchOpHashDataEvt {
            space,
            agent,
            op_hashes,
        } = input;
        let space = DnaHash::from_kitsune(&space);
        let agent = AgentPubKey::from_kitsune(&agent);
        let op_hashes = op_hashes
            .into_iter()
<<<<<<< HEAD
            .map(|h| DgdOpHash::from_kitsune(&h))
=======
            .map(|h| DhtOpHash::from_kitsune(&h))
>>>>>>> master
            // the allowance of clippy::needless_collcect refers to the following call
            .collect::<Vec<_>>();

        let evt_sender = self.evt_sender.clone();
        Ok(async move {
            let mut out = vec![];
<<<<<<< HEAD
            for (dgd_hash, op_hash, dgd_op) in evt_sender
=======
            for (dht_hash, op_hash, dht_op) in evt_sender
>>>>>>> master
                .fetch_op_hash_data(space, agent.clone(), op_hashes)
                .await?
            {
                out.push((
                    op_hash.into_kitsune(),
<<<<<<< HEAD
                    crate::wire::WireDgdOpData {
                        from_agent: agent.clone(),
                        dgd_hash,
                        op_data: dgd_op,
=======
                    crate::wire::WireDhtOpData {
                        from_agent: agent.clone(),
                        dht_hash,
                        op_data: dht_op,
>>>>>>> master
                    }
                    .encode()
                    .map_err(kitsune_p2p::KitsuneP2pError::other)?,
                ));
            }
            Ok(out)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_sign_network_data(
        &mut self,
        input: kitsune_p2p::event::SignNetworkDataEvt,
    ) -> kitsune_p2p::event::KitsuneP2pEventHandlerResult<kitsune_p2p::KitsuneSignature> {
        let space = DnaHash::from_kitsune(&input.space);
        let agent = AgentPubKey::from_kitsune(&input.agent);
        let fut = self
            .evt_sender
            .sign_network_data(space, agent, input.data.to_vec());
        Ok(async move {
            let sig = fut.await?.0;
            Ok(sig.into())
        }
        .boxed()
        .into())
    }
}

impl ghost_actor::GhostHandler<AIngleP2p> for AIngleP2pActor {}

impl AIngleP2pHandler for AIngleP2pActor {
    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_join(
        &mut self,
        dna_hash: DnaHash,
        agent_pub_key: AgentPubKey,
    ) -> AIngleP2pHandlerResult<()> {
        let space = dna_hash.into_kitsune();
        let agent = agent_pub_key.into_kitsune();

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move { Ok(kitsune_p2p.join(space, agent).await?) }
            .boxed()
            .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_leave(
        &mut self,
        dna_hash: DnaHash,
        agent_pub_key: AgentPubKey,
    ) -> AIngleP2pHandlerResult<()> {
        let space = dna_hash.into_kitsune();
        let agent = agent_pub_key.into_kitsune();

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move { Ok(kitsune_p2p.leave(space, agent).await?) }
            .boxed()
            .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_call_remote(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
        to_agent: AgentPubKey,
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        payload: ExternIO,
    ) -> AIngleP2pHandlerResult<SerializedBytes> {
        let space = dna_hash.into_kitsune();
        let to_agent = to_agent.into_kitsune();
        let from_agent = from_agent.into_kitsune();

        let req =
            crate::wire::WireMessage::call_remote(zome_name, fn_name, cap, payload).encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            let result: Vec<u8> = kitsune_p2p
                .rpc_single(space, to_agent, from_agent, req, None)
                .await?;
            Ok(UnsafeBytes::from(result).into())
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_publish(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
>>>>>>> master
        timeout_ms: Option<u64>,
    ) -> AIngleP2pHandlerResult<()> {
        let space = dna_hash.into_kitsune();
        let from_agent = from_agent.into_kitsune();
<<<<<<< HEAD
        let basis = dgd_hash.to_kitsune();

        let payload = crate::wire::WireMessage::publish(request_validation_receipt, dgd_hash, ops)
=======
        let basis = dht_hash.to_kitsune();

        let payload = crate::wire::WireMessage::publish(request_validation_receipt, dht_hash, ops)
>>>>>>> master
            .encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            kitsune_p2p
                .notify_multi(kitsune_p2p::actor::NotifyMulti {
                    space,
                    from_agent,
                    basis,
                    remote_agent_count: None, // default best-effort
                    timeout_ms,
                    payload,
                })
                .await?;
            Ok(())
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_get_validation_package(
        &mut self,
        input: actor::GetValidationPackage,
    ) -> AIngleP2pHandlerResult<ValidationPackageResponse> {
        let space = input.dna_hash.into_kitsune();
        let to_agent = input.request_from.into_kitsune();
        let from_agent = input.agent_pub_key.into_kitsune();

        let req = crate::wire::WireMessage::get_validation_package(input.header_hash).encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            let response = kitsune_p2p
                .rpc_single(space, to_agent, from_agent, req, None)
                .await?;
            let response = SerializedBytes::from(UnsafeBytes::from(response)).try_into()?;
            Ok(response)
        }
        .boxed()
        .into())
    }

<<<<<<< HEAD
    #[tracing::instrument(skip(self, dna_hash, from_agent, dgd_hash, options), level = "trace")]
=======
    #[tracing::instrument(skip(self, dna_hash, from_agent, dht_hash, options), level = "trace")]
>>>>>>> master
    fn handle_get(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: actor::GetOptions,
    ) -> AIngleP2pHandlerResult<Vec<GetElementResponse>> {
        let space = dna_hash.into_kitsune();
        let from_agent = from_agent.into_kitsune();
<<<<<<< HEAD
        let basis = dgd_hash.to_kitsune();
        let r_options: event::GetOptions = (&options).into();

        let payload = crate::wire::WireMessage::get(dgd_hash, r_options).encode()?;
=======
        let basis = dht_hash.to_kitsune();
        let r_options: event::GetOptions = (&options).into();

        let payload = crate::wire::WireMessage::get(dht_hash, r_options).encode()?;
>>>>>>> master

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            let result = kitsune_p2p
                .rpc_multi(kitsune_p2p::actor::RpcMulti {
                    space,
                    from_agent,
                    basis,
                    remote_agent_count: options.remote_agent_count,
                    timeout_ms: options.timeout_ms,
                    as_race: options.as_race,
                    race_timeout_ms: options.race_timeout_ms,
                    payload,
                })
                .instrument(tracing::debug_span!("rpc_multi"))
                .await?;

            let mut out = Vec::new();
            for item in result {
                let kitsune_p2p::actor::RpcMultiResponse { response, .. } = item;
                out.push(SerializedBytes::from(UnsafeBytes::from(response)).try_into()?);
            }

            Ok(out)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_get_meta(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: actor::GetMetaOptions,
    ) -> AIngleP2pHandlerResult<Vec<MetadataSet>> {
        let space = dna_hash.into_kitsune();
        let from_agent = from_agent.into_kitsune();
<<<<<<< HEAD
        let basis = dgd_hash.to_kitsune();
        let r_options: event::GetMetaOptions = (&options).into();

        let payload = crate::wire::WireMessage::get_meta(dgd_hash, r_options).encode()?;
=======
        let basis = dht_hash.to_kitsune();
        let r_options: event::GetMetaOptions = (&options).into();

        let payload = crate::wire::WireMessage::get_meta(dht_hash, r_options).encode()?;
>>>>>>> master

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            let result = kitsune_p2p
                .rpc_multi(kitsune_p2p::actor::RpcMulti {
                    space,
                    from_agent,
                    basis,
                    remote_agent_count: options.remote_agent_count,
                    timeout_ms: options.timeout_ms,
                    as_race: options.as_race,
                    race_timeout_ms: options.race_timeout_ms,
                    payload,
                })
                .await?;

            let mut out = Vec::new();
            for item in result {
                let kitsune_p2p::actor::RpcMultiResponse { response, .. } = item;
                out.push(SerializedBytes::from(UnsafeBytes::from(response)).try_into()?);
            }

            Ok(out)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_get_links(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
        link_key: WireLinkMetaKey,
        options: actor::GetLinksOptions,
    ) -> AIngleP2pHandlerResult<Vec<GetLinksResponse>> {
        let space = dna_hash.into_kitsune();
        let from_agent = from_agent.into_kitsune();
        let basis = link_key.basis().to_kitsune();
        let r_options: event::GetLinksOptions = (&options).into();

        let payload = crate::wire::WireMessage::get_links(link_key, r_options).encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            // TODO - We're just targeting a single remote node for now
            //        without doing any pagination / etc...
            //        Setting up RpcMulti to act like RpcSingle
            let result = kitsune_p2p
                .rpc_multi(kitsune_p2p::actor::RpcMulti {
                    space,
                    from_agent,
                    basis,
                    remote_agent_count: Some(1),
                    timeout_ms: options.timeout_ms,
                    as_race: false,
                    race_timeout_ms: options.timeout_ms,
                    payload,
                })
                .await?;

            let mut out = Vec::new();
            for item in result {
                let kitsune_p2p::actor::RpcMultiResponse { response, .. } = item;
                out.push(SerializedBytes::from(UnsafeBytes::from(response)).try_into()?);
            }

            Ok(out)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_get_agent_activity(
        &mut self,
        dna_hash: DnaHash,
        from_agent: AgentPubKey,
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: actor::GetActivityOptions,
    ) -> AIngleP2pHandlerResult<Vec<AgentActivityResponse>> {
        let space = dna_hash.into_kitsune();
        let from_agent = from_agent.into_kitsune();
<<<<<<< HEAD
        // Convert the agent key to an any dgd hash so it can be used
        // as the basis for sending this request
        let agent_hash: AnyDgdHash = agent.clone().into();
=======
        // Convert the agent key to an any dht hash so it can be used
        // as the basis for sending this request
        let agent_hash: AnyDhtHash = agent.clone().into();
>>>>>>> master
        let basis = agent_hash.to_kitsune();
        let r_options: event::GetActivityOptions = (&options).into();

        let payload =
            crate::wire::WireMessage::get_agent_activity(agent, query, r_options).encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            // TODO - We're just targeting a single remote node for now
            //        without doing any pagination / etc...
            //        Setting up RpcMulti to act like RpcSingle
            let result = kitsune_p2p
                .rpc_multi(kitsune_p2p::actor::RpcMulti {
                    space,
                    from_agent,
                    basis,
                    remote_agent_count: Some(1),
                    timeout_ms: options.timeout_ms,
                    as_race: false,
                    race_timeout_ms: options.timeout_ms,
                    payload,
                })
                .await?;

            let mut out = Vec::new();
            for item in result {
                let kitsune_p2p::actor::RpcMultiResponse { response, .. } = item;
                out.push(SerializedBytes::from(UnsafeBytes::from(response)).try_into()?);
            }

            Ok(out)
        }
        .boxed()
        .into())
    }

    #[tracing::instrument(skip(self), level = "trace")]
    fn handle_send_validation_receipt(
        &mut self,
        dna_hash: DnaHash,
        to_agent: AgentPubKey,
        from_agent: AgentPubKey,
        receipt: SerializedBytes,
    ) -> AIngleP2pHandlerResult<()> {
        let space = dna_hash.into_kitsune();
        let to_agent = to_agent.into_kitsune();
        let from_agent = from_agent.into_kitsune();

        let req = crate::wire::WireMessage::validation_receipt(receipt).encode()?;

        let kitsune_p2p = self.kitsune_p2p.clone();
        Ok(async move {
            kitsune_p2p
                .rpc_single(space, to_agent, from_agent, req, None)
                .await?;
            Ok(())
        }
        .boxed()
        .into())
    }
}