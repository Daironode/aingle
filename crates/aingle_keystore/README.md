# crate

[![Project](https://img.shields.io/badge/project-aingle-blue.svg?style=flat-square)](http://aingle.ai/)
[![Forum](https://img.shields.io/badge/chat-forum%2eaingle%2enet-blue.svg?style=flat-square)](https://forum.aingle.ai)
[![Chat](https://img.shields.io/badge/chat-chat%2eaingle%2enet-blue.svg?style=flat-square)](https://chat.aingle.ai)

[![Twitter Follow](https://img.shields.io/twitter/follow/aingle.svg?style=social&label=Follow)](https://twitter.com/aingle)
License: [![License: CAL 1.0](https://img.shields.io/badge/License-CAL%201.0-blue.svg)](https://github.com/AIngleLab/cryptographic-autonomy-license)

Current version: 0.0.1

A Keystore is a secure repository of private keys. KeystoreSender is a
reference to a Keystore. KeystoreSender allows async generation of keypairs,
and usage of those keypairs, reference by the public AgentPubKey.

## Example

```rust
use ai_hash::AgentPubKey;
use crate::*;
use aingle_middleware_bytes::prelude::*;

#[tokio::main(flavor = "multi_thread")]
async fn main() {
    tokio::task::spawn(async move {
        let _ = aingle_crypto::crypto_init_sodium();

        let keystore = test_keystore::spawn_test_keystore(vec![]).await.unwrap();
        let agent_pubkey = AgentPubKey::new_from_pure_entropy(&keystore).await.unwrap();

        #[derive(Debug, serde::Serialize, serde::Deserialize, SerializedBytes)]
        struct MyData(Vec<u8>);

        let my_data_1 = MyData(b"signature test data 1".to_vec());

        let signature = agent_pubkey.sign(&keystore, &my_data_1).await.unwrap();

        assert!(agent_pubkey.verify_signature(&signature, &my_data_1).await.unwrap());
    }).await.unwrap();
}
```

## Contribute


* Connect with us on our [forum](https://forum.aingle.ai)

## License
 [![License: CAL 1.0](https://img.shields.io/badge/License-CAL-1.0-blue.svg)](https://github.com/AIngleLab/cryptographic-autonomy-license)

Copyright (C) 2019 - 2021, AIngle


