<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

pub const X25519_PUB_KEY_BYTES: usize = 32;

#[derive(Clone, Copy, SerializedBytes)]
pub struct X25519PubKey([u8; X25519_PUB_KEY_BYTES]);

crate::secure_primitive!(X25519PubKey, X25519_PUB_KEY_BYTES);
