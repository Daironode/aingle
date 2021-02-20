#![deny(missing_docs)]
//! aingle specific wrapper around more generic p2p module

use aingle_hash::*;
<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master
use aingle_types::prelude::*;
use std::sync::Arc;

mod types;
pub use types::actor::AIngleP2pRef;
pub use types::actor::AIngleP2pSender;
pub use types::*;

mod spawn;
use ghost_actor::dependencies::tracing;
use ghost_actor::dependencies::tracing_futures::Instrument;
pub use spawn::*;
pub use test::stub_network;
pub use test::AIngleP2pCellFixturator;

pub use kitsune_p2p;

#[mockall::automock]
#[async_trait::async_trait]
/// A wrapper around AIngleP2pSender that partially applies the dna_hash / agent_pub_key.
/// I.e. a sender that is tied to a specific cell.
pub trait AIngleP2pCellT {
    /// owned getter
    fn dna_hash(&self) -> DnaHash;

    /// owned getter
    fn from_agent(&self) -> AgentPubKey;

    /// The p2p module must be informed at runtime which dna/agent pairs it should be tracking.
    async fn join(&mut self) -> actor::AIngleP2pResult<()>;

    /// If a cell is deactivated, we'll need to \"leave\" the network module as well.
    async fn leave(&mut self) -> actor::AIngleP2pResult<()>;

    /// Invoke a zome function on a remote node (if you have been granted the capability).
    async fn call_remote(
        &mut self,
        to_agent: AgentPubKey,
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        payload: ExternIO,
    ) -> actor::AIngleP2pResult<SerializedBytes>;

    /// Publish data to the correct neighborhood.
    #[allow(clippy::ptr_arg)]
    async fn publish(
        &mut self,
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
>>>>>>> master
        timeout_ms: Option<u64>,
    ) -> actor::AIngleP2pResult<()>;

    /// Request a validation package.
    async fn get_validation_package(
        &mut self,
        request_from: AgentPubKey,
        header_hash: HeaderHash,
    ) -> actor::AIngleP2pResult<ValidationPackageResponse>;

<<<<<<< HEAD
    /// Get an entry from the DGD.
    async fn get(
        &mut self,
        dgd_hash: aingle_hash::AnyDgdHash,
        options: actor::GetOptions,
    ) -> actor::AIngleP2pResult<Vec<GetElementResponse>>;

    /// Get metadata from the DGD.
    async fn get_meta(
        &mut self,
        dgd_hash: aingle_hash::AnyDgdHash,
        options: actor::GetMetaOptions,
    ) -> actor::AIngleP2pResult<Vec<MetadataSet>>;

    /// Get links from the DGD.
=======
    /// Get an entry from the DHT.
    async fn get(
        &mut self,
        dht_hash: aingle_hash::AnyDhtHash,
        options: actor::GetOptions,
    ) -> actor::AIngleP2pResult<Vec<GetElementResponse>>;

    /// Get metadata from the DHT.
    async fn get_meta(
        &mut self,
        dht_hash: aingle_hash::AnyDhtHash,
        options: actor::GetMetaOptions,
    ) -> actor::AIngleP2pResult<Vec<MetadataSet>>;

    /// Get links from the DHT.
>>>>>>> master
    async fn get_links(
        &mut self,
        link_key: WireLinkMetaKey,
        options: actor::GetLinksOptions,
    ) -> actor::AIngleP2pResult<Vec<GetLinksResponse>>;

<<<<<<< HEAD
    /// Get agent activity from the DGD.
=======
    /// Get agent activity from the DHT.
>>>>>>> master
    async fn get_agent_activity(
        &mut self,
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: actor::GetActivityOptions,
    ) -> actor::AIngleP2pResult<Vec<AgentActivityResponse>>;

    /// Send a validation receipt to a remote node.
    async fn send_validation_receipt(
        &mut self,
        to_agent: AgentPubKey,
        receipt: SerializedBytes,
    ) -> actor::AIngleP2pResult<()>;
}

/// A wrapper around AIngleP2pSender that partially applies the dna_hash / agent_pub_key.
/// I.e. a sender that is tied to a specific cell.
#[derive(Clone)]
pub struct AIngleP2pCell {
    sender: ghost_actor::GhostSender<actor::AIngleP2p>,
    dna_hash: Arc<DnaHash>,
    from_agent: Arc<AgentPubKey>,
}

#[async_trait::async_trait]
impl AIngleP2pCellT for AIngleP2pCell {
    /// owned getter
    fn dna_hash(&self) -> DnaHash {
        (*self.dna_hash).clone()
    }

