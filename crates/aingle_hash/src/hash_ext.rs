use crate::assert_length;
use crate::encode;
use crate::hash_type;
use crate::HashType;
use crate::HashableContent;
use crate::HashableContentBytes;
use crate::AIngleHash;
use crate::AIngleHashOf;
use crate::AIngleHashed;
use crate::AINGLE_HASH_CORE_LEN;
use hash_type::HashTypeAsync;
use hash_type::HashTypeSync;

/// The maximum size to hash synchronously. Anything larger than this will
/// take too long to hash within a single tokio context
pub const MAX_HASHABLE_CONTENT_LEN: usize = 16 * 1000 * 1000; // 16 MiB

impl<T: HashTypeSync> AIngleHash<T> {
    /// Synchronously hash a reference to the given content to produce a AIngleHash
    /// If the content is larger than MAX_HASHABLE_CONTENT_LEN, this will **panic**!
    pub fn with_data_sync<C: HashableContent<HashType = T>>(content: &C) -> AIngleHash<T> {
        hash_from_content(content)
    }
}

impl<T, C> AIngleHashed<C>
where
    T: HashTypeSync,
    C: HashableContent<HashType = T>,
{
    /// Compute the hash of this content and store it alongside
    pub fn from_content_sync(content: C) -> Self {
        let hash: AIngleHashOf<C> = AIngleHash::<T>::with_data_sync(&content);
        Self { content, hash }
    }

    /// Verify that the cached hash matches the content.
    /// Important to run this after e.g. deserialization.
    pub fn verify_hash_sync(&self) -> Result<(), AIngleHash<T>> {
        let hash = AIngleHash::<T>::with_data_sync(&self.content);
        if self.hash == hash {
            Ok(())
        } else {
            Err(hash)
        }
    }
}

impl<T: HashTypeAsync> AIngleHash<T> {
    /// Asynchronously hash a reference to the given content to produce a AIngleHash
    // TODO: this needs to be pushed onto a background thread if the content is large
    pub async fn with_data<C: HashableContent<HashType = T>>(content: &C) -> AIngleHash<T> {
        hash_from_content(content)
    }
}

impl<T, C> AIngleHashed<C>
where
    T: HashTypeAsync,
    C: HashableContent<HashType = T>,
{
    /// Compute the hash of this content and store it alongside
    pub async fn from_content(content: C) -> Self {
        let hash: AIngleHashOf<C> = AIngleHash::<T>::with_data(&content).await;
        Self { content, hash }
    }

    /// Verify that the cached hash matches the content.
    /// Important to run this after e.g. deserialization.
    pub async fn verify_hash(&self) -> Result<(), AIngleHash<T>> {
        let hash = AIngleHash::<T>::with_data(&self.content).await;
        if self.hash == hash {
            Ok(())
        } else {
            Err(hash)
        }
    }
}

fn hash_from_content<T: HashType, C: HashableContent<HashType = T>>(content: &C) -> AIngleHash<T> {
    match content.hashable_content() {
        HashableContentBytes::Content(sb) => {
<<<<<<< HEAD
            let bytes: Vec<u8> = aingle_middleware_bytes::UnsafeBytes::from(sb).into();
=======
            let bytes: Vec<u8> = aingle_serialized_bytes::UnsafeBytes::from(sb).into();
>>>>>>> master
            let hash = encode::blake2b_256(&bytes);
            assert_length!(AINGLE_HASH_CORE_LEN, &hash);
            AIngleHash::<T>::from_raw_32_and_type(hash, content.hash_type())
        }
        HashableContentBytes::Prehashed39(bytes) => AIngleHash::from_raw_39_panicky(bytes),
    }
}
