//! AIngle DnaError type.

use aingle_hash::DnaHash;
use aingle_zome_types::zome::ZomeName;
use thiserror::Error;

/// AIngle DnaError type.
#[derive(Clone, Debug, Error)]
pub enum DnaError {
    /// ZomeNotFound
    #[error("Zome not found: {0}")]
    ZomeNotFound(String),

    /// EmptyZome
    #[error("Zome has no code: {0}")]
    EmptyZome(String),

    /// Invalid
    #[error("DNA is invalid: {0}")]
    Invalid(String),

    /// TraitNotFound
    #[error("Trait not found: {0}")]
    TraitNotFound(String),

    /// ZomeFunctionNotFound
    #[error("Zome function not found: {0}")]
    ZomeFunctionNotFound(String),

    /// SerializedBytesError
    #[error("SerializedBytesError: {0}")]
<<<<<<< HEAD
    SerializedBytesError(#[from] aingle_middleware_bytes::SerializedBytesError),
=======
    SerializedBytesError(#[from] aingle_serialized_bytes::SerializedBytesError),
>>>>>>> master

    /// std::io::Error
    /// we don't #[from] the std::io::Error directly because it doesn't implement Clone
    #[error("std::io::Error: {0}")]
    StdIoError(String),

    /// InvalidWasmHash
    #[error("InvalidWasmHash")]
    InvalidWasmHash,

    /// NonWasmZome
    #[error("Accessed a zome expecting to find a WasmZome, but found other type. Zome name: {0}")]
    NonWasmZome(ZomeName),

    /// DnaHashMismatch
    #[error("DNA hash of file does not match contents.\nHash in file: {0}\nActual hash: {1}")]
    DnaHashMismatch(DnaHash, DnaHash),
}

impl From<std::io::Error> for DnaError {
    fn from(error: std::io::Error) -> Self {
        Self::StdIoError(error.to_string())
    }
}

/// Result type for DnaError
pub type DnaResult<T> = Result<T, DnaError>;