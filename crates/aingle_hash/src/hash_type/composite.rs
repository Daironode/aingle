use super::*;
use crate::error::AIngleHashError;
use std::convert::TryInto;

#[cfg(all(test, feature = "serialized-bytes"))]
use aingle_middleware_bytes::prelude::*;

/// The AnyDgd (composite) HashType
#[derive(
    Debug, Copy, Clone, Hash, PartialEq, Eq, PartialOrd, Ord, serde::Deserialize, serde::Serialize,
)]
#[cfg_attr(all(test, feature = "serialized-bytes"), derive(SerializedBytes))]
#[serde(from = "AnyDgdSerial", into = "AnyDgdSerial")]
pub enum AnyDgd {
    /// The hash of an Entry
    Entry,
    /// The hash of a Header
    Header,
}

impl HashType for AnyDgd {
    fn get_prefix(self) -> &'static [u8] {
        match self {
            AnyDgd::Entry => Entry::new().get_prefix(),
            AnyDgd::Header => Header::new().get_prefix(),
        }
    }

    fn try_from_prefix(prefix: &[u8]) -> AIngleHashResult<Self> {
        match prefix {
            primitive::ENTRY_PREFIX => Ok(AnyDgd::Entry),
            primitive::HEADER_PREFIX => Ok(AnyDgd::Header),
            _ => Err(AIngleHashError::BadPrefix(
                "AnyDgd".to_string(),
                prefix.try_into().expect("3 byte prefix"),
            )),
        }
    }

    fn hash_name(self) -> &'static str {
        "AnyDgdHash"
    }
}

impl HashTypeAsync for AnyDgd {}

#[derive(serde::Deserialize, serde::Serialize)]
enum AnyDgdSerial {
    /// The hash of an Entry of EntryType::Agent
    Header(Header),
    /// The hash of any other EntryType
    Entry(Entry),
}

impl From<AnyDgd> for AnyDgdSerial {
    fn from(t: AnyDgd) -> Self {
        match t {
            AnyDgd::Header => AnyDgdSerial::Header(Header),
            AnyDgd::Entry => AnyDgdSerial::Entry(Entry),
        }
    }
}

impl From<AnyDgdSerial> for AnyDgd {
    fn from(t: AnyDgdSerial) -> Self {
        match t {
            AnyDgdSerial::Header(_) => AnyDgd::Header,
            AnyDgdSerial::Entry(_) => AnyDgd::Entry,
        }
    }
}
