<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

#[derive(Clone, Debug, Serialize, Deserialize, PartialEq)]
pub struct GetAgentActivityInput {
    pub agent_pubkey: aingle_hash::AgentPubKey,
    pub chain_query_filter: crate::query::ChainQueryFilter,
    pub activity_request: crate::query::ActivityRequest,
}

impl GetAgentActivityInput {
    /// Constructor.
    pub fn new(
        agent_pubkey: aingle_hash::AgentPubKey,
        chain_query_filter: crate::query::ChainQueryFilter,
        activity_request: crate::query::ActivityRequest,
    ) -> Self {
        Self {
            agent_pubkey,
            chain_query_filter,
            activity_request,
        }
    }
}
