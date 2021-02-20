//! Types needed for all validation
use std::convert::TryFrom;

use derivative::Derivative;
use aingle_hash::DgdOpHash;
use aingle_types::dgd_op::DgdOp;

use super::workflow::error::WorkflowResult;
use super::SourceChainError;
use super::SysValidationError;
use super::ValidationOutcome;

/// Exit early with either an outcome or an error
pub enum OutcomeOrError<T, E> {
    Outcome(T),
    Err(E),
}

/// Helper macro for implementing from sub error types
/// for the error in OutcomeOrError
#[macro_export]
macro_rules! from_sub_error {
    ($error_type:ident, $sub_error_type:ident) => {
        impl<T> From<$sub_error_type> for OutcomeOrError<T, $error_type> {
            fn from(e: $sub_error_type) -> Self {
                OutcomeOrError::Err($error_type::from(e))
            }
        }
    };
}

/// Type for deriving ordering of DgdOps
/// Don't change the order of this enum unless
/// you mean to change the order we process ops
#[allow(missing_docs)]
#[derive(Clone, Copy, Debug, Eq, PartialEq, Ord, PartialOrd)]
pub enum DgdOpOrder {
    RegisterAgentActivity(aingle_zome_types::timestamp::Timestamp),
    StoreEntry(aingle_zome_types::timestamp::Timestamp),
    StoreElement(aingle_zome_types::timestamp::Timestamp),
    RegisterUpdatedContent(aingle_zome_types::timestamp::Timestamp),
    RegisterUpdatedElement(aingle_zome_types::timestamp::Timestamp),
    RegisterDeletedBy(aingle_zome_types::timestamp::Timestamp),
    RegisterDeletedEntryHeader(aingle_zome_types::timestamp::Timestamp),
    RegisterAddLink(aingle_zome_types::timestamp::Timestamp),
    RegisterRemoveLink(aingle_zome_types::timestamp::Timestamp),
}

/// Op data that will be ordered by [DgdOpOrder]
#[derive(Derivative, Debug, Clone)]
#[derivative(Eq, PartialEq, Ord, PartialOrd)]
pub struct OrderedOp<V> {
    pub order: DgdOpOrder,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub hash: DgdOpHash,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub op: DgdOp,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub value: V,
}

impl From<&DgdOp> for DgdOpOrder {
    fn from(op: &DgdOp) -> Self {
        use DgdOpOrder::*;
        match op {
            DgdOp::StoreElement(_, h, _) => StoreElement(h.timestamp()),
            DgdOp::StoreEntry(_, h, _) => StoreEntry(*h.timestamp()),
            DgdOp::RegisterAgentActivity(_, h) => RegisterAgentActivity(h.timestamp()),
            DgdOp::RegisterUpdatedContent(_, h, _) => RegisterUpdatedContent(h.timestamp),
            DgdOp::RegisterUpdatedElement(_, h, _) => RegisterUpdatedElement(h.timestamp),
            DgdOp::RegisterDeletedBy(_, h) => RegisterDeletedBy(h.timestamp),
            DgdOp::RegisterDeletedEntryHeader(_, h) => RegisterDeletedEntryHeader(h.timestamp),
            DgdOp::RegisterAddLink(_, h) => RegisterAddLink(h.timestamp),
            DgdOp::RegisterRemoveLink(_, h) => RegisterRemoveLink(h.timestamp),
        }
    }
}

impl OutcomeOrError<ValidationOutcome, SysValidationError> {
    /// Convert an OutcomeOrError<ValidationOutcome, SysValidationError> into
    /// a InvalidCommit and exit the call zome workflow early
    pub fn invalid_call_zome_commit<T>(self) -> WorkflowResult<T> {
        Err(SourceChainError::InvalidCommit(ValidationOutcome::try_from(self)?.to_string()).into())
    }
}
