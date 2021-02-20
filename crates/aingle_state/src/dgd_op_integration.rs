//! Various types for the databases involved in the DgdOp integration workflow

use fallible_iterator::FallibleIterator;
use aingle_hash::*;
use aingle_lmdb::buffer::KvBufFresh;
use aingle_lmdb::db::INTEGRATED_DGD_OPS;
use aingle_lmdb::error::DatabaseError;
use aingle_lmdb::error::DatabaseResult;
use aingle_lmdb::prelude::BufferedStore;
use aingle_lmdb::prelude::EnvironmentRead;
use aingle_lmdb::prelude::GetDb;
use aingle_lmdb::prelude::Readable;
use aingle_p2p::dgd_arc::DgdArc;
use aingle_middleware_bytes::prelude::*;
use aingle_types::prelude::*;
use aingle_zome_types::validate::ValidationStatus;

/// Database type for AuthoredDgdOps
/// Buffer for accessing [DgdOp]s that you authored and finding the amount of validation receipts
pub type AuthoredDgdOpsStore = KvBufFresh<AuthoredDgdOpsKey, AuthoredDgdOpsValue>;

/// The key type for the AuthoredDgdOps db: a DgdOpHash
pub type AuthoredDgdOpsKey = DgdOpHash;

/// A type for storing in databases that only need the hashes.
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq)]
pub struct AuthoredDgdOpsValue {
    /// Signatures and hashes of the op
    pub op: DgdOpLight,
    /// Validation receipts received
    pub receipt_count: u32,
    /// Time last published, None if never published
    pub last_publish_time: Option<Timestamp>,
}

impl AuthoredDgdOpsValue {
    /// Create a new value from a DgdOpLight with no receipts and no timestamp
    pub fn from_light(op: DgdOpLight) -> Self {
        Self {
            op,
            receipt_count: 0,
            last_publish_time: None,
        }
    }
}

/// Database type for IntegrationLimbo: the queue of ops ready to be integrated.
pub type IntegrationLimboStore = KvBufFresh<IntegrationLimboKey, IntegrationLimboValue>;

/// Database type for IntegratedDgdOps
/// [DgdOp]s that have already been integrated
pub type IntegratedDgdOpsStore = KvBufFresh<DgdOpHash, IntegratedDgdOpsValue>;

/// Buffer that adds query logic to the IntegratedDgdOpsStore
pub struct IntegratedDgdOpsBuf {
    store: IntegratedDgdOpsStore,
}

impl std::ops::Deref for IntegratedDgdOpsBuf {
    type Target = IntegratedDgdOpsStore;
    fn deref(&self) -> &Self::Target {
        &self.store
    }
}

impl std::ops::DerefMut for IntegratedDgdOpsBuf {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.store
    }
}

impl BufferedStore for IntegratedDgdOpsBuf {
    type Error = DatabaseError;
    fn flush_to_txn_ref(
        &mut self,
        writer: &mut aingle_lmdb::prelude::Writer,
    ) -> Result<(), Self::Error> {
        self.store.flush_to_txn_ref(writer)
    }
}

/// The key type for the IntegrationLimbo db is just a DgdOpHash
pub type IntegrationLimboKey = DgdOpHash;

/// A type for storing in databases that only need the hashes.
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq)]
pub struct IntegratedDgdOpsValue {
    /// The op's validation status
    pub validation_status: ValidationStatus,
    /// Signatures and hashes of the op
    pub op: DgdOpLight,
    /// Time when the op was integrated
    pub when_integrated: Timestamp,
}

/// A type for storing in databases that only need the hashes.
#[derive(Clone, Debug, Serialize, Deserialize, Eq, PartialEq)]
pub struct IntegrationLimboValue {
    /// The op's validation status
    pub validation_status: ValidationStatus,
    /// The op
    pub op: DgdOpLight,
}

impl IntegratedDgdOpsBuf {
    /// Create a new buffer for the IntegratedDgdOpsStore
    pub fn new(env: EnvironmentRead) -> DatabaseResult<Self> {
        let db = env.get_db(&*INTEGRATED_DGD_OPS).unwrap();
        Ok(Self {
            store: IntegratedDgdOpsStore::new(env, db),
        })
    }

    /// simple get by dgd_op_hash
    pub fn get(&'_ self, op_hash: &DgdOpHash) -> DatabaseResult<Option<IntegratedDgdOpsValue>> {
        self.store.get(op_hash)
    }

    /// Get ops that match optional queries:
    /// - from a time (Inclusive)
    /// - to a time (Exclusive)
    /// - match a dgd location
    pub fn query<'r, R: Readable>(
        &'r self,
        r: &'r R,
        from: Option<Timestamp>,
        to: Option<Timestamp>,
        dgd_arc: Option<DgdArc>,
    ) -> DatabaseResult<
        Box<
            dyn FallibleIterator<Item = (DgdOpHash, IntegratedDgdOpsValue), Error = DatabaseError>
                + 'r,
        >,
    > {
        Ok(Box::new(
            self.store
                .iter(r)?
                .map(move |(k, v)| Ok((DgdOpHash::from_raw_39_panicky(k.to_vec()), v)))
                .filter_map(move |(k, v)| match from {
                    Some(time) if v.when_integrated >= time => Ok(Some((k, v))),
                    None => Ok(Some((k, v))),
                    _ => Ok(None),
                })
                .filter_map(move |(k, v)| match to {
                    Some(time) if v.when_integrated < time => Ok(Some((k, v))),
                    None => Ok(Some((k, v))),
                    _ => Ok(None),
                })
                .filter_map(move |(k, v)| match dgd_arc {
                    Some(dgd_arc) if dgd_arc.contains(v.op.dgd_basis().get_loc()) => {
                        Ok(Some((k, v)))
                    }
                    None => Ok(Some((k, v))),
                    _ => Ok(None),
                }),
        ))
    }
}
#[cfg(test)]
mod tests {
    use super::*;
    use ::fixt::prelude::*;
    use chrono::Duration;
    use chrono::Utc;
    use aingle_hash::fixt::AnyDgdHashFixturator;
    use aingle_hash::fixt::DgdOpHashFixturator;
    use aingle_hash::fixt::HeaderHashFixturator;
    use aingle_lmdb::buffer::BufferedStore;
    use aingle_lmdb::env::ReadManager;
    use aingle_lmdb::env::WriteManager;
    use aingle_lmdb::test_utils::test_cell_env;
    use pretty_assertions::assert_eq;

