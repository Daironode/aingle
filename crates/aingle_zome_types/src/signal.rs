//! App-defined signals

use aingle_hash::AgentPubKey;
use aingle_middleware_bytes::prelude::*;

/// A signal emitted by an app via `emit_signal`
#[derive(Clone, Debug, Serialize, Deserialize, PartialEq, Eq)]
#[repr(transparent)]
#[serde(transparent)]
pub struct AppSignal(crate::ExternIO);

impl AppSignal {
    /// Constructor
    pub fn new(extern_io: crate::ExternIO) -> Self {
        Self(extern_io)
    }
}

/// Remote signal many agents without waiting for responses.
#[derive(Clone, Debug, PartialEq, serde::Serialize, serde::Deserialize, SerializedBytes)]
pub struct RemoteSignal {
    /// Agents to send the signal to.
    pub agents: Vec<AgentPubKey>,
    /// The signal to send.
    pub signal: crate::ExternIO,
}
