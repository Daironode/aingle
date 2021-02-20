//! Data structures representing the operations that can be performed within a AIngle DGD.
//!
//! See the [item-level documentation for `DgdOp`][DgdOp] for more details.
//!
//! [DgdOp]: enum.DgdOp.html

use crate::element::ElementGroup;
use crate::header::NewEntryHeader;
use crate::prelude::*;
use error::DgdOpError;
use error::DgdOpResult;
use aingle_hash::hash_type;
use aingle_hash::HashableContentBytes;
use aingle_zome_types::prelude::*;
use serde::Deserialize;
use serde::Serialize;

#[allow(missing_docs)]
pub mod error;

/// A unit of DGD gossip. Used to notify an authority of new (meta)data to hold
/// as well as changes to the status of already held data.
#[derive(
    Clone, Debug, Serialize, Deserialize, SerializedBytes, Eq, PartialEq, derive_more::Display,
)]
pub enum DgdOp {
    #[display(fmt = "StoreElement")]
    /// Used to notify the authority for a header that it has been created.
    ///
    /// Conceptually, authorities receiving this `DgdOp` do three things:
    ///
    /// - Ensure that the element passes validation.
    /// - Store the header into their DGD shard.
    /// - Store the entry into their CAS.
    ///   - Note: they do not become responsible for keeping the set of
    ///     references from that entry up-to-date.
    StoreElement(Signature, Header, Option<Box<Entry>>),

    #[display(fmt = "StoreEntry")]
    /// Used to notify the authority for an entry that it has been created
    /// anew. (The same entry can be created more than once.)
    ///
    /// Conceptually, authorities receiving this `DgdOp` do four things:
    ///
    /// - Ensure that the element passes validation.
    /// - Store the entry into their DGD shard.
    /// - Store the header into their CAS.
    ///   - Note: they do not become responsible for keeping the set of
    ///     references from that header up-to-date.
    /// - Add a "created-by" reference from the entry to the hash of the header.
    ///
    /// TODO: document how those "created-by" references are stored in
    /// reality.
    StoreEntry(Signature, NewEntryHeader, Box<Entry>),

    #[display(fmt = "RegisterAgentActivity")]
    /// Used to notify the authority for an agent's public key that that agent
    /// has committed a new header.
    ///
    /// Conceptually, authorities receiving this `DgdOp` do three things:
    ///
    /// - Ensure that *the header alone* passes surface-level validation.
    /// - Store the header into their DGD shard.
    ///   - FIXME: @artbrock, do they?
    /// - Add an "agent-activity" reference from the public key to the hash
    ///   of the header.
    ///
    /// TODO: document how those "agent-activity" references are stored in
    /// reality.
    RegisterAgentActivity(Signature, Header),

    #[display(fmt = "RegisterUpdatedContent")]
    /// Op for updating an entry.
    /// This is sent to the entry authority.
    // TODO: This entry is here for validation by the entry update header holder
    // link's don't do this. The entry is validated by store entry. Maybe we either
    // need to remove the Entry here or add it to link.
    RegisterUpdatedContent(Signature, header::Update, Option<Box<Entry>>),

    #[display(fmt = "RegisterUpdatedElement")]
    /// Op for updating an element.
    /// This is sent to the element authority.
    RegisterUpdatedElement(Signature, header::Update, Option<Box<Entry>>),

    #[display(fmt = "RegisterDeletedBy")]
    /// Op for registering a Header deletion with the Header authority
    RegisterDeletedBy(Signature, header::Delete),

    #[display(fmt = "RegisterDeletedEntryHeader")]
    /// Op for registering a Header deletion with the Entry authority, so that
    /// the Entry can be marked Dead if all of its Headers have been deleted
    RegisterDeletedEntryHeader(Signature, header::Delete),

    #[display(fmt = "RegisterAddLink")]
    /// Op for adding a link
    RegisterAddLink(Signature, header::CreateLink),

