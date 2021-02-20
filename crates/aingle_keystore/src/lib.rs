#![deny(missing_docs)]
#![allow(clippy::needless_doctest_main)]
//! A Keystore is a secure repository of private keys. KeystoreSender is a
//! reference to a Keystore. KeystoreSender allows async generation of keypairs,
//! and usage of those keypairs, reference by the public AgentPubKey.
//!
//! # Example
//!
//! ```
//! use aingle_hash::AgentPubKey;
//! use aingle_keystore::*;
<<<<<<< HEAD
//! use aingle_middleware_bytes::prelude::*;
=======
//! use aingle_serialized_bytes::prelude::*;
>>>>>>> master
//!
//! #[tokio::main(threaded_scheduler)]
//! async fn main() {
//!     tokio::task::spawn(async move {
//!         let keystore = test_keystore::spawn_test_keystore().await.unwrap();
//!         let agent_pubkey = AgentPubKey::new_from_pure_entropy(&keystore).await.unwrap();
//!
//!         #[derive(Debug, serde::Serialize, serde::Deserialize, SerializedBytes)]
//!         struct MyData(Vec<u8>);
//!
//!         let my_data_1 = MyData(b"signature test data 1".to_vec());
//!
//!         let signature = agent_pubkey.sign(&keystore, &my_data_1).await.unwrap();
//!
//!         /*
//!         assert!(agent_pubkey.verify_signature(&signature, &my_data_1).await.unwrap());
//!         */
//!     }).await.unwrap();
//! }
//! ```

<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

mod error;
pub use error::*;

pub mod keystore_actor;
pub use keystore_actor::KeystoreSender;
pub use keystore_actor::KeystoreSenderExt;
use keystore_actor::*;

mod agent_pubkey_ext;
pub use agent_pubkey_ext::*;

pub mod lair_keystore;
pub mod test_keystore;
