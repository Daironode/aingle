//! Module containing the AIngleP2p actor definition.
#![allow(clippy::too_many_arguments)]

use crate::*;
use aingle_types::activity::AgentActivityResponse;

/// Request a validation package.
#[derive(Clone, Debug)]
pub struct GetValidationPackage {
    /// The dna_hash / space_hash context.
    pub dna_hash: DnaHash,
    /// The agent_id / agent_pub_key context.
    pub agent_pub_key: AgentPubKey,
    /// Request the package from this agent.
    pub request_from: AgentPubKey,
    /// Request the package for this Header
    pub header_hash: HeaderHash,
}

#[derive(Clone, Debug)]
/// Get options help control how the get is processed at various levels.
/// Fields tagged with `[Network]` are network-level controls.
/// Fields tagged with `[Remote]` are controls that will be forwarded to the
/// remote agent processing this `Get` request.
pub struct GetOptions {
    /// [Network]
    /// How many remote nodes should we make requests of / aggregate.
    /// Set to `None` for a default "best-effort".
    pub remote_agent_count: Option<u8>,

    /// [Network]
    /// Timeout to await responses for aggregation.
    /// Set to `None` for a default "best-effort".
    /// Note - if all requests time-out you will receive an empty result,
    /// not a timeout error.
    pub timeout_ms: Option<u64>,

    /// [Network]
    /// We are interested in speed. If `true` and we have any results
    /// when `race_timeout_ms` is expired, those results will be returned.
    /// After `race_timeout_ms` and before `timeout_ms` the first result
    /// received will be returned.
    pub as_race: bool,

    /// [Network]
    /// See `as_race` for details.
    /// Set to `None` for a default "best-effort" race.
    pub race_timeout_ms: Option<u64>,

    /// [Remote]
    /// Whether the remote-end should follow redirects or just return the
    /// requested entry.
    pub follow_redirects: bool,

    /// [Remote]
    /// Return all live headers even if there is deletes.
    /// Useful for metadata calls.
    pub all_live_headers_with_metadata: bool,
}

impl Default for GetOptions {
    fn default() -> Self {
        Self {
            remote_agent_count: None,
            timeout_ms: None,
            as_race: true,
            race_timeout_ms: None,
            follow_redirects: true,
            all_live_headers_with_metadata: false,
        }
    }
}

impl From<aingle_zome_types::entry::GetOptions> for GetOptions {
    fn from(_: aingle_zome_types::entry::GetOptions) -> Self {
        Self::default()
    }
}

<<<<<<< HEAD
/// Get metadata from the DGD.
=======
/// Get metadata from the DHT.
>>>>>>> master
/// Fields tagged with `[Network]` are network-level controls.
/// Fields tagged with `[Remote]` are controls that will be forwarded to the
/// remote agent processing this `GetLinks` request.
#[derive(Clone, Debug)]
pub struct GetMetaOptions {
    /// [Network]
    /// How many remote nodes should we make requests of / aggregate.
    /// Set to `None` for a default "best-effort".
    pub remote_agent_count: Option<u8>,

    /// [Network]
    /// Timeout to await responses for aggregation.
    /// Set to `None` for a default "best-effort".
    /// Note - if all requests time-out you will receive an empty result,
    /// not a timeout error.
    pub timeout_ms: Option<u64>,

    /// [Network]
    /// We are interested in speed. If `true` and we have any results
    /// when `race_timeout_ms` is expired, those results will be returned.
    /// After `race_timeout_ms` and before `timeout_ms` the first result
    /// received will be returned.
    pub as_race: bool,

    /// [Network]
    /// See `as_race` for details.
    /// Set to `None` for a default "best-effort" race.
    pub race_timeout_ms: Option<u64>,

    /// [Remote]
    /// Tells the remote-end which metadata to return
    pub metadata_request: MetadataRequest,
}

impl Default for GetMetaOptions {
    fn default() -> Self {
        Self {
            remote_agent_count: None,
            timeout_ms: None,
            as_race: true,
            race_timeout_ms: None,
            metadata_request: MetadataRequest::default(),
        }
    }
}

#[derive(Debug, Clone)]
<<<<<<< HEAD
/// Get links from the DGD.
=======
/// Get links from the DHT.
>>>>>>> master
/// Fields tagged with `[Network]` are network-level controls.
/// Fields tagged with `[Remote]` are controls that will be forwarded to the
/// remote agent processing this `GetLinks` request.
pub struct GetLinksOptions {
    /// [Network]
    /// Timeout to await responses for aggregation.
    /// Set to `None` for a default "best-effort".
    /// Note - if all requests time-out you will receive an empty result,
    /// not a timeout error.
    pub timeout_ms: Option<u64>,
}

impl Default for GetLinksOptions {
    fn default() -> Self {
        Self { timeout_ms: None }
    }
}

