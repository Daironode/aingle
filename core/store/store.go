
package store

import (
	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/core/store/overlaydb"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract/event"
	cstates " github.com/Daironode/aingle/smartcontract/states"
)

type ExecuteResult struct {
	WriteSet        *overlaydb.MemDB
	Hash            common.Uint256
	MerkleRoot      common.Uint256
	CrossStates     []common.Uint256
	CrossStatesRoot common.Uint256
	Notify          []*event.ExecuteNotify
}

// LedgerStore provides func with store package.
type LedgerStore interface {
	InitLedgerStoreWithGenesisBlock(genesisblock *types.Block, defaultBookkeeper []keypair.PublicKey) error
	Close() error
	AddHeaders(headers []*types.Header) error
	AddBlock(block *types.Block, ccMsg *types.CrossChainMsg, stateMerkleRoot common.Uint256) error
	ExecuteBlock(b *types.Block) (ExecuteResult, error)                                       // called by consensus
	SubmitBlock(b *types.Block, crossChainMsg *types.CrossChainMsg, exec ExecuteResult) error // called by consensus
	GetStateMerkleRoot(height uint32) (result common.Uint256, err error)
	GetCurrentBlockHash() common.Uint256
	GetCurrentBlockHeight() uint32
	GetCurrentHeaderHeight() uint32
	GetCurrentHeaderHash() common.Uint256
	GetBlockHash(height uint32) common.Uint256
	GetHeaderByHash(blockHash common.Uint256) (*types.Header, error)
	GetRawHeaderByHash(blockHash common.Uint256) (*types.RawHeader, error)
	GetHeaderByHeight(height uint32) (*types.Header, error)
	GetBlockByHash(blockHash common.Uint256) (*types.Block, error)
	GetBlockByHeight(height uint32) (*types.Block, error)
	GetTransaction(txHash common.Uint256) (*types.Transaction, uint32, error)
	IsContainBlock(blockHash common.Uint256) (bool, error)
	IsContainTransaction(txHash common.Uint256) (bool, error)
	GetBlockRootWithNewTxRoots(startHeight uint32, txRoots []common.Uint256) common.Uint256
	GetMerkleProof(m, n uint32) ([]common.Uint256, error)
	GetContractState(contractHash common.Address) (*payload.DeployCode, error)
	GetBookkeeperState() (*states.BookkeeperState, error)
	GetStorageItem(key *states.StorageKey) (*states.StorageItem, error)
	PreExecuteContract(tx *types.Transaction) (*cstates.PreExecResult, error)
	PreExecuteContractBatch(txes []*types.Transaction, atomic bool) ([]*cstates.PreExecResult, uint32, error)
	GetEventNotifyByTx(tx common.Uint256) (*event.ExecuteNotify, error)
	GetEventNotifyByBlock(height uint32) ([]*event.ExecuteNotify, error)

	//cross chain states root
	GetCrossStatesRoot(height uint32) (common.Uint256, error)
	GetCrossChainMsg(height uint32) (*types.CrossChainMsg, error)
	GetCrossStatesProof(height uint32, key []byte) ([]byte, error)
	EnableBlockPrune(numBeforeCurr uint32)
}
