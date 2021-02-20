use crate::header::HeaderHashes;
use crate::zome_io::ExternIO;
use crate::CallbackResult;
<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

#[derive(Clone, PartialEq, Serialize, Deserialize, SerializedBytes, Debug)]
pub enum PostCommitCallbackResult {
    Success,
    Fail(HeaderHashes, String),
}

impl From<ExternIO> for PostCommitCallbackResult {
    fn from(guest_output: ExternIO) -> Self {
        match guest_output.decode() {
            Ok(v) => v,
            Err(e) => Self::Fail(vec![].into(), format!("{:?}", e)),
        }
    }
}

impl CallbackResult for PostCommitCallbackResult {
    fn is_definitive(&self) -> bool {
        matches!(self, PostCommitCallbackResult::Fail(_, _))
    }
}
