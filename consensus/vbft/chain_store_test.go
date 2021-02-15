
package vbft

import (
	"os"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	vconfig " github.com/Daironode/aingle/consensus/vbft/config"
	" github.com/Daironode/aingle/core/genesis"
	" github.com/Daironode/aingle/core/ledger"
)

var testBookkeeperAccounts []*account.Account

func newTestChainStore(t *testing.T) *ChainStore {
	log.InitLog(log.InfoLog, log.Stdout)
	var err error
	acct := account.NewAccount("SHA256withECDSA")
	if acct == nil {
		t.Fatalf("GetDefaultAccount error: acc is nil")
	}
	os.RemoveAll(config.DEFAULT_DATA_DIR)
	db, err := ledger.NewLedger(config.DEFAULT_DATA_DIR, 0)
	if err != nil {
		t.Fatalf("NewLedger error %s", err)
	}

	var bookkeepers []keypair.PublicKey
	if len(testBookkeeperAccounts) == 0 {
		for i := 0; i < 7; i++ {
			acc := account.NewAccount("")
			testBookkeeperAccounts = append(testBookkeeperAccounts, acc)
			bookkeepers = append(bookkeepers, acc.PublicKey)
		}
	}

	genesisConfig := config.DefConfig.Genesis

	// update peers in genesis
	for i, p := range genesisConfig.VBFT.Peers {
		if i > 0 && i <= len(testBookkeeperAccounts) {
			p.PeerPubkey = vconfig.PubkeyID(testBookkeeperAccounts[i-1].PublicKey)
		}
	}
	block, err := genesis.BuildGenesisBlock(bookkeepers, genesisConfig)
	if err != nil {
		t.Fatalf("BuildGenesisBlock error %s", err)
	}

	err = db.Init(bookkeepers, block)
	if err != nil {
		t.Fatalf("InitLedgerStoreWithGenesisBlock error %s", err)
	}
	chainstore, err := OpenBlockStore(db, nil)
	if err != nil {
		t.Fatalf("openblockstore failed: %v\n", err)
	}
	return chainstore
}

func cleanTestChainStore() {
	os.RemoveAll(config.DEFAULT_DATA_DIR)
	testBookkeeperAccounts = make([]*account.Account, 0)
}

func TestGetChainedBlockNum(t *testing.T) {
	chainstore := newTestChainStore(t)
	if chainstore == nil {
		t.Error("newChainStrore error")
		return
	}
	defer cleanTestChainStore()

	blocknum := chainstore.GetChainedBlockNum()
	t.Logf("TestGetChainedBlockNum :%d", blocknum)
}
