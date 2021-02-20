use crate::core::ribosome::CallContext;
use crate::core::ribosome::RibosomeT;
use aingle_types::prelude::*;
use std::sync::Arc;
use aingle_wasmer_host::prelude::WasmError;

#[allow(clippy::extra_unused_lifetimes)]
pub fn get<'a>(
    _ribosome: Arc<impl RibosomeT>,
    call_context: Arc<CallContext>,
    input: GetInput,
) -> Result<Option<Element>, WasmError> {
<<<<<<< HEAD
    let GetInput{ any_dgd_hash, get_options } = input;
=======
    let GetInput{ any_dht_hash, get_options } = input;
>>>>>>> master

    // Get the network from the context
    let network = call_context.host_access.network().clone();

    // timeouts must be handled by the network
    tokio_safe_block_on::tokio_safe_block_forever_on(async move {
        let maybe_element = call_context
            .host_access
            .workspace()
            .write()
            .await
            .cascade(network)
<<<<<<< HEAD
            .dgd_get(any_dgd_hash, get_options)
=======
            .dht_get(any_dht_hash, get_options)
>>>>>>> master
            .await
            .map_err(|cascade_error| WasmError::Host(cascade_error.to_string()))?;

        Ok(maybe_element)
    })
}

// we are relying on the create tests to show the commit/get round trip
// See commit_entry.rs