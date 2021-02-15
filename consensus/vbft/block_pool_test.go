
package vbft

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	vconfig " github.com/Daironode/aingle/consensus/vbft/config"
	" github.com/Daironode/aingle/core/ledger"
	" github.com/Daironode/aingle/core/signature"
	" github.com/Daironode/aingle/core/types"
)

func buildTestBlockPool(t *testing.T) (*BlockPool, error) {
	store := newTestChainStore(t)
	return newBlockPool(nil, 64, store)
}

func buildTestBlock(t *testing.T, lastBlock *types.Block, lgr *ledger.Ledger) (*Block, error) {
	timestamp := uint32(time.Now().Unix())
	if timestamp <= lastBlock.Header.Timestamp {
		timestamp = lastBlock.Header.Timestamp + 1
	}
	txs := []*types.Transaction{}
	txHash := []common.Uint256{}
	txRoot := common.ComputeMerkleRoot(txHash)
	blockRoot := lgr.GetBlockRootWithNewTxRoots(lastBlock.Header.Height, []common.Uint256{lastBlock.Header.TransactionsRoot, txRoot})

	consensusPayload, err := json.Marshal(&vconfig.VbftBlockInfo{})
	if err != nil {
		t.Fatalf("failed to build consensus payload: %s", err)
	}

	blkHeader := &types.Header{
		PrevBlockHash:    lastBlock.Header.Hash(),
		TransactionsRoot: txRoot,
		BlockRoot:        blockRoot,
		Timestamp:        timestamp,
		Height:           lastBlock.Header.Height + 1,
		ConsensusData:    common.GetNonce(),
		ConsensusPayload: consensusPayload,
		Bookkeepers:      make([]keypair.PublicKey, 0),
		SigData:          make([][]byte, 0),
	}

	// add sigs
	hash := blkHeader.Hash()
	for i := 0; i < 5; i++ {
		acc := testBookkeeperAccounts[i]
		sig, err := signature.Sign(acc, hash[:])
		if err != nil {
			t.Fatalf("bookkeeper %d sign block: %s", i, err)
		}
		blkHeader.Bookkeepers = append(blkHeader.Bookkeepers, acc.PublicKey)
		blkHeader.SigData = append(blkHeader.SigData, sig)
	}

	blk := &types.Block{
		Header:       blkHeader,
		Transactions: txs,
	}
	block := &Block{
		Block: blk,
	}
	return block, nil
}
func TestAddBlock(t *testing.T) {
	blockpool, err := buildTestBlockPool(t)
	if err != nil {
		t.Errorf("buildTestBlockPool err:%s", err)
	}
	defer cleanTestChainStore()

	lastBlock, _ := blockpool.getSealedBlock(0)
	if lastBlock == nil {
		t.Errorf("getblock err")
	}
	t.Logf("block height:%d", blockpool.chainStore.GetChainedBlockNum())
	blk, err := buildTestBlock(t, lastBlock.Block, blockpool.chainStore.db)
	if err != nil {
		t.Errorf("buildTestBlock err:%s", err)
	}
	err = blockpool.chainStore.AddBlock(blk)
	if err != nil {
		t.Errorf("AddBlock err:%s", err)
	}
	merkleRoot, err := blockpool.getExecMerkleRoot(blockpool.chainStore.GetChainedBlockNum())
	if err != nil {
		t.Errorf("getExecMerkleRoot err:%s", err)
	}
	t.Logf("block height:%d,merkleRoot:%s", blockpool.chainStore.GetChainedBlockNum(), merkleRoot.ToHexString())
	err = blockpool.submitBlock(blockpool.chainStore.GetChainedBlockNum())
	if err != nil {
		t.Errorf("submitBlock err:%s", err)
	}
}