    /// owned getter
    fn from_agent(&self) -> AgentPubKey {
        (*self.from_agent).clone()
    }

    /// The p2p module must be informed at runtime which dna/agent pairs it should be tracking.
    async fn join(&mut self) -> actor::AIngleP2pResult<()> {
        self.sender
            .join((*self.dna_hash).clone(), (*self.from_agent).clone())
            .await
    }

    /// If a cell is deactivated, we'll need to \"leave\" the network module as well.
    async fn leave(&mut self) -> actor::AIngleP2pResult<()> {
        self.sender
            .leave((*self.dna_hash).clone(), (*self.from_agent).clone())
            .await
    }

    /// Invoke a zome function on a remote node (if you have been granted the capability).
    async fn call_remote(
        &mut self,
        to_agent: AgentPubKey,
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        payload: ExternIO,
    ) -> actor::AIngleP2pResult<SerializedBytes> {
        self.sender
            .call_remote(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
                to_agent,
                zome_name,
                fn_name,
                cap,
                payload,
            )
            .await
    }

    /// Publish data to the correct neighborhood.
    async fn publish(
        &mut self,
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
>>>>>>> master
        timeout_ms: Option<u64>,
    ) -> actor::AIngleP2pResult<()> {
        self.sender
            .publish(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
                request_validation_receipt,
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                ops,
                timeout_ms,
            )
            .await
    }

    /// Request a validation package.
    async fn get_validation_package(
        &mut self,
        request_from: AgentPubKey,
        header_hash: HeaderHash,
    ) -> actor::AIngleP2pResult<ValidationPackageResponse> {
        self.sender
            .get_validation_package(actor::GetValidationPackage {
                dna_hash: (*self.dna_hash).clone(),
                agent_pub_key: (*self.from_agent).clone(),
                request_from,
                header_hash,
            })
            .await
    }

<<<<<<< HEAD
    /// Get an entry from the DGD.
    async fn get(
        &mut self,
        dgd_hash: aingle_hash::AnyDgdHash,
=======
    /// Get an entry from the DHT.
    async fn get(
        &mut self,
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: actor::GetOptions,
    ) -> actor::AIngleP2pResult<Vec<GetElementResponse>> {
        self.sender
            .get(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                options,
            )
            .instrument(tracing::debug_span!("AIngleP2p::get"))
            .await
    }

<<<<<<< HEAD
    /// Get metadata from the DGD.
    async fn get_meta(
        &mut self,
        dgd_hash: aingle_hash::AnyDgdHash,
=======
    /// Get metadata from the DHT.
    async fn get_meta(
        &mut self,
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: actor::GetMetaOptions,
    ) -> actor::AIngleP2pResult<Vec<MetadataSet>> {
        self.sender
            .get_meta(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
<<<<<<< HEAD
                dgd_hash,
=======
                dht_hash,
>>>>>>> master
                options,
            )
            .await
    }

<<<<<<< HEAD
    /// Get links from the DGD.
=======
    /// Get links from the DHT.
>>>>>>> master
    async fn get_links(
        &mut self,
        link_key: WireLinkMetaKey,
        options: actor::GetLinksOptions,
    ) -> actor::AIngleP2pResult<Vec<GetLinksResponse>> {
        self.sender
            .get_links(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
                link_key,
                options,
            )
            .await
    }

<<<<<<< HEAD
    /// Get agent activity from the DGD.
=======
    /// Get agent activity from the DHT.
>>>>>>> master
    async fn get_agent_activity(
        &mut self,
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: actor::GetActivityOptions,
    ) -> actor::AIngleP2pResult<Vec<AgentActivityResponse>> {
        self.sender
            .get_agent_activity(
                (*self.dna_hash).clone(),
                (*self.from_agent).clone(),
                agent,
                query,
                options,
            )
            .await
    }

    /// Send a validation receipt to a remote node.
    async fn send_validation_receipt(
        &mut self,
        to_agent: AgentPubKey,
        receipt: SerializedBytes,
    ) -> actor::AIngleP2pResult<()> {
        self.sender
            .send_validation_receipt(
                (*self.dna_hash).clone(),
                to_agent,
                (*self.from_agent).clone(),
                receipt,
            )
            .await
    }
}

<<<<<<< HEAD
pub use kitsune_p2p::dgd_arc;
=======
pub use kitsune_p2p::dht_arc;
>>>>>>> master

mod test;