use crate::assert_length;
use crate::error::AIngleHashError;
use crate::HashType;
use crate::AIngleHash;
use crate::PrimitiveHashType;
use crate::AINGLE_HASH_CORE_LEN;
use crate::AINGLE_HASH_FULL_LEN;
use crate::AINGLE_HASH_PREFIX_LEN;
use std::convert::TryFrom;
use std::convert::TryInto;

impl<P: PrimitiveHashType> TryFrom<&str> for AIngleHash<P> {
    type Error = AIngleHashError;
    fn try_from(s: &str) -> Result<Self, AIngleHashError> {
        let hash_type = P::new();
        AIngleHash::from_raw_39(aingle_hash_decode(hash_type.get_prefix(), s)?)
    }
}

impl<P: PrimitiveHashType> TryFrom<&String> for AIngleHash<P> {
    type Error = AIngleHashError;
    fn try_from(s: &String) -> Result<Self, AIngleHashError> {
        Self::try_from(s as &str)
    }
}

impl<P: PrimitiveHashType> TryFrom<String> for AIngleHash<P> {
    type Error = AIngleHashError;
    fn try_from(s: String) -> Result<Self, AIngleHashError> {
        Self::try_from(&s)
    }
}

impl<T: HashType> std::fmt::Display for AIngleHash<T> {
    fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
        write!(f, "{}", aingle_hash_encode(self.get_raw_39()))
    }
}

/// internal REPR for aingle hash
pub fn aingle_hash_encode(data: &[u8]) -> String {
    format!("u{}", base64::encode_config(data, base64::URL_SAFE_NO_PAD),)
}

/// internal PARSE for aingle hash REPR
pub fn aingle_hash_decode(prefix: &[u8], s: &str) -> Result<Vec<u8>, AIngleHashError> {
    if &s[..1] != "u" {
        return Err(AIngleHashError::NoU);
    }
    let s = match base64::decode_config(&s[1..], base64::URL_SAFE_NO_PAD) {
        Err(_) => return Err(AIngleHashError::BadBase64),
        Ok(s) => s,
    };
    if s.len() != AINGLE_HASH_FULL_LEN {
        return Err(AIngleHashError::BadSize);
    }
    let actual_prefix: [u8; AINGLE_HASH_PREFIX_LEN] = s[..AINGLE_HASH_PREFIX_LEN].try_into().unwrap();
    if actual_prefix != prefix {
        return Err(AIngleHashError::BadPrefix(
            format!("{:?}", prefix),
            actual_prefix,
        ));
    }
<<<<<<< HEAD
    let loc_bytes = aingle_dgd_location_bytes(
=======
    let loc_bytes = aingle_dht_location_bytes(
>>>>>>> master
        &s[AINGLE_HASH_PREFIX_LEN..AINGLE_HASH_PREFIX_LEN + AINGLE_HASH_CORE_LEN],
    );
    let loc_bytes: &[u8] = &loc_bytes;
    if loc_bytes != &s[AINGLE_HASH_PREFIX_LEN + AINGLE_HASH_CORE_LEN..] {
        return Err(AIngleHashError::BadChecksum);
    }
    assert_length!(AINGLE_HASH_FULL_LEN, &s);
    Ok(s.to_vec())
}

<<<<<<< HEAD
/// internal compute the aingle dgd location u32
pub fn aingle_dgd_location_bytes(data: &[u8]) -> Vec<u8> {
=======
/// internal compute the aingle dht location u32
pub fn aingle_dht_location_bytes(data: &[u8]) -> Vec<u8> {
>>>>>>> master
    // Assert the data size is relatively small so we are
    // comfortable executing this synchronously / blocking tokio thread.
    assert_eq!(32, data.len(), "only 32 byte hashes supported");

    let hash = blake2b_128(data);
    let mut out = vec![hash[0], hash[1], hash[2], hash[3]];
    for i in (4..16).step_by(4) {
        out[0] ^= hash[i];
        out[1] ^= hash[i + 1];
        out[2] ^= hash[i + 2];
        out[3] ^= hash[i + 3];
    }
    out
}

/// internal compute a 32 byte blake2b hash
pub fn blake2b_256(data: &[u8]) -> Vec<u8> {
    let hash = blake2b_simd::Params::new().hash_length(32).hash(data);
    hash.as_bytes().to_vec()
}

/// internal compute a 16 byte blake2b hash
pub fn blake2b_128(data: &[u8]) -> Vec<u8> {
    let hash = blake2b_simd::Params::new().hash_length(16).hash(data);
    hash.as_bytes().to_vec()
}
