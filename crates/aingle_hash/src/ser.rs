//! Defines the serialization rules for AIngleHashes

use crate::HashType;
use crate::AIngleHash;
use aingle_middleware_bytes::SerializedBytes;
use aingle_middleware_bytes::SerializedBytesError;
use aingle_middleware_bytes::UnsafeBytes;

impl<T: HashType> serde::Serialize for AIngleHash<T> {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        serializer.serialize_bytes(self.get_raw_39())
    }
}

impl<'de, T: HashType> serde::Deserialize<'de> for AIngleHash<T> {
    fn deserialize<D>(deserializer: D) -> Result<AIngleHash<T>, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        deserializer.deserialize_bytes(AIngleHashVisitor(std::marker::PhantomData))
    }
}

struct AIngleHashVisitor<T: HashType>(std::marker::PhantomData<T>);

impl<'de, T: HashType> serde::de::Visitor<'de> for AIngleHashVisitor<T> {
    type Value = AIngleHash<T>;

    fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
        formatter.write_str("a AIngleHash of primitive hash_type")
    }

    fn visit_bytes<E>(self, h: &[u8]) -> Result<Self::Value, E>
    where
        E: serde::de::Error,
    {
        if !h.len() == 39 {
            Err(serde::de::Error::custom(
                "AIngleHash serialized representation must be exactly 39 bytes",
            ))
        } else {
            AIngleHash::from_raw_39(h.to_vec())
                .map_err(|e| serde::de::Error::custom(format!("AIngleHash error: {:?}", e)))
        }
    }

    fn visit_seq<A>(self, mut seq: A) -> Result<Self::Value, A::Error>
    where
        A: serde::de::SeqAccess<'de>,
    {
        let mut vec = Vec::with_capacity(seq.size_hint().unwrap_or(0));

        while let Some(b) = seq.next_element()? {
            vec.push(b);
        }

        self.visit_bytes(&vec)
    }
}

impl<T: HashType> std::convert::TryFrom<&AIngleHash<T>> for SerializedBytes {
    type Error = SerializedBytesError;
    fn try_from(t: &AIngleHash<T>) -> std::result::Result<SerializedBytes, SerializedBytesError> {
        match aingle_middleware_bytes::encode(t) {
            Ok(v) => Ok(SerializedBytes::from(UnsafeBytes::from(v))),
            Err(e) => Err(SerializedBytesError::Serialize(e.to_string())),
        }
    }
}

impl<T: HashType> std::convert::TryFrom<AIngleHash<T>> for SerializedBytes {
    type Error = SerializedBytesError;
    fn try_from(t: AIngleHash<T>) -> std::result::Result<SerializedBytes, SerializedBytesError> {
        SerializedBytes::try_from(&t)
    }
}

impl<T: HashType> std::convert::TryFrom<SerializedBytes> for AIngleHash<T> {
    type Error = SerializedBytesError;
    fn try_from(sb: SerializedBytes) -> std::result::Result<AIngleHash<T>, SerializedBytesError> {
        match aingle_middleware_bytes::decode(sb.bytes()) {
            Ok(v) => Ok(v),
            Err(e) => Err(SerializedBytesError::Deserialize(e.to_string())),
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::*;
    use aingle_middleware_bytes::prelude::*;
    use std::convert::TryInto;

    #[derive(serde::Deserialize, Debug)]
    #[serde(transparent)]
    struct TestByteArray(#[serde(with = "serde_bytes")] Vec<u8>);

    #[test]
    #[cfg(feature = "serialized-bytes")]
    fn test_serialized_bytes_roundtrip() {
        use aingle_middleware_bytes::SerializedBytes;
        use std::convert::TryInto;

        let h_orig = DnaHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]);
        let h: SerializedBytes = h_orig.clone().try_into().unwrap();
        let h: DnaHash = h.try_into().unwrap();

        assert_eq!(h_orig, h);
        assert_eq!(*h.hash_type(), hash_type::Dna::new());
    }

    #[test]
    fn test_rmp_roundtrip() {
        let h_orig = AgentPubKey::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]);
        let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
        let h: AgentPubKey = aingle_middleware_bytes::decode(&buf).unwrap();

        assert_eq!(h_orig, h);
        assert_eq!(*h.hash_type(), hash_type::Agent::new());

