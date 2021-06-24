# ai_hash

[![Project](https://img.shields.io/badge/project-aingle-blue.svg?style=flat-square)](http://aingle.ai/)
[![Forum](https://img.shields.io/badge/chat-forum%2eaingle%2enet-blue.svg?style=flat-square)](https://forum.aingle.ai)
[![Chat](https://img.shields.io/badge/chat-chat%2eaingle%2enet-blue.svg?style=flat-square)](https://chat.aingle.ai)

[![Twitter Follow](https://img.shields.io/twitter/follow/aingle.svg?style=social&label=Follow)](https://twitter.com/aingle)
License: [![License: CAL 1.0](https://img.shields.io/badge/License-CAL%201.0-blue.svg)](https://github.com/AIngleLab/cryptographic-autonomy-license)

ai_hash::AiHash is a hashing framework for AIngle.

Note that not all AiHashes are simple hashes of the full content as you might expect in a "content-addressable" application. The main exception is `AgentPubKey`, which is simply the key itself to enable self-proving signatures. As an exception it is also named exceptionally, i.e. it doesn't end in "Hash".

## Hash Types

Each AiHash has a HashType. There are two flavors of HashType: *primitive*, and *composite*

### Primitive HashTypes

Each primitive HashType has a unique 3-byte prefix associated with it, to easily distiguish between hashes in any environment. These prefixes are multihash compatible. The primitive types are:

| hash type | AiHash alias | prefix |
|-----------|----------------|--------|
| Agent     | AgentPubKey    | uhCAk  |
| Entry     | EntryHash      | uhCEk  |
| SgdOp     | SgdOpHash      | uhCQk  |
| Saf       | SafHash        | uhC0k  |
| NetId     | NetIdHash      | uhCIk  |
| Header    | HeaderHash     | uhCkk  |
| Wasm      | SafWasmHash    | uhCok  |

The "AiHash alias" column lists the type aliases provided to refer to each type of AiHash. For instance, `HeaderHash` is the following type alias:

```rust
pub type HeaderHash = AiHash<hash_type::Header>;
```

(the prefixes listed are the base64 representations)

### Composite HashTypes

Composite hash types are used in contexts when one of several primitive hash types would be valid. They are implemented as Rust enums. The composite types are:

`EntryHash`: used to hash Entries. An Entry can hash to either a `ContentHash` or an `AgentPubKey`.

`AnySgdHash`: used to hash arbitrary SGD data. SGD data is either a Header or an Entry, therefore AnySgdHash can refer to either a `HeaderHash` or an `EntryHash`.

## Serialization

AiHash implements `Display` providing a `to_string()` function accessing the hash as a user friendly string. It also provides TryFrom for string types allowing you to parse this string representation.

AiHash includes a 4 byte (or u32) sgd "location" that serves dual purposes. - It is used as a checksum when parsing string representations. - It is used as a u32 in our sgd sharding algorithm.

AiHash implements [SerializedBytes](https://lib.rs/crates/aingle_middleware_bytes) to make it easy to cross ffi barriers such as WASM and the UI websocket.

## Example

```rust
use ai_hash::*;
use std::convert::TryInto;
use aingle_middleware_bytes::SerializedBytes;

let entry: EntryHash =
    "uhCEkWCsAgoKkkfwyJAglj30xX_GLLV-3BXuFy436a2SqpcEwyBzm"
    .try_into()
    .unwrap();

assert_eq!(3860645936, entry.get_loc());

let bytes: SerializedBytes = entry.try_into().unwrap();

assert_eq!(
    "{\"type\":\"EntryHash\",\"hash\":[88,43,0,130,130,164,145,252,50,36,8,37,143,125,49,95,241,139,45,95,183,5,123,133,203,141,250,107,100,170,165,193,48,200,28,230]}",
    &format!("{:?}", bytes),
);
```

## Advanced

Calculating hashes takes time - In a futures context we don't want to block. AiHash provides sync (blocking) and async (non-blocking) apis for hashing.

```rust
use ai_hash::*;

let entry_content = b"test entry content";

let content_hash = EntryHash::with_data_sync(entry_content.to_vec()).into();

assert_eq!(
    "EntryHash(uhCEkhPbA5vaw3Fk-ZvPSKuyyjg8eoX98fve75qiUEFgAE3BO7D4d)",
    &format!("{:?}", content_hash),
);
```

### Sometimes your data doesn't want to be re-hashed:

```rust
use ai_hash::*;

// pretend our pub key is all 0xdb bytes
let agent_pub_key = vec![0xdb; 32];

let agent_id: AiHash = AgentPubKey::from_raw_32(agent_pub_key).into();

assert_eq!(
    "AgentPubKey(uhCAk29vb29vb29vb29vb29vb29vb29vb29vb29vb29vb29uTp5Iv)",
    &format!("{:?}", agent_id),
);
```

## Contribute


* Connect with us on our [forum](https://forum.aingle.ai)

## License
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

Copyright (C) 2019 - 2021, AIngle


