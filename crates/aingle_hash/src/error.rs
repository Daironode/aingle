//! AIngleHash Error Type.

use crate::AINGLE_HASH_PREFIX_LEN;

/// AIngleHash Error Type.
#[derive(thiserror::Error, Debug, Clone, PartialEq)]
pub enum AIngleHashError {
    /// aingle hashes begin with a lower case u (base64url_no_pad)
    #[error("AIngle Hash missing 'u' prefix")]
    NoU,

    /// could not base64 decode the aingle hash
    #[error("AIngle Hash has invalid base64 encoding")]
    BadBase64,

    /// this string is not the right size for a aingle hash
    #[error("AIngle Hash has incorrect size")]
    BadSize,

    /// this hash does not match a known aingle hash prefix
    #[error("AIngle Hash {0} has unknown prefix {1:?}")]
    BadPrefix(String, [u8; AINGLE_HASH_PREFIX_LEN]),

    /// checksum validation failed
    #[error("AIngle Hash checksum validation failed")]
    BadChecksum,
}

/// AIngleHash Result type
pub type AIngleHashResult<T> = Result<T, AIngleHashError>;
