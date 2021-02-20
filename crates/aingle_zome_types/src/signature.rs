//! Signature for authenticity of data
use aingle_hash::AgentPubKey;
<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

/// Input structure for creating a signature.
#[derive(Debug, PartialEq, Serialize, Deserialize, SerializedBytes, Clone)]
pub struct Sign {
    /// The public key associated with the private key that should be used to
    /// generate the signature.
    pub key: aingle_hash::AgentPubKey,

    /// The data that should be signed.
    #[serde(with = "serde_bytes")]
    pub data: Vec<u8>,
}

impl Sign {
    /// construct a new Sign struct.
    pub fn new<S>(key: aingle_hash::AgentPubKey, input: S) -> Result<Self, SerializedBytesError>
    where
        S: Serialize + std::fmt::Debug,
    {
        Ok(Self {
            key,
<<<<<<< HEAD
            data: aingle_middleware_bytes::encode(&input)?,
=======
            data: aingle_serialized_bytes::encode(&input)?,
>>>>>>> master
        })
    }

    /// construct a new Sign struct from raw bytes.
    pub fn new_raw(key: aingle_hash::AgentPubKey, data: Vec<u8>) -> Self {
        Self { key, data }
    }

    /// key getter
    pub fn key(&self) -> &AgentPubKey {
        &self.key
    }

    /// data getter
    pub fn data(&self) -> &[u8] {
        &self.data
    }
}

/// The raw bytes of a signature.
#[derive(Clone, Serialize, Deserialize, SerializedBytes, Hash, PartialEq, Eq, PartialOrd, Ord)]
pub struct Signature(#[serde(with = "serde_bytes")] pub Vec<u8>);

impl From<Vec<u8>> for Signature {
    fn from(bytes: Vec<u8>) -> Self {
        Self(bytes)
    }
}

impl AsRef<[u8]> for Signature {
    fn as_ref(&self) -> &[u8] {
        &self.0
    }
}

impl std::fmt::Debug for Signature {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_fmt(format_args!("Signature(0x"))?;
        for byte in &self.0 {
            f.write_fmt(format_args!("{:02x}", byte))?;
        }
        f.write_fmt(format_args!(")"))?;
        Ok(())
    }
}

/// Mirror struct for Sign that includes a signature to verify against a key and data.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize, SerializedBytes)]
pub struct VerifySignature {
    /// The public key associated with the private key that should be used to
    /// verify the signature.
    pub key: aingle_hash::AgentPubKey,

    /// The signature being verified.
    pub signature: Signature,

    /// The signed data
    #[serde(with = "serde_bytes")]
    pub data: Vec<u8>,
}

impl AsRef<Signature> for VerifySignature {
    fn as_ref(&self) -> &Signature {
        &self.signature
    }
}

impl AsRef<aingle_hash::AgentPubKey> for VerifySignature {
    fn as_ref(&self) -> &AgentPubKey {
        &self.key
    }
}

impl VerifySignature {
    /// Alias for as_ref for data.
    pub fn as_data_ref(&self) -> &[u8] {
        &self.data.as_ref()
    }

    /// Alias for as_ref for signature.
    pub fn as_signature_ref(&self) -> &Signature {
        &self.as_ref()
    }

    /// Alias for as_ref for agent key.
    pub fn as_key_ref(&self) -> &aingle_hash::AgentPubKey {
        &self.as_ref()
    }

    /// construct a new VerifySignature struct.
    pub fn new<D>(
        key: aingle_hash::AgentPubKey,
        signature: Signature,
        data: D,
    ) -> Result<Self, SerializedBytesError>
    where
        D: serde::Serialize + std::fmt::Debug,
    {
        Ok(Self {
            key,
            signature,
<<<<<<< HEAD
            data: aingle_middleware_bytes::encode(&data)?,
=======
            data: aingle_serialized_bytes::encode(&data)?,
>>>>>>> master
        })
    }

    /// construct a new Sign struct from raw bytes.
    pub fn new_raw(key: aingle_hash::AgentPubKey, signature: Signature, data: Vec<u8>) -> Self {
        Self {
            key,
            signature,
            data,
        }
    }
}