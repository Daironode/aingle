use crate::HasHash;
use crate::HashableContent;
use crate::AIngleHashOf;
use aingle_serialized_bytes::prelude::*;

/// Represents some piece of content along with its hash representation, so that
/// hashes need not be calculated multiple times.
/// Provides an easy constructor which consumes the content.
// TODO: consider making lazy with OnceCell
#[derive(Serialize, Deserialize)]
pub struct AIngleHashed<C: HashableContent> {
    pub(crate) content: C,
    pub(crate) hash: AIngleHashOf<C>,
}

impl<C: HashableContent> HasHash<C::HashType> for AIngleHashed<C> {
    fn as_hash(&self) -> &AIngleHashOf<C> {
        &self.hash
    }

    fn into_hash(self) -> AIngleHashOf<C> {
        self.hash
    }
}

impl<C> AIngleHashed<C>
where
    C: HashableContent,
{
    /// Combine content with its precalculated hash
    pub fn with_pre_hashed(content: C, hash: AIngleHashOf<C>) -> Self {
        Self { content, hash }
    }

    // NB: as_hash and into_hash are provided by the HasHash impl

    /// Accessor for content
    pub fn as_content(&self) -> &C {
        &self.content
    }

    /// Convert to content
    pub fn into_content(self) -> C {
        self.content
    }

    /// Deconstruct as a tuple
    pub fn into_inner(self) -> (C, AIngleHashOf<C>) {
        (self.content, self.hash)
    }
}

impl<C> Clone for AIngleHashed<C>
where
    C: HashableContent + Clone,
{
    fn clone(&self) -> Self {
        Self {
            content: self.content.clone(),
            hash: self.hash.clone(),
        }
    }
}

impl<C> std::fmt::Debug for AIngleHashed<C>
where
    C: HashableContent + std::fmt::Debug,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_fmt(format_args!("AIngleHashed({:?})", self.content))?;
        Ok(())
    }
}

impl<C> std::convert::From<AIngleHashed<C>> for (C, AIngleHashOf<C>)
where
    C: HashableContent,
{
    fn from(g: AIngleHashed<C>) -> (C, AIngleHashOf<C>) {
        g.into_inner()
    }
}

impl<C> std::ops::Deref for AIngleHashed<C>
where
    C: HashableContent,
{
    type Target = C;

    fn deref(&self) -> &Self::Target {
        self.as_content()
    }
}

impl<C> std::convert::AsRef<C> for AIngleHashed<C>
where
    C: HashableContent,
{
    fn as_ref(&self) -> &C {
        self.as_content()
    }
}

impl<C> std::borrow::Borrow<C> for AIngleHashed<C>
where
    C: HashableContent,
{
    fn borrow(&self) -> &C {
        self.as_content()
    }
}

impl<C> std::cmp::PartialEq for AIngleHashed<C>
where
    C: HashableContent,
{
    fn eq(&self, other: &Self) -> bool {
        self.hash == other.hash
    }
}

impl<C> std::cmp::Eq for AIngleHashed<C> where C: HashableContent {}

impl<C> std::hash::Hash for AIngleHashed<C>
where
    C: HashableContent,
{
    fn hash<StdH: std::hash::Hasher>(&self, state: &mut StdH) {
        std::hash::Hash::hash(&self.hash, state)
    }
}
