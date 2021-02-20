use crate::entry::Entry;
use crate::header::CreateLink;
use crate::header::DeleteLink;
use crate::zome_io::ExternIO;
use crate::CallbackResult;
<<<<<<< HEAD
use aingle_hash::AnyDgdHash;
use aingle_middleware_bytes::prelude::*;
=======
use aingle_hash::AnyDhtHash;
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

#[derive(Serialize, Deserialize, SerializedBytes, Debug)]
pub struct ValidateCreateLinkData {
    pub link_add: CreateLink,
    pub base: Entry,
    pub target: Entry,
}

#[derive(Serialize, Deserialize, SerializedBytes, Debug)]
pub struct ValidateDeleteLinkData {
    pub delete_link: DeleteLink,
}

#[derive(Clone, Debug, PartialEq, Serialize, Deserialize, SerializedBytes)]
pub enum ValidateLinkCallbackResult {
    Valid,
    Invalid(String),
<<<<<<< HEAD
    UnresolvedDependencies(Vec<AnyDgdHash>),
=======
    UnresolvedDependencies(Vec<AnyDhtHash>),
>>>>>>> master
}

impl CallbackResult for ValidateLinkCallbackResult {
    fn is_definitive(&self) -> bool {
        matches!(self, ValidateLinkCallbackResult::Invalid(_))
    }
}

impl From<ExternIO> for ValidateLinkCallbackResult {
    fn from(guest_output: ExternIO) -> Self {
        match guest_output.decode() {
            Ok(v) => v,
            Err(e) => Self::Invalid(format!("{:?}", e)),
        }
    }
}
