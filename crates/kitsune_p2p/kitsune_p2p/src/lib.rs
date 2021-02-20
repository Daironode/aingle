#![deny(missing_docs)]
<<<<<<< HEAD
//! P2p / dgd communication framework.
=======
//! P2p / dht communication framework.
>>>>>>> master

/// re-exported dependencies
pub mod dependencies {
    pub use ::kitsune_p2p_proxy;
    pub use ::kitsune_p2p_types;
    pub use ::url2;
}

mod types;
pub use types::*;

mod config;
pub use config::*;

mod spawn;
pub use spawn::*;

#[cfg(test)]
pub mod test_util;

#[cfg(test)]
mod test;

pub mod fixt;