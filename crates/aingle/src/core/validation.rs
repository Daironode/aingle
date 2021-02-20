//! Types needed for all validation
use std::convert::TryFrom;

use derivative::Derivative;
<<<<<<< HEAD
use aingle_hash::DgdOpHash;
use aingle_types::dgd_op::DgdOp;
=======
use aingle_hash::DhtOpHash;
use aingle_types::dht_op::DhtOp;
>>>>>>> master

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

<<<<<<< HEAD
/// Type for deriving ordering of DgdOps
=======
/// Type for deriving ordering of DhtOps
>>>>>>> master
/// Don't change the order of this enum unless
/// you mean to change the order we process ops
#[allow(missing_docs)]
#[derive(Clone, Copy, Debug, Eq, PartialEq, Ord, PartialOrd)]
<<<<<<< HEAD
pub enum DgdOpOrder {
=======
pub enum DhtOpOrder {
>>>>>>> master
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

<<<<<<< HEAD
/// Op data that will be ordered by [DgdOpOrder]
#[derive(Derivative, Debug, Clone)]
#[derivative(Eq, PartialEq, Ord, PartialOrd)]
pub struct OrderedOp<V> {
    pub order: DgdOpOrder,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub hash: DgdOpHash,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub op: DgdOp,
=======
/// Op data that will be ordered by [DhtOpOrder]
#[derive(Derivative, Debug, Clone)]
#[derivative(Eq, PartialEq, Ord, PartialOrd)]
pub struct OrderedOp<V> {
    pub order: DhtOpOrder,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub hash: DhtOpHash,
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub op: DhtOp,
>>>>>>> master
    #[derivative(PartialEq = "ignore", PartialOrd = "ignore", Ord = "ignore")]
    pub value: V,
}

<<<<<<< HEAD
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
=======
impl From<&DhtOp> for DhtOpOrder {
    fn from(op: &DhtOp) -> Self {
        use DhtOpOrder::*;
        match op {
            DhtOp::StoreElement(_, h, _) => StoreElement(h.timestamp()),
            DhtOp::StoreEntry(_, h, _) => StoreEntry(*h.timestamp()),
            DhtOp::RegisterAgentActivity(_, h) => RegisterAgentActivity(h.timestamp()),
            DhtOp::RegisterUpdatedContent(_, h, _) => RegisterUpdatedContent(h.timestamp),
            DhtOp::RegisterUpdatedElement(_, h, _) => RegisterUpdatedElement(h.timestamp),
            DhtOp::RegisterDeletedBy(_, h) => RegisterDeletedBy(h.timestamp),
            DhtOp::RegisterDeletedEntryHeader(_, h) => RegisterDeletedEntryHeader(h.timestamp),
            DhtOp::RegisterAddLink(_, h) => RegisterAddLink(h.timestamp),
            DhtOp::RegisterRemoveLink(_, h) => RegisterRemoveLink(h.timestamp),
>>>>>>> master
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
