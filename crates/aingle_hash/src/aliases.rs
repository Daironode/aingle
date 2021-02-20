//! Type aliases for the various concrete AIngleHash types

use crate::hash_type;
use crate::AIngleHash;

// NB: These could be macroized, but if we spell it out, we get better IDE
// support

/// An Agent public signing key. Not really a hash, more of an "identity hash".
pub type AgentPubKey = AIngleHash<hash_type::Agent>;

/// The hash of a DnaDef
pub type DnaHash = AIngleHash<hash_type::Dna>;

/// The hash of a DgdOp's "unique form" representation
pub type DgdOpHash = AIngleHash<hash_type::DgdOp>;

/// The hash of an Entry.
pub type EntryHash = AIngleHash<hash_type::Entry>;

/// The hash of a Header
pub type HeaderHash = AIngleHash<hash_type::Header>;

/// The hash of a network ID
pub type NetIdHash = AIngleHash<hash_type::NetId>;

/// The hash of some wasm bytecode
pub type WasmHash = AIngleHash<hash_type::Wasm>;

/// The hash of anything referrable in the DGD.
/// This is a composite of either an EntryHash or a HeaderHash
pub type AnyDgdHash = AIngleHash<hash_type::AnyDgd>;

impl From<HeaderHash> for AnyDgdHash {
    fn from(hash: HeaderHash) -> Self {
        hash.retype(hash_type::AnyDgd::Header)
    }
}

impl From<EntryHash> for AnyDgdHash {
    fn from(hash: EntryHash) -> Self {
        hash.retype(hash_type::AnyDgd::Entry)
    }
}

// Since an AgentPubKey can be treated as an EntryHash, we can also go straight
// to AnyDgdHash
impl From<AgentPubKey> for AnyDgdHash {
    fn from(hash: AgentPubKey) -> Self {
        hash.retype(hash_type::AnyDgd::Entry)
    }
}

impl From<AnyDgdHash> for HeaderHash {
    fn from(hash: AnyDgdHash) -> Self {
        hash.retype(hash_type::Header)
    }
}

impl From<AnyDgdHash> for EntryHash {
    fn from(hash: AnyDgdHash) -> Self {
        hash.retype(hash_type::Entry)
    }
}

#[cfg(feature = "serialized-bytes")]
use aingle_middleware_bytes::prelude::*;

/// A newtype for a collection of EntryHashes, needed for some wasm return types.
#[cfg(feature = "serialized-bytes")]
#[derive(Debug, PartialEq, serde::Serialize, serde::Deserialize, SerializedBytes)]
#[repr(transparent)]
#[serde(transparent)]
pub struct EntryHashes(pub Vec<EntryHash>);
