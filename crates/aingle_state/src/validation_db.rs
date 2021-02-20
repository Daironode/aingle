//! # Validation Database Types

use aingle_hash::AgentPubKey;
<<<<<<< HEAD
use aingle_hash::AnyDgdHash;
use aingle_hash::DgdOpHash;
=======
use aingle_hash::AnyDhtHash;
use aingle_hash::DhtOpHash;
>>>>>>> master
use aingle_lmdb::buffer::KvBufFresh;
use aingle_lmdb::db::VALIDATION_LIMBO;
use aingle_lmdb::error::DatabaseResult;
use aingle_lmdb::prelude::EnvironmentRead;
use aingle_lmdb::prelude::GetDb;
<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master
use aingle_types::prelude::*;
use shrinkwraprs::Shrinkwrap;

#[derive(Shrinkwrap)]
#[shrinkwrap(mutable)]
/// The database for putting ops into to await validation
pub struct ValidationLimboStore(pub KvBufFresh<ValidationLimboKey, ValidationLimboValue>);

/// Key to the validation limbo
<<<<<<< HEAD
pub type ValidationLimboKey = DgdOpHash;
=======
pub type ValidationLimboKey = DhtOpHash;
>>>>>>> master

/// A type for storing in databases that only need the hashes.
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq)]
pub struct ValidationLimboValue {
    /// Status of this op in the limbo
    pub status: ValidationLimboStatus,
    /// The actual op
<<<<<<< HEAD
    pub op: DgdOpLight,
    /// Where the op was sent to
    pub basis: AnyDgdHash,
=======
    pub op: DhtOpLight,
    /// Where the op was sent to
    pub basis: AnyDhtHash,
>>>>>>> master
    /// When the op was added to limbo
    pub time_added: Timestamp,
    /// Last time we tried to validated the op
    pub last_try: Option<Timestamp>,
    /// Number of times we have tried to validate the op
    pub num_tries: u32,
    /// The agent that sent you this op
    pub from_agent: Option<AgentPubKey>,
}

<<<<<<< HEAD
/// The status of a [DgdOp] in limbo
=======
/// The status of a [DhtOp] in limbo
>>>>>>> master
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq)]
pub enum ValidationLimboStatus {
    /// Is awaiting to be system validated
    Pending,
    /// Is waiting for dependencies so the op can proceed to system validation
<<<<<<< HEAD
    AwaitingSysDeps(AnyDgdHash),
    /// Is awaiting to be app validated
    SysValidated,
    /// Is waiting for dependencies so the op can proceed to app validation
    AwaitingAppDeps(Vec<AnyDgdHash>),
=======
    AwaitingSysDeps(AnyDhtHash),
    /// Is awaiting to be app validated
    SysValidated,
    /// Is waiting for dependencies so the op can proceed to app validation
    AwaitingAppDeps(Vec<AnyDhtHash>),
>>>>>>> master
}

impl ValidationLimboStore {
    /// Create a new Validation Limbo db
    pub fn new(env: EnvironmentRead) -> DatabaseResult<Self> {
        let db = env.get_db(&*VALIDATION_LIMBO)?;
        Ok(Self(KvBufFresh::new(env, db)))
    }
}