    #[display(fmt = "RegisterRemoveLink")]
    /// Op for removing a link
    RegisterRemoveLink(Signature, header::DeleteLink),
}

/// Show that this type is used as the basis
type DgdBasis = AnyDgdHash;

/// A type for storing in databases that don't need the actual
/// data. Everything is a hash of the type except the signatures.
#[allow(missing_docs)]
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq, derive_more::Display)]
pub enum DgdOpLight {
    #[display(fmt = "StoreElement")]
    StoreElement(HeaderHash, Option<EntryHash>, DgdBasis),
    #[display(fmt = "StoreEntry")]
    StoreEntry(HeaderHash, EntryHash, DgdBasis),
    #[display(fmt = "RegisterAgentActivity")]
    RegisterAgentActivity(HeaderHash, DgdBasis),
    #[display(fmt = "RegisterUpdatedContent")]
    RegisterUpdatedContent(HeaderHash, EntryHash, DgdBasis),
    #[display(fmt = "RegisterUpdatedElement")]
    RegisterUpdatedElement(HeaderHash, EntryHash, DgdBasis),
    #[display(fmt = "RegisterDeletedBy")]
    RegisterDeletedBy(HeaderHash, DgdBasis),
    #[display(fmt = "RegisterDeletedEntryHeader")]
    RegisterDeletedEntryHeader(HeaderHash, DgdBasis),
    #[display(fmt = "RegisterAddLink")]
    RegisterAddLink(HeaderHash, DgdBasis),
    #[display(fmt = "RegisterRemoveLink")]
    RegisterRemoveLink(HeaderHash, DgdBasis),
}

