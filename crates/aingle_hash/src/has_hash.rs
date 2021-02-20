//! Definition of the HasHash trait

use crate::HashType;
use crate::AIngleHash;

/// Anything which has an owned AIngleHashOf.
pub trait HasHash<T: HashType> {
    /// Get the hash by reference
    fn as_hash(&self) -> &AIngleHash<T>;

    /// Convert to the owned hash
    fn into_hash(self) -> AIngleHash<T>;
}
