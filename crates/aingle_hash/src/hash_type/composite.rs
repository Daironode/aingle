use super::*;
use crate::error::AIngleHashError;
use std::convert::TryInto;

#[cfg(all(test, feature = "serialized-bytes"))]
<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;

/// The AnyDgd (composite) HashType
=======
use aingle_serialized_bytes::prelude::*;

/// The AnyDht (composite) HashType
>>>>>>> master
#[derive(
    Debug, Copy, Clone, Hash, PartialEq, Eq, PartialOrd, Ord, serde::Deserialize, serde::Serialize,
)]
#[cfg_attr(all(test, feature = "serialized-bytes"), derive(SerializedBytes))]
<<<<<<< HEAD
#[serde(from = "AnyDgdSerial", into = "AnyDgdSerial")]
pub enum AnyDgd {
=======
#[serde(from = "AnyDhtSerial", into = "AnyDhtSerial")]
pub enum AnyDht {
>>>>>>> master
    /// The hash of an Entry
    Entry,
    /// The hash of a Header
    Header,
}

<<<<<<< HEAD
impl HashType for AnyDgd {
    fn get_prefix(self) -> &'static [u8] {
        match self {
            AnyDgd::Entry => Entry::new().get_prefix(),
            AnyDgd::Header => Header::new().get_prefix(),
=======
impl HashType for AnyDht {
    fn get_prefix(self) -> &'static [u8] {
        match self {
            AnyDht::Entry => Entry::new().get_prefix(),
            AnyDht::Header => Header::new().get_prefix(),
>>>>>>> master
        }
    }

    fn try_from_prefix(prefix: &[u8]) -> AIngleHashResult<Self> {
        match prefix {
<<<<<<< HEAD
            primitive::ENTRY_PREFIX => Ok(AnyDgd::Entry),
            primitive::HEADER_PREFIX => Ok(AnyDgd::Header),
            _ => Err(AIngleHashError::BadPrefix(
                "AnyDgd".to_string(),
=======
            primitive::ENTRY_PREFIX => Ok(AnyDht::Entry),
            primitive::HEADER_PREFIX => Ok(AnyDht::Header),
            _ => Err(AIngleHashError::BadPrefix(
                "AnyDht".to_string(),
>>>>>>> master
                prefix.try_into().expect("3 byte prefix"),
            )),
        }
    }

    fn hash_name(self) -> &'static str {
<<<<<<< HEAD
        "AnyDgdHash"
    }
}

impl HashTypeAsync for AnyDgd {}

#[derive(serde::Deserialize, serde::Serialize)]
enum AnyDgdSerial {
=======
        "AnyDhtHash"
    }
}

impl HashTypeAsync for AnyDht {}

#[derive(serde::Deserialize, serde::Serialize)]
enum AnyDhtSerial {
>>>>>>> master
    /// The hash of an Entry of EntryType::Agent
    Header(Header),
    /// The hash of any other EntryType
    Entry(Entry),
}

<<<<<<< HEAD
impl From<AnyDgd> for AnyDgdSerial {
    fn from(t: AnyDgd) -> Self {
        match t {
            AnyDgd::Header => AnyDgdSerial::Header(Header),
            AnyDgd::Entry => AnyDgdSerial::Entry(Entry),
=======
impl From<AnyDht> for AnyDhtSerial {
    fn from(t: AnyDht) -> Self {
        match t {
            AnyDht::Header => AnyDhtSerial::Header(Header),
            AnyDht::Entry => AnyDhtSerial::Entry(Entry),
>>>>>>> master
        }
    }
}

<<<<<<< HEAD
impl From<AnyDgdSerial> for AnyDgd {
    fn from(t: AnyDgdSerial) -> Self {
        match t {
            AnyDgdSerial::Header(_) => AnyDgd::Header,
            AnyDgdSerial::Entry(_) => AnyDgd::Entry,
=======
impl From<AnyDhtSerial> for AnyDht {
    fn from(t: AnyDhtSerial) -> Self {
        match t {
            AnyDhtSerial::Header(_) => AnyDht::Header,
            AnyDhtSerial::Entry(_) => AnyDht::Entry,
>>>>>>> master
        }
    }
}