impl DgdOp {
    fn as_unique_form(&self) -> UniqueForm<'_> {
        match self {
            Self::StoreElement(_, header, _) => UniqueForm::StoreElement(header),
            Self::StoreEntry(_, header, _) => UniqueForm::StoreEntry(header),
            Self::RegisterAgentActivity(_, header) => UniqueForm::RegisterAgentActivity(header),
            Self::RegisterUpdatedContent(_, header, _) => {
                UniqueForm::RegisterUpdatedContent(header)
            }
            Self::RegisterUpdatedElement(_, header, _) => {
                UniqueForm::RegisterUpdatedElement(header)
            }
            Self::RegisterDeletedBy(_, header) => UniqueForm::RegisterDeletedBy(header),
            Self::RegisterDeletedEntryHeader(_, header) => {
                UniqueForm::RegisterDeletedEntryHeader(header)
            }
            Self::RegisterAddLink(_, header) => UniqueForm::RegisterAddLink(header),
            Self::RegisterRemoveLink(_, header) => UniqueForm::RegisterRemoveLink(header),
        }
    }

    /// Returns the basis hash which determines which agents will receive this DgdOp
    pub fn dgd_basis(&self) -> AnyDgdHash {
        self.as_unique_form().basis()
    }

    /// Convert a [DgdOp] to a [DgdOpLight] and basis
    pub fn to_light(
        // Hoping one day we can work out how to go from `&Create`
        // to `&Header::Create(Create)` so punting on a reference
        &self,
    ) -> DgdOpLight {
        let basis = self.dgd_basis();
        match self {
            DgdOp::StoreElement(_, h, _) => {
                let e = h.entry_data().map(|(e, _)| e.clone());
                let h = HeaderHash::with_data_sync(h);
                DgdOpLight::StoreElement(h, e, basis)
            }
            DgdOp::StoreEntry(_, h, _) => {
                let e = h.entry().clone();
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::StoreEntry(h, e, basis)
            }
            DgdOp::RegisterAgentActivity(_, h) => {
                let h = HeaderHash::with_data_sync(h);
                DgdOpLight::RegisterAgentActivity(h, basis)
            }
            DgdOp::RegisterUpdatedContent(_, h, _) => {
                let e = h.entry_hash.clone();
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterUpdatedContent(h, e, basis)
            }
            DgdOp::RegisterUpdatedElement(_, h, _) => {
                let e = h.entry_hash.clone();
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterUpdatedElement(h, e, basis)
            }
            DgdOp::RegisterDeletedBy(_, h) => {
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterDeletedBy(h, basis)
            }
            DgdOp::RegisterDeletedEntryHeader(_, h) => {
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterDeletedEntryHeader(h, basis)
            }
            DgdOp::RegisterAddLink(_, h) => {
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterAddLink(h, basis)
            }
            DgdOp::RegisterRemoveLink(_, h) => {
                let h = HeaderHash::with_data_sync(&Header::from(h.clone()));
                DgdOpLight::RegisterRemoveLink(h, basis)
            }
        }
    }

    /// Get the signature for this op
    pub fn signature(&self) -> &Signature {
        match self {
            DgdOp::StoreElement(s, _, _)
            | DgdOp::StoreEntry(s, _, _)
            | DgdOp::RegisterAgentActivity(s, _)
            | DgdOp::RegisterUpdatedContent(s, _, _)
            | DgdOp::RegisterUpdatedElement(s, _, _)
            | DgdOp::RegisterDeletedBy(s, _)
            | DgdOp::RegisterDeletedEntryHeader(s, _)
            | DgdOp::RegisterAddLink(s, _)
            | DgdOp::RegisterRemoveLink(s, _) => s,
        }
    }

    /// Extract inner Signature, Header and Option<Entry> from an op
    pub fn into_inner(self) -> (Signature, Header, Option<Entry>) {
        match self {
            DgdOp::StoreElement(s, h, e) => (s, h, e.map(|e| *e)),
            DgdOp::StoreEntry(s, h, e) => (s, h.into(), Some(*e)),
            DgdOp::RegisterAgentActivity(s, h) => (s, h, None),
            DgdOp::RegisterUpdatedContent(s, h, e) => (s, h.into(), e.map(|e| *e)),
            DgdOp::RegisterUpdatedElement(s, h, e) => (s, h.into(), e.map(|e| *e)),
            DgdOp::RegisterDeletedBy(s, h) => (s, h.into(), None),
            DgdOp::RegisterDeletedEntryHeader(s, h) => (s, h.into(), None),
            DgdOp::RegisterAddLink(s, h) => (s, h.into(), None),
            DgdOp::RegisterRemoveLink(s, h) => (s, h.into(), None),
        }
    }

    /// Get the header from this op
    /// This requires cloning and converting the header
    /// as some ops don't hold the Header type
    pub fn header(&self) -> Header {
        match self {
            DgdOp::StoreElement(_, h, _) => h.clone(),
            DgdOp::StoreEntry(_, h, _) => h.clone().into(),
            DgdOp::RegisterAgentActivity(_, h) => h.clone(),
            DgdOp::RegisterUpdatedContent(_, h, _) => h.clone().into(),
            DgdOp::RegisterUpdatedElement(_, h, _) => h.clone().into(),
            DgdOp::RegisterDeletedBy(_, h) => h.clone().into(),
            DgdOp::RegisterDeletedEntryHeader(_, h) => h.clone().into(),
            DgdOp::RegisterAddLink(_, h) => h.clone().into(),
            DgdOp::RegisterRemoveLink(_, h) => h.clone().into(),
        }
    }
}

