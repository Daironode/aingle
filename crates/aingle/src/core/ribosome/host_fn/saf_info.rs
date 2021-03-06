use crate::core::ribosome::CallContext;
use crate::core::ribosome::RibosomeT;
use std::sync::Arc;
use aingle_wasmer_host::prelude::WasmError;
use aingle_zome_types::info::SafInfo;

pub fn saf_info(
    _ribosome: Arc<impl RibosomeT>,
    _call_context: Arc<CallContext>,
    _input: (),
) -> Result<SafInfo, WasmError> {
    unimplemented!()
}

