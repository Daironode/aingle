//! crate::saf::wasm is a module for managing webassembly code
//!  - within the in-memory saf struct
//!  - and serialized to json
use backtrace::Backtrace;
use ai_hash::*;
use aingle_middleware_bytes::prelude::*;
use serde::Deserialize;
use serde::Serialize;
use std::fmt;
use std::hash::Hash;
use std::hash::Hasher;
use std::sync::Arc;
use tracing::*;

/// Represents web assembly code.
#[derive(Serialize, Deserialize, Clone, Eq)]
pub struct SafWasm {
    /// the wasm bytes from a .wasm file
    pub code: Arc<Box<[u8]>>,
}

/// A SafWasm paired with its WasmHash
pub type SafWasmHashed = AiHashed<SafWasm>;

impl HashableContent for SafWasm {
    type HashType = hash_type::Wasm;

    fn hash_type(&self) -> Self::HashType {
        hash_type::Wasm
    }

    fn hashable_content(&self) -> HashableContentBytes {
        HashableContentBytes::Content(
            self.try_into()
                .expect("Could not serialize HashableContent"),
        )
    }
}

impl TryFrom<&SafWasm> for SerializedBytes {
    type Error = SerializedBytesError;
    fn try_from(saf_wasm: &SafWasm) -> Result<Self, Self::Error> {
        Ok(SerializedBytes::from(UnsafeBytes::from(
            saf_wasm.code.to_vec(),
        )))
    }
}
impl TryFrom<SafWasm> for SerializedBytes {
    type Error = SerializedBytesError;
    fn try_from(saf_wasm: SafWasm) -> Result<Self, Self::Error> {
        Self::try_from(&saf_wasm)
    }
}

impl TryFrom<SerializedBytes> for SafWasm {
    type Error = SerializedBytesError;
    fn try_from(serialized_bytes: SerializedBytes) -> Result<Self, Self::Error> {
        Ok(SafWasm {
            code: Arc::new(serialized_bytes.bytes().to_owned().into_boxed_slice()),
        })
    }
}

impl SafWasm {
    /// Provide basic placeholder for wasm entries in saf structs, used for testing only.
    pub fn new_invalid() -> Self {
        debug!(
            "SafWasm::new_invalid() called from:\n{:?}",
            Backtrace::new()
        );
        SafWasm {
            code: Arc::new(Box::new([])),
        }
    }

    /// get a new Arc to the Vec<u8> bytes for the wasm
    pub fn code(&self) -> Arc<Box<[u8]>> {
        Arc::clone(&self.code)
    }
}

impl fmt::Debug for SafWasm {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "<<<SAF WASM CODE>>>")
    }
}

impl PartialEq for SafWasm {
    fn eq(&self, other: &SafWasm) -> bool {
        self.code == other.code
    }
}

impl Hash for SafWasm {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.code.hash(state);
    }
}

impl From<Vec<u8>> for SafWasm {
    fn from(wasm: Vec<u8>) -> Self {
        Self {
            code: Arc::new(wasm.into_boxed_slice()),
        }
    }
}