#[derive(Debug, Clone)]
<<<<<<< HEAD
/// Get agent activity from the DGD.
=======
/// Get agent activity from the DHT.
>>>>>>> master
/// Fields tagged with `[Network]` are network-level controls.
/// Fields tagged with `[Remote]` are controls that will be forwarded to the
/// remote agent processing this `GetLinks` request.
pub struct GetActivityOptions {
    /// [Network]
    /// Timeout to await responses for aggregation.
    /// Set to `None` for a default "best-effort".
    /// Note - if all requests time-out you will receive an empty result,
    /// not a timeout error.
    pub timeout_ms: Option<u64>,
    /// Number of times to retry getting elements in parallel.
<<<<<<< HEAD
    /// For a small dgd a large parallel get can overwhelm a single
=======
    /// For a small dht a large parallel get can overwhelm a single
>>>>>>> master
    /// agent and it can be worth retrying the elements that didn't
    /// get found.
    pub retry_gets: u8,
    /// [Remote]
    /// Include the all valid activity headers in the response.
    /// If this is false the call becomes a lightweight response with
    /// just the chain status and highest observed header.
    /// This is useful when you want to ask an authority about the
    /// status of a chain but do not need all the headers.
    pub include_valid_activity: bool,
    /// Include any rejected headers in the response.
    pub include_rejected_activity: bool,
    /// Include the full signed headers and hashes in the response
    /// instead of just the hashes.
    pub include_full_headers: bool,
}

impl Default for GetActivityOptions {
    fn default() -> Self {
        Self {
            timeout_ms: None,
            retry_gets: 0,
            include_valid_activity: true,
            include_rejected_activity: false,
            include_full_headers: false,
        }
    }
}

ghost_actor::ghost_chan! {
    /// The AIngleP2pSender struct allows controlling the AIngleP2p
    /// actor instance.
    pub chan AIngleP2p<AIngleP2pError> {
        /// The p2p module must be informed at runtime which dna/agent pairs it should be tracking.
        fn join(dna_hash: DnaHash, agent_pub_key: AgentPubKey) -> ();

        /// If a cell is deactivated, we'll need to \"leave\" the network module as well.
        fn leave(dna_hash: DnaHash, agent_pub_key: AgentPubKey) -> ();

        /// Invoke a zome function on a remote node (if you have been granted the capability).
        fn call_remote(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            to_agent: AgentPubKey,
            zome_name: ZomeName,
            fn_name: FunctionName,
            cap: Option<CapSecret>,
            payload: ExternIO,
        ) -> SerializedBytes;

        /// Publish data to the correct neighborhood.
        fn publish(
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
        ) -> ();

        /// Request a validation package.
        fn get_validation_package(input: GetValidationPackage) -> ValidationPackageResponse;

<<<<<<< HEAD
        /// Get an entry from the DGD.
        fn get(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            dgd_hash: aingle_hash::AnyDgdHash,
            options: GetOptions,
        ) -> Vec<GetElementResponse>;

        /// Get metadata from the DGD.
        fn get_meta(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            dgd_hash: aingle_hash::AnyDgdHash,
            options: GetMetaOptions,
        ) -> Vec<MetadataSet>;

        /// Get links from the DGD.
=======
        /// Get an entry from the DHT.
        fn get(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            dht_hash: aingle_hash::AnyDhtHash,
            options: GetOptions,
        ) -> Vec<GetElementResponse>;

        /// Get metadata from the DHT.
        fn get_meta(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            dht_hash: aingle_hash::AnyDhtHash,
            options: GetMetaOptions,
        ) -> Vec<MetadataSet>;

        /// Get links from the DHT.
>>>>>>> master
        fn get_links(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            link_key: WireLinkMetaKey,
            options: GetLinksOptions,
        ) -> Vec<GetLinksResponse>;

<<<<<<< HEAD
        /// Get agent activity from the DGD.
=======
        /// Get agent activity from the DHT.
>>>>>>> master
        fn get_agent_activity(
            dna_hash: DnaHash,
            from_agent: AgentPubKey,
            agent: AgentPubKey,
            query: ChainQueryFilter,
            options: GetActivityOptions,
        ) -> Vec<AgentActivityResponse>;

        /// Send a validation receipt to a remote node.
        fn send_validation_receipt(dna_hash: DnaHash, to_agent: AgentPubKey, from_agent: AgentPubKey, receipt: SerializedBytes) -> ();
    }
}

/// Convenience type for referring to the AIngleP2p GhostSender
pub type AIngleP2pRef = ghost_actor::GhostSender<AIngleP2p>;

/// Extension trait for converting GhostSender<AIngleP2p> into AIngleP2pCell
pub trait AIngleP2pRefToCell {
    /// Partially apply dna_hash && agent_pub_key to this sender,
    /// binding it to a specific cell context.
    fn into_cell(self, dna_hash: DnaHash, from_agent: AgentPubKey) -> crate::AIngleP2pCell;

    /// Clone and partially apply dna_hash && agent_pub_key to this sender,
    /// binding it to a specific cell context.
    fn to_cell(&self, dna_hash: DnaHash, from_agent: AgentPubKey) -> crate::AIngleP2pCell;
}

impl AIngleP2pRefToCell for AIngleP2pRef {
    fn into_cell(self, dna_hash: DnaHash, from_agent: AgentPubKey) -> crate::AIngleP2pCell {
        crate::AIngleP2pCell {
            sender: self,
            dna_hash: Arc::new(dna_hash),
            from_agent: Arc::new(from_agent),
        }
    }

    fn to_cell(&self, dna_hash: DnaHash, from_agent: AgentPubKey) -> crate::AIngleP2pCell {
        self.clone().into_cell(dna_hash, from_agent)
    }
}
