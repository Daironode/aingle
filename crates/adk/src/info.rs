use crate::prelude::*;

/// Trivial wrapper for `__agent_info` host function.
/// Agent info input struct is `()` so the function call simply looks like this:
///
/// ```ignore
/// let agent_info = agent_info()?;
/// ```
///
/// the [ `AgentInfo` ] is the current agent's original pubkey/address that they joined the network with
/// and their most recent pubkey/address.
pub fn agent_info() -> ExternResult<AgentInfo> {
    ADK.with(|h| h.borrow().agent_info(()))
}

/// @todo Not implemented
pub fn app_info() -> ExternResult<AppInfo> {
    ADK.with(|h| h.borrow().app_info(()))
}

/// @todo Not implemented
pub fn saf_info() -> ExternResult<SafInfo> {
    ADK.with(|h| h.borrow().saf_info(()))
}

/// Get the zome information.
/// There are no inputs to [ `zome_info` ].
///
/// Zome information includes saf name, hash, zome name and properties.
///
/// In general any aingle compatible wasm can be compiled and run in any zome so the zome info
/// needs to be looked up at runtime to e.g. know where to send/receive `call_remote` rpc calls to.
pub fn zome_info() -> ExternResult<ZomeInfo> {
    ADK.with(|h| h.borrow().zome_info(()))
}

/// @todo Not implemented
pub fn call_info() -> ExternResult<CallInfo> {
    ADK.with(|h| h.borrow().call_info(()))
}
