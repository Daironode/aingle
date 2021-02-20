#![allow(missing_docs)]

<<<<<<< HEAD
use crate::encode::aingle_dgd_location_bytes;
use crate::hash_type;
use crate::AgentPubKey;
use crate::AnyDgdHash;
use crate::DgdOpHash;
=======
use crate::encode::aingle_dht_location_bytes;
use crate::hash_type;
use crate::AgentPubKey;
use crate::AnyDhtHash;
use crate::DhtOpHash;
>>>>>>> master
use crate::DnaHash;
use crate::EntryHash;
use crate::HeaderHash;
use crate::NetIdHash;
use crate::WasmHash;
use ::fixt::prelude::*;
use std::convert::TryFrom;

pub type HashTypeEntry = hash_type::Entry;
<<<<<<< HEAD
pub type HashTypeAnyDgd = hash_type::AnyDgd;
=======
pub type HashTypeAnyDht = hash_type::AnyDht;
>>>>>>> master

// TODO: use strum to do this:
//
// fixturator!(
//     HashTypeEntry;
//     unit variants [ Agent Content ] empty Content;
// );

fixturator!(
<<<<<<< HEAD
    HashTypeAnyDgd;
    curve Empty HashTypeAnyDgd::Header;
    curve Unpredictable HashTypeAnyDgd::Header;
    curve Predictable HashTypeAnyDgd::Header;
=======
    HashTypeAnyDht;
    curve Empty HashTypeAnyDht::Header;
    curve Unpredictable HashTypeAnyDht::Header;
    curve Predictable HashTypeAnyDht::Header;
>>>>>>> master
);

/// A type alias for a Vec<u8> whose fixturator is expected to only return
/// a Vec of length 36
pub type ThirtySixHashBytes = Vec<u8>;

// Simply generate "bytes" which is a Vec<u8> of 36 bytes
fixturator!(
    ThirtySixHashBytes,
    append_location([0; 32].to_vec()),
    {
        let mut u8_fixturator = U8Fixturator::new(Unpredictable);
        let mut bytes = vec![];
        for _ in 0..32 {
            bytes.push(u8_fixturator.next().unwrap());
        }
        append_location(bytes)
    },
    {
        let mut index = get_fixt_index!();
        let mut u8_fixturator = U8Fixturator::new_indexed(Predictable, index);
        let mut bytes = vec![];
        for _ in 0..32 {
            bytes.push(u8_fixturator.next().unwrap());
        }
        index += 1;
        set_fixt_index!(index);
        append_location(bytes)
    }
);

fn append_location(mut base: Vec<u8>) -> Vec<u8> {
<<<<<<< HEAD
    let mut loc_bytes = aingle_dgd_location_bytes(&base);
=======
    let mut loc_bytes = aingle_dht_location_bytes(&base);
>>>>>>> master
    base.append(&mut loc_bytes);
    base
}

fixturator!(
    AgentPubKey;
    curve Empty AgentPubKey::from_raw_36(ThirtySixHashBytesFixturator::new_indexed(Empty, get_fixt_index!()).next().unwrap());
    curve Unpredictable AgentPubKey::from_raw_36(ThirtySixHashBytesFixturator::new_indexed(Unpredictable, get_fixt_index!()).next().unwrap());
    curve Predictable {
        // these agent keys match what the mock keystore spits out for the first two agents
        // don't mess with this unless you also update the keystore!!!
        let agents = vec![
            AgentPubKey::try_from("uhCAkmrkoAHPVf_eufG7eC5fm6QKrW5pPMoktvG5LOC0SnJ4vV1Uv")
            .unwrap(),
            AgentPubKey::try_from("uhCAke1j8Z2a-_min0h0pGuEMcYlo_V1l1mt9OtBuywKmHlg4L_R-")
                .unwrap(),
        ];
        agents[get_fixt_index!() % agents.len()].clone()
    };
);

fixturator!(
    EntryHash;
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
    DnaHash;
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
<<<<<<< HEAD
    DgdOpHash;
=======
    DhtOpHash;
>>>>>>> master
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
    HeaderHash;
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
    NetIdHash;
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
    WasmHash;
    constructor fn from_raw_36(ThirtySixHashBytes);
);

fixturator!(
<<<<<<< HEAD
    AnyDgdHash;
    constructor fn from_raw_36_and_type(ThirtySixHashBytes, HashTypeAnyDgd);
=======
    AnyDhtHash;
    constructor fn from_raw_36_and_type(ThirtySixHashBytes, HashTypeAnyDht);
>>>>>>> master
);
