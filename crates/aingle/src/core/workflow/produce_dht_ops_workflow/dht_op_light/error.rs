use crate::core::SourceChainError;
use aingle_hash::AnyDhtHash;
use aingle_hash::HeaderHash;
use aingle_cascade::error::CascadeError;
use aingle_lmdb::error::DatabaseError;
use aingle_serialized_bytes::SerializedBytesError;
use aingle_types::dht_op::error::DhtOpError;
use aingle_zome_types::header::conversions::WrongHeaderError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum DhtOpConvertError {
    #[error(transparent)]
    DatabaseError(#[from] DatabaseError),
    #[error(transparent)]
    SerializedBytesError(#[from] SerializedBytesError),
    #[error("The header is expected to contain EntryData, but doesn't: {0}")]
    MissingEntryDataForHeader(HeaderHash),
    #[error(
        "Data for a DhtOp was missing from the source chain. Make sure that elements are always integrated before metadata"
    )]
    MissingData(AnyDhtHash),
    #[error("Tried to create a StoreEntry with a header that is not Create or Update")]
    HeaderEntryMismatch,
    #[error(
        "Entry was missing for StoreEntry when private. Maybe the database doesn't have access"
    )]
    StoreEntryOnPrivate,
    #[error("A DeleteLink contained a link_add_address to a header that is not a CreateLink")]
    DeleteLinkRequiresCreateLink,
    #[error("The Header: {0} is the wrong type for this DhtOp: {1}")]
    HeaderMismatch(String, String),
    #[error(transparent)]
    SourceChainError(#[from] SourceChainError),
    #[error(transparent)]
    CascadeError(#[from] CascadeError),
    #[error(transparent)]
    DhtOpError(#[from] DhtOpError),
    #[error("Tried to use the wrong header for this op: {0}")]
    WrongHeaderError(#[from] WrongHeaderError),
}

pub type DhtOpConvertResult<T> = Result<T, DhtOpConvertError>;