impl DgdOpLight {
    /// Get the dgd basis for where to send this op
    pub fn dgd_basis(&self) -> &AnyDgdHash {
        match self {
            DgdOpLight::StoreElement(_, _, b)
            | DgdOpLight::StoreEntry(_, _, b)
            | DgdOpLight::RegisterAgentActivity(_, b)
            | DgdOpLight::RegisterUpdatedContent(_, _, b)
            | DgdOpLight::RegisterUpdatedElement(_, _, b)
            | DgdOpLight::RegisterDeletedBy(_, b)
            | DgdOpLight::RegisterDeletedEntryHeader(_, b)
            | DgdOpLight::RegisterAddLink(_, b)
            | DgdOpLight::RegisterRemoveLink(_, b) => b,
        }
    }
    /// Get the header hash from this op
    pub fn header_hash(&self) -> &HeaderHash {
        match self {
            DgdOpLight::StoreElement(h, _, _)
            | DgdOpLight::StoreEntry(h, _, _)
            | DgdOpLight::RegisterAgentActivity(h, _)
            | DgdOpLight::RegisterUpdatedContent(h, _, _)
            | DgdOpLight::RegisterUpdatedElement(h, _, _)
            | DgdOpLight::RegisterDeletedBy(h, _)
            | DgdOpLight::RegisterDeletedEntryHeader(h, _)
            | DgdOpLight::RegisterAddLink(h, _)
            | DgdOpLight::RegisterRemoveLink(h, _) => h,
        }
    }
}

