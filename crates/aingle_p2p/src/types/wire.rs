use crate::*;
use aingle_zome_types::zome::FunctionName;

#[derive(Debug, serde::Serialize, serde::Deserialize, SerializedBytes)]
<<<<<<< HEAD
pub(crate) struct WireDgdOpData {
    pub from_agent: aingle_hash::AgentPubKey,
    pub dgd_hash: aingle_hash::AnyDgdHash,
    pub op_data: aingle_types::dgd_op::DgdOp,
}

impl WireDgdOpData {
=======
pub(crate) struct WireDhtOpData {
    pub from_agent: aingle_hash::AgentPubKey,
    pub dht_hash: aingle_hash::AnyDhtHash,
    pub op_data: aingle_types::dht_op::DhtOp,
}

impl WireDhtOpData {
>>>>>>> master
    pub fn encode(self) -> Result<Vec<u8>, SerializedBytesError> {
        Ok(UnsafeBytes::from(SerializedBytes::try_from(self)?).into())
    }

    pub fn decode(data: Vec<u8>) -> Result<Self, SerializedBytesError> {
        let request: SerializedBytes = UnsafeBytes::from(data).into();
        Ok(request.try_into()?)
    }
}

#[derive(Debug, serde::Serialize, serde::Deserialize, SerializedBytes)]
#[serde(tag = "type", content = "content")]
pub(crate) enum WireMessage {
    CallRemote {
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        #[serde(with = "serde_bytes")]
        data: Vec<u8>,
    },
    Publish {
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
>>>>>>> master
    },
    ValidationReceipt {
        #[serde(with = "serde_bytes")]
        receipt: Vec<u8>,
    },
    Get {
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        options: event::GetOptions,
    },
    GetMeta {
        dgd_hash: aingle_hash::AnyDgdHash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        options: event::GetOptions,
    },
    GetMeta {
        dht_hash: aingle_hash::AnyDhtHash,
>>>>>>> master
        options: event::GetMetaOptions,
    },
    GetLinks {
        link_key: WireLinkMetaKey,
        options: event::GetLinksOptions,
    },
    GetAgentActivity {
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: event::GetActivityOptions,
    },
    GetValidationPackage {
        header_hash: HeaderHash,
    },
}

impl WireMessage {
    pub fn encode(&self) -> Result<Vec<u8>, SerializedBytesError> {
<<<<<<< HEAD
        aingle_middleware_bytes::encode(&self)
    }

    pub fn decode(data: &[u8]) -> Result<Self, SerializedBytesError> {
        aingle_middleware_bytes::decode(&data)
=======
        aingle_serialized_bytes::encode(&self)
    }

    pub fn decode(data: &[u8]) -> Result<Self, SerializedBytesError> {
        aingle_serialized_bytes::decode(&data)
>>>>>>> master
    }

    pub fn call_remote(
        zome_name: ZomeName,
        fn_name: FunctionName,
        cap: Option<CapSecret>,
        payload: ExternIO,
    ) -> WireMessage {
        Self::CallRemote {
            zome_name,
            fn_name,
            cap,
            data: payload.into_vec(),
        }
    }

    pub fn publish(
        request_validation_receipt: bool,
<<<<<<< HEAD
        dgd_hash: aingle_hash::AnyDgdHash,
        ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
    ) -> WireMessage {
        Self::Publish {
            request_validation_receipt,
            dgd_hash,
=======
        dht_hash: aingle_hash::AnyDhtHash,
        ops: Vec<(aingle_hash::DhtOpHash, aingle_types::dht_op::DhtOp)>,
    ) -> WireMessage {
        Self::Publish {
            request_validation_receipt,
            dht_hash,
>>>>>>> master
            ops,
        }
    }

    pub fn validation_receipt(receipt: SerializedBytes) -> WireMessage {
        Self::ValidationReceipt {
            receipt: UnsafeBytes::from(receipt).into(),
        }
    }

<<<<<<< HEAD
    pub fn get(dgd_hash: aingle_hash::AnyDgdHash, options: event::GetOptions) -> WireMessage {
        Self::Get { dgd_hash, options }
    }

    pub fn get_meta(
        dgd_hash: aingle_hash::AnyDgdHash,
        options: event::GetMetaOptions,
    ) -> WireMessage {
        Self::GetMeta { dgd_hash, options }
=======
    pub fn get(dht_hash: aingle_hash::AnyDhtHash, options: event::GetOptions) -> WireMessage {
        Self::Get { dht_hash, options }
    }

    pub fn get_meta(
        dht_hash: aingle_hash::AnyDhtHash,
        options: event::GetMetaOptions,
    ) -> WireMessage {
        Self::GetMeta { dht_hash, options }
>>>>>>> master
    }

    pub fn get_links(link_key: WireLinkMetaKey, options: event::GetLinksOptions) -> WireMessage {
        Self::GetLinks { link_key, options }
    }

    pub fn get_agent_activity(
        agent: AgentPubKey,
        query: ChainQueryFilter,
        options: event::GetActivityOptions,
    ) -> WireMessage {
        Self::GetAgentActivity {
            agent,
            query,
            options,
        }
    }
    pub fn get_validation_package(header_hash: HeaderHash) -> WireMessage {
        Self::GetValidationPackage { header_hash }
    }
}