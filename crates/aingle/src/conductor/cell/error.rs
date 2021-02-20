use crate::conductor::api::error::ConductorApiError;
use crate::conductor::entry_def_store::error::EntryDefStoreError;
use crate::core::ribosome::error::RibosomeError;
use crate::core::ribosome::guest_callback::init::InitResult;
use crate::core::workflow::error::WorkflowError;
use crate::core::workflow::produce_dgd_ops_workflow::dgd_op_light::error::DgdOpConvertError;
use crate::core::SourceChainError;
use aingle_cascade::error::CascadeError;
use aingle_lmdb::error::DatabaseError;
use aingle_p2p::AIngleP2pError;
use aingle_types::prelude::*;
use aingle_zome_types::cell::CellId;

use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum CellError {
    #[error("error dealing with workspace state: {0}")]
    DatabaseError(#[from] DatabaseError),
    #[error(transparent)]
    CascadeError(#[from] CascadeError),
    #[error("The Dna was not found in the store")]
    DnaMissing,
    #[error("Failed to join the create cell task: {0}")]
    JoinError(#[from] tokio::task::JoinError),
    #[error("Genesis failed: {0}")]
    Genesis(Box<ConductorApiError>),
    #[error(transparent)]
    HeaderError(#[from] HeaderError),
    #[error("This cell has not had a successful genesis and cannot be created")]
    CellWithoutGenesis(CellId),
    #[error(
        "The cell failed to cleanup its environment because: {0}. Recommend manually deleting the database at: {1}"
    )]
    Cleanup(String, PathBuf),
    #[error(transparent)]
    DnaError(#[from] DnaError),
    #[error(transparent)]
    EntryDefStoreError(#[from] EntryDefStoreError),
    #[error(transparent)]
    WorkflowError(#[from] Box<WorkflowError>),
    #[error(transparent)]
    WorkspaceError(#[from] aingle_state::workspace::WorkspaceError),
    #[error(transparent)]
    RibosomeError(#[from] RibosomeError),
    #[error(transparent)]
    SourceChainError(#[from] SourceChainError),
    #[error("The cell tried to run the initialize zomes callback but failed because {0:?}")]
    InitFailed(InitResult),
    #[error(transparent)]
    AIngleP2pError(#[from] AIngleP2pError),
    #[error(transparent)]
    ConductorApiError(#[from] Box<ConductorApiError>),
    #[error(transparent)]
    SerializedBytesError(#[from] aingle_middleware_bytes::SerializedBytesError),
    #[error(transparent)]
    DgdOpConvertError(#[from] DgdOpConvertError),
    #[error("Todo")]
    Todo,
}

pub type CellResult<T> = Result<T, CellError>;
