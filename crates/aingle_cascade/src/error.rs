<<<<<<< HEAD
use aingle_hash::{AnyDgdHash, HeaderHash};
use aingle_lmdb::error::DatabaseError;
use aingle_p2p::AIngleP2pError;
use aingle_middleware_bytes::SerializedBytesError;
=======
use aingle_hash::{AnyDhtHash, HeaderHash};
use aingle_lmdb::error::DatabaseError;
use aingle_p2p::AIngleP2pError;
use aingle_serialized_bytes::SerializedBytesError;
>>>>>>> master
use aingle_state::source_chain::SourceChainError;
use aingle_types::prelude::*;
use aingle_zome_types::header::conversions::WrongHeaderError;
// use aingle::conductor::CellError;
<<<<<<< HEAD
// use aingle::core::workflow::produce_dgd_ops_workflow::dgd_op_light::error::DgdOpConvertError;
=======
// use aingle::core::workflow::produce_dht_ops_workflow::dht_op_light::error::DhtOpConvertError;
>>>>>>> master
use thiserror::Error;
use tokio::task::JoinError;

#[derive(Error, Debug)]
pub enum CascadeError {
    #[error(transparent)]
    DatabaseError(#[from] DatabaseError),

    #[error(transparent)]
    ElementGroupError(#[from] ElementGroupError),

    #[error(transparent)]
    HeaderError(#[from] HeaderError),

    #[error("Expected this Header to contain an Entry: {0}")]
    EntryMissing(HeaderHash),

    #[error(transparent)]
<<<<<<< HEAD
    DgdOpError(#[from] DgdOpError),

    #[error("Got an invalid response from an authority for the request hash: {0:?}")]
    InvalidResponse(AnyDgdHash),
=======
    DhtOpError(#[from] DhtOpError),

    #[error("Got an invalid response from an authority for the request hash: {0:?}")]
    InvalidResponse(AnyDhtHash),
>>>>>>> master

    #[error(transparent)]
    JoinError(#[from] JoinError),

    #[error(transparent)]
    SourceChainError(#[from] SourceChainError),

    #[error(transparent)]
    NetworkError(#[from] AIngleP2pError),

    #[error(transparent)]
    SerializedBytesError(#[from] SerializedBytesError),

    #[error(transparent)]
    WrongHeaderError(#[from] WrongHeaderError),

    #[error("Cell is an authority for is missing or incorrect: {0}")]
    AuthorityDataError(#[from] AuthorityDataError),
}

pub type CascadeResult<T> = Result<T, CascadeError>;

#[derive(Error, Debug)]
pub enum AuthorityDataError {
    // #[error(transparent)]
<<<<<<< HEAD
    // DgdOpConvertError(#[from] DgdOpConvertError),
=======
    // DhtOpConvertError(#[from] DhtOpConvertError),
>>>>>>> master
    #[error(transparent)]
    WrongHeaderError(#[from] WrongHeaderError),
    #[error(transparent)]
    HeaderError(#[from] HeaderError),
    #[error("Missing element data: {0:?}")]
    MissingData(String),
    #[error("Missing metadata: {0:?}")]
    MissingMetadata(String),
}

impl AuthorityDataError {
    pub fn missing_data<T: std::fmt::Debug>(data: T) -> CascadeError {
        Self::MissingData(format!("Missing header {:?}", data)).into()
    }
    pub fn missing_data_entry<T: std::fmt::Debug>(data: T) -> CascadeError {
        Self::MissingData(format!("Missing entry for header {:?}", data)).into()
    }
    pub fn missing_metadata<T: std::fmt::Debug>(data: T) -> CascadeError {
        Self::MissingMetadata(format!("{:?}", data)).into()
    }
}
