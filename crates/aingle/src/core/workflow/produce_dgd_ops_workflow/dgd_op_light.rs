use error::DgdOpConvertError;
use error::DgdOpConvertResult;
use aingle_hash::EntryHash;
use aingle_hash::HeaderHash;
use aingle_state::element_buf::ElementBuf;
use aingle_types::prelude::*;

pub mod error;

use aingle_lmdb::prelude::PrefixType;
use tracing::*;

#[cfg(test)]
mod tests;

/// Convert a DgdOpLight into a DgdOp (render all the hashes to values)
/// This only checks the ElementVault so can only be used with ops that you are
/// an authority or author of.
pub fn light_to_op<P: PrefixType>(
    op: DgdOpLight,
    cas: &ElementBuf<P>,
) -> DgdOpConvertResult<DgdOp> {
    let op_name = format!("{:?}", op);
    match op {
        DgdOpLight::StoreElement(h, _, _) => {
            let (header, entry) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner();
            // TODO: Could use this signature? Is it the same?
            // Should we not be storing the signature in the DgdOpLight?
            let (header, sig) = header.into_header_and_signature();
            let entry = entry.into_option().map(Box::new);
            Ok(DgdOp::StoreElement(sig, header.into_content(), entry))
        }
        DgdOpLight::StoreEntry(h, _, _) => {
            let (header, entry) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner();
            let (header, sig) = header.into_header_and_signature();
            let header = match header.into_content() {
                Header::Create(c) => NewEntryHeader::Create(c),
                Header::Update(c) => NewEntryHeader::Update(c),
                _ => return Err(DgdOpConvertError::HeaderEntryMismatch),
            };

            let entry = match header.visibility() {
                // Entry must be here because it's a StoreEntry
                EntryVisibility::Public => entry
                    .into_option()
                    .ok_or_else(|| DgdOpConvertError::MissingData(header.entry().clone().into()))?,
                // If the entry is not here and you were meant to have access
                // it's because you were using a database without access to private entries
                // If not then you should handle this error
                EntryVisibility::Private => entry
                    .into_option()
                    .ok_or(DgdOpConvertError::StoreEntryOnPrivate)?,
            };
            Ok(DgdOp::StoreEntry(sig, header, Box::new(entry)))
        }
        DgdOpLight::RegisterAgentActivity(h, _) => {
            let (header, sig) = cas
                .get_header(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_header_and_signature();
            Ok(DgdOp::RegisterAgentActivity(sig, header.into_content()))
        }
        DgdOpLight::RegisterUpdatedContent(h, _, _) => {
            let (header, entry) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner();
            let (header, sig) = header.into_header_and_signature();
            let header = match header.into_content() {
                Header::Update(u) => u,
                h => {
                    return Err(DgdOpConvertError::HeaderMismatch(
                        format!("{:?}", h),
                        op_name,
                    ));
                }
            };
            let entry = entry.into_option().map(Box::new);
            Ok(DgdOp::RegisterUpdatedContent(sig, header, entry))
        }
        DgdOpLight::RegisterUpdatedElement(h, _, _) => {
            let (header, entry) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner();
            let (header, sig) = header.into_header_and_signature();
            let header = match header.into_content() {
                Header::Update(u) => u,
                h => {
                    return Err(DgdOpConvertError::HeaderMismatch(
                        format!("{:?}", h),
                        op_name,
                    ));
                }
            };
            let entry = entry.into_option().map(Box::new);
            Ok(DgdOp::RegisterUpdatedElement(sig, header, entry))
        }
        DgdOpLight::RegisterDeletedBy(header_hash, _) => {
            let (header, sig) = get_element_delete(header_hash, op_name, &cas)?;
            Ok(DgdOp::RegisterDeletedBy(sig, header))
        }
        DgdOpLight::RegisterDeletedEntryHeader(header_hash, _) => {
            let (header, sig) = get_element_delete(header_hash, op_name, &cas)?;
            Ok(DgdOp::RegisterDeletedEntryHeader(sig, header))
        }
        DgdOpLight::RegisterAddLink(h, _) => {
            let (header, sig) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner()
                .0
                .into_header_and_signature();
            let header = match header.into_content() {
                Header::CreateLink(u) => u,
                h => {
                    return Err(DgdOpConvertError::HeaderMismatch(
                        format!("{:?}", h),
                        op_name,
                    ));
                }
            };
            Ok(DgdOp::RegisterAddLink(sig, header))
        }
        DgdOpLight::RegisterRemoveLink(h, _) => {
            let (header, sig) = cas
                .get_element(&h)?
                .ok_or_else(|| DgdOpConvertError::MissingData(h.into()))?
                .into_inner()
                .0
                .into_header_and_signature();
            let header = match header.into_content() {
                Header::DeleteLink(u) => u,
                h => {
                    return Err(DgdOpConvertError::HeaderMismatch(
                        format!("{:?}", h),
                        op_name,
                    ));
                }
            };
            Ok(DgdOp::RegisterRemoveLink(sig, header))
        }
    }
}

fn get_element_delete<P: PrefixType>(
    header_hash: HeaderHash,
    op_name: String,
    cas: &ElementBuf<P>,
) -> DgdOpConvertResult<(header::Delete, Signature)> {
    let (header, sig) = cas
        .get_header(&header_hash)?
        .ok_or_else(|| DgdOpConvertError::MissingData(header_hash.into()))?
        .into_header_and_signature();
    match header.into_content() {
        Header::Delete(u) => Ok((u, sig)),
        h => Err(DgdOpConvertError::HeaderMismatch(
            format!("{:?}", h),
            op_name,
        )),
    }
}

#[instrument(skip(cas))]
async fn get_entry_hash_for_header(
    header_hash: &HeaderHash,
    cas: &ElementBuf,
) -> DgdOpConvertResult<EntryHash> {
    debug!(%header_hash);
    let entry = cas
        .get_header(header_hash)?
        .and_then(|e| e.header().entry_data().map(|(hash, _)| hash.clone()));
    entry.ok_or_else(|| DgdOpConvertError::MissingEntryDataForHeader(header_hash.clone()))
}
