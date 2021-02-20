//! Defines the AIngleHash type, used for all hashes in AIngle.
//!
//! AIngleHashes come in a variety of types. See the `hash_type::primitive`
//! module for the full list.
//!
//! AIngleHashes are serialized as a plain 39-byte sequence.
//! The structure is like so:
//!
//! ```text
//! PPPCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCLLLL
//! ^  ^                                  ^
//!  \  \---------"untyped"--------------/
//!   \                                 /
//!    \-------------"full"------------/
//!
//! P: 3 byte prefix to indicate hash type
//! C: 32 byte hash, the "core"
<<<<<<< HEAD
//! L: 4 byte hash of the core hash, for DGD location
=======
//! L: 4 byte hash of the core hash, for DHT location
>>>>>>> master
//! ```
//!
//! The 36 bytes which exclude the initial 3-byte type prefix are known
//! throughout the codebase as the "untyped" hash
//!
//! The complete 39 bytes together are known as the "full" hash

use crate::encode;
use crate::error::AIngleHashResult;
use crate::has_hash::HasHash;
use crate::HashType;
use crate::PrimitiveHashType;

/// Length of the prefix bytes (3)
pub const AINGLE_HASH_PREFIX_LEN: usize = 3;

/// Length of the core bytes (32)
pub const AINGLE_HASH_CORE_LEN: usize = 32;

/// Length of the location bytes (4)
pub const AINGLE_HASH_LOC_LEN: usize = 4;

/// Length of the core bytes + the loc bytes (36 = 32 + 4),
/// i.e. everything except the type prefix
pub const AINGLE_HASH_UNTYPED_LEN: usize = AINGLE_HASH_CORE_LEN + AINGLE_HASH_LOC_LEN; // 36

/// Length of the full AIngleHash bytes (39 = 3 + 32 + 4)
pub const AINGLE_HASH_FULL_LEN: usize = AINGLE_HASH_PREFIX_LEN + AINGLE_HASH_CORE_LEN + AINGLE_HASH_LOC_LEN;

/// Helper for ensuring the the proper number of bytes is used in various situations
#[macro_export]
macro_rules! assert_length {
    ($len:expr, $hash:expr) => {
        debug_assert_eq!(
            $hash.len(),
            $len,
            "invalid byte count for AIngleHash {:?}",
            $hash
        );
    };
}

/// A AIngleHash contains a vector of 36 bytes representing a 32-byte blake2b hash
<<<<<<< HEAD
/// plus 4 bytes representing a DGD location. It also contains a zero-sized
=======
/// plus 4 bytes representing a DHT location. It also contains a zero-sized
>>>>>>> master
/// type which specifies what it is a hash of.
///
/// There is custom de/serialization implemented in [ser.rs]
#[derive(Clone, Hash, PartialEq, Eq, PartialOrd, Ord)]
pub struct AIngleHash<T: HashType> {
    hash: Vec<u8>,
    hash_type: T,
}

impl<T: HashType> AIngleHash<T> {
    /// Raw constructor: Create a AIngleHash from 39 bytes, using the prefix
    /// bytes to determine the hash_type
    pub fn from_raw_39(hash: Vec<u8>) -> AIngleHashResult<Self> {
        assert_length!(AINGLE_HASH_FULL_LEN, &hash);
        let hash_type = T::try_from_prefix(&hash[0..3])?;
        Ok(Self { hash, hash_type })
    }
    /// Raw constructor: Create a AIngleHash from 39 bytes, using the prefix
    /// bytes to determine the hash_type. Panics if hash_type does not match.
    pub fn from_raw_39_panicky(hash: Vec<u8>) -> Self {
        Self::from_raw_39(hash).expect("the specified hash_type does not match the prefix bytes")
    }

    /// Use a precomputed hash + location byte array in vec form,
    /// along with a type, to construct a hash. Used in this crate only, for testing.
    pub fn from_raw_36_and_type(mut bytes: Vec<u8>, hash_type: T) -> Self {
        assert_length!(AINGLE_HASH_UNTYPED_LEN, &bytes);
        let mut hash = hash_type.get_prefix().to_vec();
        hash.append(&mut bytes);
        assert_length!(AINGLE_HASH_FULL_LEN, &hash);
        Self { hash, hash_type }
    }

    /// Change the type of this AIngleHash, keeping the same bytes
    pub fn retype<TT: HashType>(mut self, hash_type: TT) -> AIngleHash<TT> {
        let prefix = hash_type.get_prefix();
        self.hash[0..AINGLE_HASH_PREFIX_LEN].copy_from_slice(&prefix[0..AINGLE_HASH_PREFIX_LEN]);
        AIngleHash {
            hash: self.hash,
            hash_type,
        }
    }

    /// The HashType of this hash
    pub fn hash_type(&self) -> &T {
        &self.hash_type
    }

    /// Get the raw 39-byte Vec including the 3 byte prefix, base 32 bytes, and the 4 byte loc
    pub fn get_raw_39(&self) -> &[u8] {
        &self.hash[..]
    }

