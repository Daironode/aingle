use aingle_middleware_bytes::SerializedBytesError;
use aingle_zome_types::header::conversions::WrongHeaderError;
use aingle_zome_types::Header;
use thiserror::Error;

#[derive(PartialEq, Eq, Clone, Debug, Error)]
pub enum DgdOpError {
    #[error("Tried to create a DgdOp from a Element that requires an Entry. Header type {0:?}")]
    HeaderWithoutEntry(Header),
    #[error(transparent)]
    SerializedBytesError(#[from] SerializedBytesError),
    #[error(transparent)]
    WrongHeaderError(#[from] WrongHeaderError),
}

pub type DgdOpResult<T> = Result<T, DgdOpError>;
