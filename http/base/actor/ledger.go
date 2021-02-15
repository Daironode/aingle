
package actor

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/ledger"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract/event"
	cstate " github.com/Daironode/aingle/smartcontract/states"
)

const (
	REQ_TIMEOUT    = 5
	ERR_ACTOR_COMM = "[http] Actor comm error: %v"
)

//GetHeaderByHeight from ledger
func GetHeaderByHeight(height uint32) (*types.Header, error) {
	return ledger.DefLedger.GetHeaderByHeight(height)
}

//GetBlockByHeight from ledger
func GetBlockByHeight(height uint32) (*types.Block, error) {
	return ledger.DefLedger.GetBlockByHeight(height)
}

//GetBlockHashFromStore from ledger
func GetBlockHashFromStore(height uint32) common.Uint256 {
	return ledger.DefLedger.GetBlockHash(height)
}

//CurrentBlockHash from ledger
func CurrentBlockHash() common.Uint256 {
	return ledger.DefLedger.GetCurrentBlockHash()
}

//GetBlockFromStore from ledger
func GetBlockFromStore(hash common.Uint256) (*types.Block, error) {
	return ledger.DefLedger.GetBlockByHash(hash)
}

//GetCurrentBlockHeight from ledger
func GetCurrentBlockHeight() uint32 {
	return ledger.DefLedger.GetCurrentBlockHeight()
}

//GetTransaction from ledger
func GetTransaction(hash common.Uint256) (*types.Transaction, error) {
	return ledger.DefLedger.GetTransaction(hash)
}

//GetStorageItem from ledger
func GetStorageItem(address common.Address, key []byte) ([]byte, error) {
	return ledger.DefLedger.GetStorageItem(address, key)
}

//GetContractStateFromStore from ledger
func GetContractStateFromStore(hash common.Address) (*payload.DeployCode, error) {
	hash = updateNativeSCAddr(hash)
	return ledger.DefLedger.GetContractState(hash)
}

//GetTxnWithHeightByTxHash from ledger
func GetTxnWithHeightByTxHash(hash common.Uint256) (uint32, *types.Transaction, error) {
	tx, height, err := ledger.DefLedger.GetTransactionWithHeight(hash)
	return height, tx, err
}

//PreExecuteContract from ledger
func PreExecuteContract(tx *types.Transaction) (*cstate.PreExecResult, error) {
	return ledger.DefLedger.PreExecuteContract(tx)
}

func PreExecuteContractBatch(tx []*types.Transaction, atomic bool) ([]*cstate.PreExecResult, uint32, error) {
	return ledger.DefLedger.PreExecuteContractBatch(tx, atomic)
}

//GetEventNotifyByTxHash from ledger
func GetEventNotifyByTxHash(txHash common.Uint256) (*event.ExecuteNotify, error) {
	return ledger.DefLedger.GetEventNotifyByTx(txHash)
}

//GetEventNotifyByHeight from ledger
func GetEventNotifyByHeight(height uint32) ([]*event.ExecuteNotify, error) {
	return ledger.DefLedger.GetEventNotifyByBlock(height)
}

//GetMerkleProof from ledger
func GetMerkleProof(proofHeight uint32, rootHeight uint32) ([]common.Uint256, error) {
	return ledger.DefLedger.GetMerkleProof(proofHeight, rootHeight)
}

func GetCrossChainMsg(height uint32) (*types.CrossChainMsg, error) {
	return ledger.DefLedger.GetCrossChainMsg(height)
}

func GetCrossStatesProof(height uint32, key []byte) ([]byte, error) {
	return ledger.DefLedger.GetCrossStatesProof(height, key)
}
