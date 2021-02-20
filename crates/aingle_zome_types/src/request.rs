//! Types for requesting metadata

<<<<<<< HEAD
use aingle_middleware_bytes::prelude::*;
=======
use aingle_serialized_bytes::prelude::*;
>>>>>>> master

#[derive(Debug, Hash, PartialEq, Eq, Clone, Serialize, Deserialize)]
/// Metadata that can be requested on a basis
pub struct MetadataRequest {
    /// Get all the headers on an entry.
    /// Invalid request on a header.
    pub all_valid_headers: bool,
    // TODO: Implement after validation
    /// Placeholder
    pub all_invalid_headers: bool,
    /// Get all the deletes on a header
    pub all_deletes: bool,
    /// Get all the updates on an entry or header
    pub all_updates: bool,
    /// Placeholder
    pub follow_redirects: bool,
    /// Request the status of an entry.
    /// This is faster then getting all the headers
    /// and checking for live headers.
<<<<<<< HEAD
    pub entry_dgd_status: bool,
=======
    pub entry_dht_status: bool,
>>>>>>> master
}

impl Default for MetadataRequest {
    fn default() -> Self {
        Self {
            all_valid_headers: true,
            all_invalid_headers: false,
            all_deletes: true,
            all_updates: true,
            follow_redirects: false,
<<<<<<< HEAD
            entry_dgd_status: false,
=======
            entry_dht_status: false,
>>>>>>> master
        }
    }
}
