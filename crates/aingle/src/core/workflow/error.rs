// Error types are self-explanatory
#![allow(missing_docs)]

use super::app_validation_workflow::AppValidationError;
<<<<<<< HEAD
use super::produce_dgd_ops_workflow::dgd_op_light::error::DgdOpConvertError;
=======
use super::produce_dht_ops_workflow::dht_op_light::error::DhtOpConvertError;
>>>>>>> master
use crate::conductor::api::error::ConductorApiError;
use crate::conductor::CellError;
use crate::core::queue_consumer::QueueTriggerClosedError;
use crate::core::ribosome::error::RibosomeError;
use crate::core::SysValidationError;
use aingle_cascade::error::CascadeError;
use aingle_lmdb::error::DatabaseError;
use aingle_p2p::AIngleP2pError;
use aingle_state::source_chain::SourceChainError;
use aingle_state::workspace::WorkspaceError;
use aingle_types::prelude::*;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum WorkflowError {
    #[error(transparent)]
    AppValidationError(#[from] AppValidationError),

    #[error("Agent is invalid: {0:?}")]
    AgentInvalid(AgentPubKey),

    #[error("Conductor API error: {0}")]
    ConductorApi(#[from] Box<ConductorApiError>),

    #[error(transparent)]
    CascadeError(#[from] CascadeError),

    #[error("Workspace error: {0}")]
    WorkspaceError(#[from] WorkspaceError),

    #[error("Database error: {0}")]
    DatabaseError(#[from] DatabaseError),

    #[error(transparent)]
    RibosomeError(#[from] RibosomeError),

    #[error("Source chain error: {0}")]
    SourceChainError(#[from] SourceChainError),

    #[error("Capability token missing")]
    CapabilityMissing,

    #[error(transparent)]
    SerializedBytesError(#[from] SerializedBytesError),

    #[error(transparent)]
<<<<<<< HEAD
    DgdOpConvertError(#[from] DgdOpConvertError),
=======
    DhtOpConvertError(#[from] DhtOpConvertError),
>>>>>>> master

    #[error(transparent)]
    CellError(#[from] CellError),

    #[error(transparent)]
    JoinError(#[from] tokio::task::JoinError),

    #[error(transparent)]
    QueueTriggerClosedError(#[from] QueueTriggerClosedError),

    #[error(transparent)]
    AIngleP2pError(#[from] AIngleP2pError),

    #[error(transparent)]
<<<<<<< HEAD
    DgdOpError(#[from] DgdOpError),
=======
    DhtOpError(#[from] DhtOpError),
>>>>>>> master

    #[error(transparent)]
    SysValidationError(#[from] SysValidationError),
}

/// Internal type to handle running workflows
pub type WorkflowResult<T> = Result<T, WorkflowError>;
