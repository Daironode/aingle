//! Definition of the HasHash trait

use crate::HashType;
use crate::AiHash;

/// Anything which has an owned AiHashOf.
pub trait HasHash<T: HashType> {
    /// Get the hash by reference
    fn as_hash(&self) -> &AiHash<T>;

    /// Convert to the owned hash
    fn into_hash(self) -> AiHash<T>;
}
