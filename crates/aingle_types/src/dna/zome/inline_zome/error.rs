#![allow(missing_docs)]

<<<<<<< HEAD
use aingle_middleware_bytes::SerializedBytesError;
=======
use aingle_serialized_bytes::SerializedBytesError;
>>>>>>> master
use aingle_zome_types::prelude::*;
use thiserror::Error;

pub type InlineZomeResult<T> = Result<T, InlineZomeError>;

#[derive(Error, Debug)]
pub enum InlineZomeError {
    #[error("No such InlineZome callback: {0}")]
    NoSuchCallback(FunctionName),

    #[error("Error during host fn call: {0}")]
    HostFnApiError(#[from] HostFnApiError),

    #[error(transparent)]
    SerializationError(#[from] SerializedBytesError),
}