    /// Get 36-byte Vec which excludes the 3 byte prefix
    pub fn get_raw_36(&self) -> &[u8] {
        let bytes = &self.hash[AINGLE_HASH_PREFIX_LEN..];
        assert_length!(AINGLE_HASH_UNTYPED_LEN, bytes);
        bytes
    }

    /// Fetch just the core 32 bytes (without the 4 location bytes)
    pub fn get_raw_32(&self) -> &[u8] {
        let bytes = &self.hash[AINGLE_HASH_PREFIX_LEN..AINGLE_HASH_PREFIX_LEN + AINGLE_HASH_CORE_LEN];
        assert_length!(AINGLE_HASH_CORE_LEN, bytes);
        bytes
    }

<<<<<<< HEAD
    /// Fetch the aingle DGD location for this hash
=======
    /// Fetch the aingle dht location for this hash
>>>>>>> master
    pub fn get_loc(&self) -> u32 {
        bytes_to_loc(&self.hash[AINGLE_HASH_FULL_LEN - AINGLE_HASH_LOC_LEN..])
    }

    /// consume into the inner byte vector
    pub fn into_inner(self) -> Vec<u8> {
        assert_length!(AINGLE_HASH_FULL_LEN, &self.hash);
        self.hash
    }
}

impl<T: HashType> AIngleHash<T> {
    /// Construct a AIngleHash from a 32-byte hash.
    /// The 3 prefix bytes will be added based on the provided HashType,
    /// and the 4 location bytes will be computed.
    ///
    /// For convenience, 36 bytes can also be passed in, in which case
    /// the location bytes will used as provided, not computed.
    pub fn from_raw_32_and_type(mut hash: Vec<u8>, hash_type: T) -> Self {
        if hash.len() == AINGLE_HASH_CORE_LEN {
<<<<<<< HEAD
            hash.append(&mut encode::aingle_dgd_location_bytes(&hash));
=======
            hash.append(&mut encode::aingle_dht_location_bytes(&hash));
>>>>>>> master
        }

        assert_length!(AINGLE_HASH_UNTYPED_LEN, &hash);

        AIngleHash::from_raw_36_and_type(hash, hash_type)
    }
}

impl<P: PrimitiveHashType> AIngleHash<P> {
    /// Construct from 36 raw bytes, using the known PrimitiveHashType
    pub fn from_raw_36(hash: Vec<u8>) -> Self {
        assert_length!(AINGLE_HASH_UNTYPED_LEN, &hash);
        Self::from_raw_36_and_type(hash, P::new())
    }
    /// Construct a AIngleHash from a prehashed raw 32-byte slice.
    /// The location bytes will be calculated.
    pub fn from_raw_32(hash: Vec<u8>) -> Self {
        Self::from_raw_32_and_type(hash, P::new())
    }
}

impl<T: HashType> AsRef<[u8]> for AIngleHash<T> {
    fn as_ref(&self) -> &[u8] {
        assert_length!(AINGLE_HASH_FULL_LEN, &self.hash);
        &self.hash
    }
}

impl<T: HashType> IntoIterator for AIngleHash<T> {
    type Item = u8;
    type IntoIter = std::vec::IntoIter<Self::Item>;

    fn into_iter(self) -> Self::IntoIter {
        self.hash.into_iter()
    }
}

impl<T: HashType> HasHash<T> for AIngleHash<T> {
    fn as_hash(&self) -> &AIngleHash<T> {
        &self
    }
    fn into_hash(self) -> AIngleHash<T> {
        self
    }
}

// NB: See encode/encode_raw module for Display impl
impl<T: HashType> std::fmt::Debug for AIngleHash<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_fmt(format_args!("{}({})", self.hash_type().hash_name(), self))?;
        Ok(())
    }
}

/// internal convert 4 location bytes into a u32 location
fn bytes_to_loc(bytes: &[u8]) -> u32 {
    (bytes[0] as u32)
        + ((bytes[1] as u32) << 8)
        + ((bytes[2] as u32) << 16)
        + ((bytes[3] as u32) << 24)
}

#[cfg(test)]
mod tests {
    use crate::*;

    #[cfg(not(feature = "string-encoding"))]
    fn assert_type<T: HashType>(t: &str, h: AIngleHash<T>) {
        assert_eq!(3_688_618_971, h.get_loc());
        assert_eq!(
            "[219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219]",
            format!("{:?}", h.get_raw_32()),
        );
    }

    #[test]
    #[cfg(not(feature = "string-encoding"))]
    fn test_enum_types() {
        assert_type(
            "DnaHash",
            DnaHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
        );
        assert_type(
            "NetIdHash",
            NetIdHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
        );
        assert_type(
            "AgentPubKey",
            AgentPubKey::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
        );
        assert_type(
            "EntryHash",
            EntryHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
        );
        assert_type(
<<<<<<< HEAD
            "DgdOpHash",
            DgdOpHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
=======
            "DhtOpHash",
            DhtOpHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
>>>>>>> master
        );
    }

    #[test]
    #[should_panic]
    fn test_fails_with_bad_size() {
        DnaHash::from_raw_36(vec![0xdb; 35]);
    }
}