        // Make sure that the representation is a raw 39-byte array
        let array: TestByteArray = aingle_middleware_bytes::decode(&buf).unwrap();
        assert_eq!(array.0.len(), AINGLE_HASH_FULL_LEN);
        assert_eq!(
            array.0,
            vec![
                132, 32, 36, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219,
                219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219, 219,
                219, 219, 219, 219, 219, 219,
            ]
        );
    }

    #[test]
    fn test_composite_hashtype_roundtrips() {
        {
            let h_orig = AnyDgdHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::AnyDgd::Header,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let h: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
            assert_eq!(h_orig, h);
            assert_eq!(*h.hash_type(), hash_type::AnyDgd::Header);
        }
        {
            let h_orig = AnyDgdHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::AnyDgd::Entry,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let h: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
            assert_eq!(h_orig, h);
            assert_eq!(*h.hash_type(), hash_type::AnyDgd::Entry);
        }
        {
            let h_orig = AnyDgdHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::AnyDgd::Entry,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let h: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
            assert_eq!(h_orig, h);
            assert_eq!(*h.hash_type(), hash_type::AnyDgd::Entry);
        }
    }

    #[test]
    fn test_any_dgd_deserialization() {
        {
            let h_orig = EntryHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::Entry,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let _: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
        }
        {
            let h_orig = HeaderHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::Header,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let _: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
        }
    }

    #[test]
    #[should_panic]
    fn test_any_dgd_deserialization_crossover_error() {
        {
            let h_orig = DgdOpHash::from_raw_36_and_type(
                vec![0xdb; AINGLE_HASH_UNTYPED_LEN],
                hash_type::DgdOp,
            );
            let buf = aingle_middleware_bytes::encode(&h_orig).unwrap();
            let _: AnyDgdHash = aingle_middleware_bytes::decode(&buf).unwrap();
        }
    }

    #[test]
    fn test_struct_to_struct_roundtrip() {
        #[derive(Debug, PartialEq, Eq, serde::Serialize, serde::Deserialize, SerializedBytes)]
        struct TestData {
            e: EntryHash,
            h: HeaderHash,
        }

        let orig = TestData {
            e: EntryHash::from_raw_36_and_type(vec![0xdb; AINGLE_HASH_UNTYPED_LEN], hash_type::Entry),
            h: HeaderHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]),
        };

        let sb: SerializedBytes = (&orig).try_into().unwrap();
        let res: TestData = sb.try_into().unwrap();

        assert_eq!(orig, res);
        assert_eq!(*orig.e.hash_type(), hash_type::Entry);
        assert_eq!(*orig.h.hash_type(), hash_type::Header);
    }

    #[test]
    fn test_json_to_rust() {
        #[derive(Debug, PartialEq, Eq, serde::Serialize, serde::Deserialize, SerializedBytes)]
        struct Data {
            any_hash: AnyDgdHash,
            content: String,
        }

        let any_hash = AnyDgdHash::from_raw_36_and_type(
            b"000000000000000000000000000000000000".to_vec(),
            hash_type::AnyDgd::Header,
        );
        let hash_type_sb: SerializedBytes = any_hash.hash_type().try_into().unwrap();
        let hash_type_json = r#"{"Header":[132,41,36]}"#;
        assert_eq!(format!("{:?}", hash_type_sb), hash_type_json.to_string());

        let hash_type_from_sb: hash_type::AnyDgd = hash_type_sb.try_into().unwrap();
        assert_eq!(hash_type_from_sb, hash_type::AnyDgd::Header);

        let hash_type_from_json: hash_type::AnyDgd = serde_json::from_str(&hash_type_json).unwrap();
        assert_eq!(hash_type_from_json, hash_type::AnyDgd::Header);
    }

    #[test]
    fn test_generic_content_roundtrip() {
        #[derive(Debug, Default, PartialEq, Eq, serde::Serialize, serde::Deserialize)]
        struct Generic<K> {
            bytes: Vec<u8>,
            __marker: std::marker::PhantomData<K>,
        }

        impl<K> Generic<K>
        where
            K: serde::Serialize + serde::de::DeserializeOwned + std::fmt::Debug,
            // V: Serialize + DeserializeOwned + std::fmt::Debug,
        {
            fn new() -> Self {
                Self {
                    bytes: Vec::new(),
                    __marker: Default::default(),
                }
            }

            fn get(&self) -> K {
                aingle_middleware_bytes::decode(&self.bytes).unwrap()
            }

            fn put(&mut self, k: &K) {
                self.bytes = aingle_middleware_bytes::encode(k).unwrap();
            }
        }

        let mut g: Generic<HeaderHash> = Generic::new();
        let h = HeaderHash::from_raw_36(vec![0xdb; AINGLE_HASH_UNTYPED_LEN]);
        g.put(&h);
        assert_eq!(h, g.get());
    }
}