    #[tokio::test(threaded_scheduler)]
    async fn test_query() {
        let test_env = test_cell_env();
        let env = test_env.env();
        let env_ref = env.guard();

        // Create some integration values
        let mut expected = Vec::new();
        let mut basis = AnyDgdHashFixturator::new(Predictable);
        let now = Utc::now();
        let same_basis = basis.next().unwrap();
        let mut times = Vec::new();
        times.push(now - Duration::hours(100));
        times.push(now);
        times.push(now + Duration::hours(100));
        let times_exp = times.clone();
        let values = times
            .into_iter()
            .map(|when_integrated| IntegratedDgdOpsValue {
                validation_status: ValidationStatus::Valid,
                op: DgdOpLight::RegisterAgentActivity(fixt!(HeaderHash), basis.next().unwrap()),
                when_integrated: when_integrated.into(),
            });

        // Put them in the db
        {
            let mut dgd_hash = DgdOpHashFixturator::new(Predictable);
            let mut buf = IntegratedDgdOpsBuf::new(env.clone().into()).unwrap();
            for mut value in values {
                buf.put(dgd_hash.next().unwrap(), value.clone()).unwrap();
                expected.push(value.clone());
                value.op = DgdOpLight::RegisterAgentActivity(fixt!(HeaderHash), same_basis.clone());
                buf.put(dgd_hash.next().unwrap(), value.clone()).unwrap();
                expected.push(value.clone());
            }
            env_ref
                .with_commit(|writer| buf.flush_to_txn(writer))
                .unwrap();
        }

        // Check queries
        {
            let reader = env_ref.reader().unwrap();
            let buf = IntegratedDgdOpsBuf::new(env.clone().into()).unwrap();
            // No filter
            let mut r = buf
                .query(&reader, None, None, None)
                .unwrap()
                .map(|(_, v)| Ok(v))
                .collect::<Vec<_>>()
                .unwrap();
            r.sort_by_key(|v| v.when_integrated.clone());
            assert_eq!(&r[..], &expected[..]);
            // From now
            let mut r = buf
                .query(&reader, Some(times_exp[1].clone().into()), None, None)
                .unwrap()
                .map(|(_, v)| Ok(v))
                .collect::<Vec<_>>()
                .unwrap();
            r.sort_by_key(|v| v.when_integrated.clone());
            assert!(r.contains(&expected[2]));
            assert!(r.contains(&expected[4]));
            assert!(r.contains(&expected[3]));
            assert!(r.contains(&expected[5]));
            assert_eq!(r.len(), 4);
            // From ages ago till 1hr in future
            let ages_ago = times_exp[0] - Duration::weeks(5);
            let future = times_exp[1] + Duration::hours(1);
            let mut r = buf
                .query(&reader, Some(ages_ago.into()), Some(future.into()), None)
                .unwrap()
                .map(|(_, v)| Ok(v))
                .collect::<Vec<_>>()
                .unwrap();
            r.sort_by_key(|v| v.when_integrated.clone());

            assert!(r.contains(&expected[0]));
            assert!(r.contains(&expected[1]));
            assert!(r.contains(&expected[2]));
            assert!(r.contains(&expected[3]));
            assert_eq!(r.len(), 4);
            // Same basis
            let ages_ago = times_exp[0] - Duration::weeks(5);
            let future = times_exp[1] + Duration::hours(1);
            let mut r = buf
                .query(
                    &reader,
                    Some(ages_ago.into()),
                    Some(future.into()),
                    Some(DgdArc::new(same_basis.get_loc(), 1)),
                )
                .unwrap()
                .map(|(_, v)| Ok(v))
                .collect::<Vec<_>>()
                .unwrap();
            r.sort_by_key(|v| v.when_integrated.clone());
            assert!(r.contains(&expected[1]));
            assert!(r.contains(&expected[3]));
            assert_eq!(r.len(), 2);
            // Same basis all
            let mut r = buf
                .query(
                    &reader,
                    None,
                    None,
                    Some(DgdArc::new(same_basis.get_loc(), 1)),
                )
                .unwrap()
                .map(|(_, v)| Ok(v))
                .collect::<Vec<_>>()
                .unwrap();
            r.sort_by_key(|v| v.when_integrated.clone());
            assert!(r.contains(&expected[1]));
            assert!(r.contains(&expected[3]));
            assert!(r.contains(&expected[5]));
            assert_eq!(r.len(), 3);
        }
    }
}
