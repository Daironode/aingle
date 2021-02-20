//! The workflow and queue consumer for DgdOp integration

use super::error::WorkflowResult;
use super::integrate_dgd_ops_workflow::integrate_single_data;
use super::produce_dgd_ops_workflow::dgd_op_light::error::DgdOpConvertResult;
use super::sys_validation_workflow::counterfeit_check;
use crate::core::queue_consumer::TriggerSender;
use aingle_hash::AgentPubKey;
use aingle_hash::DgdOpHash;
use aingle_cascade::integrate_single_metadata;
use aingle_lmdb::buffer::BufferedStore;
use aingle_lmdb::buffer::KvBufFresh;
use aingle_lmdb::db::INTEGRATED_DGD_OPS;
use aingle_lmdb::db::INTEGRATION_LIMBO;
use aingle_lmdb::env::EnvironmentWrite;
use aingle_lmdb::error::DatabaseResult;
use aingle_lmdb::prelude::EnvironmentRead;
use aingle_lmdb::prelude::GetDb;
use aingle_lmdb::prelude::IntegratedPrefix;
use aingle_lmdb::prelude::PendingPrefix;
use aingle_lmdb::prelude::Writer;
use aingle_state::prelude::*;
use aingle_types::prelude::*;
use aingle_zome_types::query::HighestObserved;
use tracing::instrument;

#[cfg(test)]
mod test;

#[instrument(skip(state_env, sys_validation_trigger, ops))]
pub async fn incoming_dgd_ops_workflow(
    state_env: &EnvironmentWrite,
    mut sys_validation_trigger: TriggerSender,
    ops: Vec<(aingle_hash::DgdOpHash, aingle_types::dgd_op::DgdOp)>,
    from_agent: Option<AgentPubKey>,
) -> WorkflowResult<()> {
    // set up our workspace
    let mut workspace = IncomingDgdOpsWorkspace::new(state_env.clone().into())?;

    // add incoming ops to the validation limbo
    for (hash, op) in ops {
        if !workspace.op_exists(&hash)? {
            tracing::debug!(?hash, ?op);
            if should_keep(&op).await? {
                workspace.add_to_pending(hash, op, from_agent.clone())?;
            } else {
                tracing::warn!(
                    msg = "Dropping op because it failed counterfeit checks",
                    ?op
                );
            }
        }
    }

    // commit our transaction
    let writer: crate::core::queue_consumer::OneshotWriter = state_env.clone().into();

    writer.with_writer(|writer| Ok(workspace.flush_to_txn(writer)?))?;

    // trigger validation of queued ops
    sys_validation_trigger.trigger();

    Ok(())
}

#[instrument(skip(op))]
/// If this op fails the counterfeit check it should be dropped
async fn should_keep(op: &DgdOp) -> WorkflowResult<bool> {
    let header = op.header();
    let signature = op.signature();
    Ok(counterfeit_check(signature, &header).await?)
}

#[allow(missing_docs)]
pub struct IncomingDgdOpsWorkspace {
    pub integration_limbo: IntegrationLimboStore,
    pub integrated_dgd_ops: IntegratedDgdOpsStore,
    pub validation_limbo: ValidationLimboStore,
    pub element_pending: ElementBuf<PendingPrefix>,
    pub meta_pending: MetadataBuf<PendingPrefix>,
    pub meta_integrated: MetadataBuf<IntegratedPrefix>,
}

impl Workspace for IncomingDgdOpsWorkspace {
    fn flush_to_txn_ref(&mut self, writer: &mut Writer) -> WorkspaceResult<()> {
        self.validation_limbo.0.flush_to_txn_ref(writer)?;
        self.element_pending.flush_to_txn_ref(writer)?;
        self.meta_pending.flush_to_txn_ref(writer)?;
        self.meta_integrated.flush_to_txn_ref(writer)?;
        Ok(())
    }
}

impl IncomingDgdOpsWorkspace {
    pub fn new(env: EnvironmentRead) -> WorkspaceResult<Self> {
        let db = env.get_db(&*INTEGRATED_DGD_OPS)?;
        let integrated_dgd_ops = KvBufFresh::new(env.clone(), db);

        let db = env.get_db(&*INTEGRATION_LIMBO)?;
        let integration_limbo = KvBufFresh::new(env.clone(), db);

        let validation_limbo = ValidationLimboStore::new(env.clone())?;

        let element_pending = ElementBuf::pending(env.clone())?;
        let meta_pending = MetadataBuf::pending(env.clone())?;

        let meta_integrated = MetadataBuf::vault(env)?;

        Ok(Self {
            integration_limbo,
            integrated_dgd_ops,
            validation_limbo,
            element_pending,
            meta_pending,
            meta_integrated,
        })
    }

    fn add_to_pending(
        &mut self,
        hash: DgdOpHash,
        op: DgdOp,
        from_agent: Option<AgentPubKey>,
    ) -> DgdOpConvertResult<()> {
        let basis = op.dgd_basis();
        let op_light = op.to_light();
        tracing::debug!(?op_light);

        // register the highest observed header in an agents chain
        if let DgdOp::RegisterAgentActivity(_, header) = &op {
            self.meta_integrated.register_activity_observed(
                header.author(),
                HighestObserved {
                    header_seq: header.header_seq(),
                    hash: vec![op_light.header_hash().clone()],
                },
            )?;
        }

        integrate_single_data(op, &mut self.element_pending)?;
        integrate_single_metadata(
            op_light.clone(),
            &self.element_pending,
            &mut self.meta_pending,
        )?;
        let vlv = ValidationLimboValue {
            status: ValidationLimboStatus::Pending,
            op: op_light,
            basis,
            time_added: timestamp::now(),
            last_try: None,
            num_tries: 0,
            from_agent,
        };
        self.validation_limbo.put(hash, vlv)?;
        Ok(())
    }

    pub fn op_exists(&self, hash: &DgdOpHash) -> DatabaseResult<bool> {
        Ok(self.integrated_dgd_ops.contains(&hash)?
            || self.integration_limbo.contains(&hash)?
            || self.validation_limbo.contains(&hash)?)
    }
}
