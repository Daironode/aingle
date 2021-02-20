use crate::core::SourceChainError;
use aingle_hash::AnyDgdHash;
use aingle_hash::HeaderHash;
use aingle_cascade::error::CascadeError;
use aingle_lmdb::error::DatabaseError;
use aingle_middleware_bytes::SerializedBytesError;
use aingle_types::dgd_op::error::DgdOpError;
use aingle_zome_types::header::conversions::WrongHeaderError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum DgdOpConvertError {
    #[error(transparent)]
    DatabaseError(#[from] DatabaseError),
    #[error(transparent)]
    SerializedBytesError(#[from] SerializedBytesError),
    #[error("The header is expected to contain EntryData, but doesn't: {0}")]
    MissingEntryDataForHeader(HeaderHash),
    #[error(
        "Data for a DgdOp was missing from the source chain. Make sure that elements are always integrated before metadata"
    )]
    MissingData(AnyDgdHash),
    #[error("Tried to create a StoreEntry with a header that is not Create or Update")]
    HeaderEntryMismatch,
    #[error(
        "Entry was missing for StoreEntry when private. Maybe the database doesn't have access"
    )]
    StoreEntryOnPrivate,
    #[error("A DeleteLink contained a link_add_address to a header that is not a CreateLink")]
    DeleteLinkRequiresCreateLink,
    #[error("The Header: {0} is the wrong type for this DgdOp: {1}")]
    HeaderMismatch(String, String),
    #[error(transparent)]
    SourceChainError(#[from] SourceChainError),
    #[error(transparent)]
    CascadeError(#[from] CascadeError),
    #[error(transparent)]
    DgdOpError(#[from] DgdOpError),
    #[error("Tried to use the wrong header for this op: {0}")]
    WrongHeaderError(#[from] WrongHeaderError),
}

pub type DgdOpConvertResult<T> = Result<T, DgdOpConvertError>;