// FIXME: need to use this in HashableContent
#[allow(missing_docs)]
#[derive(Serialize, Debug)]
pub enum UniqueForm<'a> {
    // As an optimization, we don't include signatures. They would be redundant
    // with headers and therefore would waste hash/comparison time to include.
    StoreElement(&'a Header),
    StoreEntry(&'a NewEntryHeader),
    RegisterAgentActivity(&'a Header),
    RegisterUpdatedContent(&'a header::Update),
    RegisterUpdatedElement(&'a header::Update),
    RegisterDeletedBy(&'a header::Delete),
    RegisterDeletedEntryHeader(&'a header::Delete),
    RegisterAddLink(&'a header::CreateLink),
    RegisterRemoveLink(&'a header::DeleteLink),
}

impl<'a> UniqueForm<'a> {
    fn basis(&'a self) -> AnyDgdHash {
        match self {
            UniqueForm::StoreElement(header) => HeaderHash::with_data_sync(*header).into(),
            UniqueForm::StoreEntry(header) => header.entry().clone().into(),
            UniqueForm::RegisterAgentActivity(header) => header.author().clone().into(),
            UniqueForm::RegisterUpdatedContent(header) => {
                header.original_entry_address.clone().into()
            }
            UniqueForm::RegisterUpdatedElement(header) => {
                header.original_header_address.clone().into()
            }
            UniqueForm::RegisterDeletedBy(header) => header.deletes_address.clone().into(),
            UniqueForm::RegisterDeletedEntryHeader(header) => {
                header.deletes_entry_address.clone().into()
            }
            UniqueForm::RegisterAddLink(header) => header.base_address.clone().into(),
            UniqueForm::RegisterRemoveLink(header) => header.base_address.clone().into(),
        }
    }
}

/// Produce all DgdOps for a Element
pub fn produce_ops_from_element(element: &Element) -> DgdOpResult<Vec<DgdOp>> {
    let op_lights = produce_op_lights_from_elements(vec![element])?;
    let (shh, maybe_entry) = element.clone().into_inner();
    let (header, signature): (Header, Signature) = shh.into_inner().0.into();

    let mut ops = Vec::with_capacity(op_lights.len());

    for op_light in op_lights {
        let signature = signature.clone();
        let header = header.clone();
        let op = match op_light {
            DgdOpLight::StoreElement(_, _, _) => {
                let maybe_entry_box = maybe_entry.clone().into_option().map(Box::new);
                DgdOp::StoreElement(signature, header, maybe_entry_box)
            }
            DgdOpLight::StoreEntry(_, _, _) => {
                let new_entry_header = header.clone().try_into()?;
                let box_entry = match maybe_entry.clone().into_option() {
                    Some(entry) => Box::new(entry),
                    None => {
                        // Entry is private so continue
                        continue;
                    }
                };
                DgdOp::StoreEntry(signature, new_entry_header, box_entry)
            }
            DgdOpLight::RegisterAgentActivity(_, _) => {
                DgdOp::RegisterAgentActivity(signature, header)
            }
            DgdOpLight::RegisterUpdatedContent(_, _, _) => {
                let entry_update = header.try_into()?;
                let maybe_entry_box = maybe_entry.clone().into_option().map(Box::new);
                DgdOp::RegisterUpdatedContent(signature, entry_update, maybe_entry_box)
            }
            DgdOpLight::RegisterUpdatedElement(_, _, _) => {
                let entry_update = header.try_into()?;
                let maybe_entry_box = maybe_entry.clone().into_option().map(Box::new);
                DgdOp::RegisterUpdatedElement(signature, entry_update, maybe_entry_box)
            }
            DgdOpLight::RegisterDeletedEntryHeader(_, _) => {
                let element_delete = header.try_into()?;
                DgdOp::RegisterDeletedEntryHeader(signature, element_delete)
            }
            DgdOpLight::RegisterDeletedBy(_, _) => {
                let element_delete = header.try_into()?;
                DgdOp::RegisterDeletedBy(signature, element_delete)
            }
            DgdOpLight::RegisterAddLink(_, _) => {
                let link_add = header.try_into()?;
                DgdOp::RegisterAddLink(signature, link_add)
            }
            DgdOpLight::RegisterRemoveLink(_, _) => {
                let link_remove = header.try_into()?;
                DgdOp::RegisterRemoveLink(signature, link_remove)
            }
        };
        ops.push(op);
    }
    Ok(ops)
}

/// Produce all the op lights for tese elements
pub fn produce_op_lights_from_elements(headers: Vec<&Element>) -> DgdOpResult<Vec<DgdOpLight>> {
    let length = headers.len();
    let headers_and_hashes = headers.into_iter().map(|e| {
        (
            e.header_address(),
            e.header(),
            e.header().entry_data().map(|(h, _)| h.clone()),
        )
    });
    produce_op_lights_from_iter(headers_and_hashes, length)
}

/// Produce all the op lights from this element group
/// with a shared entry
pub fn produce_op_lights_from_element_group(
    elements: &ElementGroup<'_>,
) -> DgdOpResult<Vec<DgdOpLight>> {
    let len = elements.len();
    let headers_and_hashes = elements.headers_and_hashes();
    let maybe_entry_hash = Some(elements.entry_hash());
    produce_op_lights_from_parts(headers_and_hashes, maybe_entry_hash, len)
}

/// Data minimal clone (no cloning entries) cheap &Element to DgdOpLight conversion
fn produce_op_lights_from_parts<'a>(
    headers_and_hashes: impl Iterator<Item = (&'a HeaderHash, &'a Header)>,
    maybe_entry_hash: Option<&EntryHash>,
    length: usize,
) -> DgdOpResult<Vec<DgdOpLight>> {
    let iter = headers_and_hashes.map(|(head, hash)| (head, hash, maybe_entry_hash.cloned()));
    produce_op_lights_from_iter(iter, length)
}
fn produce_op_lights_from_iter<'a>(
    iter: impl Iterator<Item = (&'a HeaderHash, &'a Header, Option<EntryHash>)>,
    length: usize,
) -> DgdOpResult<Vec<DgdOpLight>> {
    // Each header will have at least 2 ops
    let mut ops = Vec::with_capacity(length * 2);

    for (header_hash, header, maybe_entry_hash) in iter {
        let header_hash = header_hash.clone();

        let store_element_basis = UniqueForm::StoreElement(header).basis();
        let register_activity_basis = UniqueForm::RegisterAgentActivity(header).basis();

        ops.push(DgdOpLight::StoreElement(
            header_hash.clone(),
            maybe_entry_hash.clone(),
            store_element_basis,
        ));
        ops.push(DgdOpLight::RegisterAgentActivity(
            header_hash.clone(),
            register_activity_basis,
        ));

        match header {
            Header::Dna(_)
            | Header::OpenChain(_)
            | Header::CloseChain(_)
            | Header::AgentValidationPkg(_)
            | Header::InitZomesComplete(_) => {}
            Header::CreateLink(link_add) => ops.push(DgdOpLight::RegisterAddLink(
                header_hash,
                UniqueForm::RegisterAddLink(link_add).basis(),
            )),
            Header::DeleteLink(link_remove) => ops.push(DgdOpLight::RegisterRemoveLink(
                header_hash,
                UniqueForm::RegisterRemoveLink(link_remove).basis(),
            )),
            Header::Create(entry_create) => ops.push(DgdOpLight::StoreEntry(
                header_hash,
                maybe_entry_hash.ok_or_else(|| DgdOpError::HeaderWithoutEntry(header.clone()))?,
                UniqueForm::StoreEntry(&NewEntryHeader::Create(entry_create.clone())).basis(),
            )),
            Header::Update(entry_update) => {
                let entry_hash = maybe_entry_hash
                    .ok_or_else(|| DgdOpError::HeaderWithoutEntry(header.clone()))?;
                ops.push(DgdOpLight::StoreEntry(
                    header_hash.clone(),
                    entry_hash.clone(),
                    UniqueForm::StoreEntry(&NewEntryHeader::Update(entry_update.clone())).basis(),
                ));
                ops.push(DgdOpLight::RegisterUpdatedContent(
                    header_hash.clone(),
                    entry_hash.clone(),
                    UniqueForm::RegisterUpdatedContent(entry_update).basis(),
                ));
                ops.push(DgdOpLight::RegisterUpdatedElement(
                    header_hash,
                    entry_hash,
                    UniqueForm::RegisterUpdatedElement(entry_update).basis(),
                ));
            }
            Header::Delete(entry_delete) => {
                // TODO: VALIDATION: This only works if entry_delete.remove_address is either Create
                // or Update
                ops.push(DgdOpLight::RegisterDeletedBy(
                    header_hash.clone(),
                    UniqueForm::RegisterDeletedBy(entry_delete).basis(),
                ));
                ops.push(DgdOpLight::RegisterDeletedEntryHeader(
                    header_hash,
                    UniqueForm::RegisterDeletedEntryHeader(entry_delete).basis(),
                ));
            }
        }
    }
    Ok(ops)
}

// This has to be done manually because the macro
// implements both directions and that isn't possible with references
// TODO: Maybe add a one-way version to aingle_middleware_bytes?
impl<'a> TryFrom<&UniqueForm<'a>> for SerializedBytes {
    type Error = SerializedBytesError;
    fn try_from(u: &UniqueForm<'a>) -> Result<Self, Self::Error> {
        match aingle_middleware_bytes::encode(u) {
            Ok(v) => Ok(SerializedBytes::from(
                aingle_middleware_bytes::UnsafeBytes::from(v),
            )),
            Err(e) => Err(SerializedBytesError::Serialize(e.to_string())),
        }
    }
}

/// A DgdOp paired with its DgdOpHash
pub type DgdOpHashed = AIngleHashed<DgdOp>;

impl HashableContent for DgdOp {
    type HashType = hash_type::DgdOp;

    fn hash_type(&self) -> Self::HashType {
        hash_type::DgdOp
    }

    fn hashable_content(&self) -> HashableContentBytes {
        HashableContentBytes::Content(
            (&self.as_unique_form())
                .try_into()
                .expect("Could not serialize HashableContent"),
        )
    }
}

impl HashableContent for UniqueForm<'_> {
    type HashType = hash_type::DgdOp;

    fn hash_type(&self) -> Self::HashType {
        hash_type::DgdOp
    }

    fn hashable_content(&self) -> HashableContentBytes {
        HashableContentBytes::Content(
            self.try_into()
                .expect("Could not serialize HashableContent"),
        )
    }
}
