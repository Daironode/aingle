use aingle_p2p::AIngleP2pCell;
use aingle_p2p::AIngleP2pCellT;
use aingle_state::prelude::*;
use aingle_types::prelude::*;
use aingle_zome_types::TryInto;
use tracing::*;

use crate::core::queue_consumer::WorkComplete;

use super::error::WorkflowResult;

#[cfg(test)]
mod tests;

#[instrument(skip(vault, network))]
/// Send validation receipts to their authors in serial and without waiting for
/// responses.
/// TODO: Currently still waiting for responses because we don't have a network call
/// that doesn't.
pub async fn validation_receipt_workflow(
    vault: EnvWrite,
    network: &mut AIngleP2pCell,
) -> WorkflowResult<WorkComplete> {
    // Get the env and keystore
    let keystore = vault.keystore();
    // Who we are.
    let validator = network.from_agent();

    // Get out all ops that are marked for sending receipt.
    // FIXME: Test this query.
    let receipts = vault
        .async_reader({
            let validator = validator.clone();
            move |txn| {
                let mut stmt = txn.prepare(
                    "
            SELECT Header.author, SgdOp.hash, SgdOp.validation_status,
            SgdOp.when_integrated_ns
            From SgdOp
            JOIN Header ON SgdOp.header_hash = Header.hash
            WHERE
            SgdOp.require_receipt = 1
            AND
            SgdOp.when_integrated_ns IS NOT NULL
            AND
            SgdOp.validation_status IS NOT NULL
            ",
                )?;
                let ops = stmt
                    .query_and_then([], |r| {
                        let author: AgentPubKey = r.get("author")?;
                        let sgd_op_hash = r.get("hash")?;
                        let validation_status = r.get("validation_status")?;
                        let when_integrated = from_blob::<Timestamp>(r.get("when_integrated_ns")?)?;
                        StateQueryResult::Ok((
                            ValidationReceipt {
                                sgd_op_hash,
                                validation_status,
                                validator: validator.clone(),
                                when_integrated,
                            },
                            author,
                        ))
                    })?
                    .collect::<StateQueryResult<Vec<_>>>()?;
                StateQueryResult::Ok(ops)
            }
        })
        .await?;

    // Send the validation receipts
    for (receipt, author) in receipts {
        // Don't send receipt to self.
        if author == validator {
            continue;
        }

        let op_hash = receipt.sgd_op_hash.clone();

        // Sign on the dotted line.
        let receipt = receipt.sign(&keystore).await?;

        // Send it and don't wait for response.
        // TODO: When networking has a send without response we can use that
        // instead of waiting for response.
        if let Err(e) = network
            .send_validation_receipt(author, receipt.try_into()?)
            .await
        {
            // No one home, they will need to publish again.
            info!(failed_send_receipt = ?e);
        }
        // Attempted to send the receipt so we now mark
        // it to not send in the future.
        vault
            .async_commit(|txn| set_require_receipt(txn, op_hash, false))
            .await?;
    }

    Ok(WorkComplete::Complete)
}
