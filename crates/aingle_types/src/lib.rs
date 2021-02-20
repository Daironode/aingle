//! Common types used by other AIngle crates.
//!
//! This crate is a complement to the
//! [aingle_zome_types crate](https://crates.io/crates/aingle_zome_types),
//! which contains only the essential types which are used in AIngle DNA
//! code. This crate expands on those types to include all types which AIngle
//! itself depends on.

#![deny(missing_docs)]

pub mod activity;
pub mod app;
pub mod autonomic;
pub mod chain;
pub mod db;
<<<<<<< HEAD
pub mod dgd_op;
=======
pub mod dht_op;
>>>>>>> master
pub mod dna;
pub mod element;
pub mod entry;
pub mod fixt;
pub mod header;
pub mod link;
mod macros;
pub mod metadata;
pub mod prelude;
pub mod signal;
pub mod timestamp;
pub mod validate;

// #[cfg(test)]
pub mod test_utils;

pub use entry::EntryHashed;
pub use timestamp::{Timestamp, TimestampKey};
