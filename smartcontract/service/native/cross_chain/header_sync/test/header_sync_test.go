 
package test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/common"
	vconfig " github.com/Daironode/aingle/consensus/vbft/config"
	" github.com/Daironode/aingle/core/signature"
	" github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/core/store/leveldbstore"
	" github.com/Daironode/aingle/core/store/overlaydb"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/smartcontract/service/native"
	cccom " github.com/Daironode/aingle/smartcontract/service/native/cross_chain/common"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/header_sync"
	" github.com/Daironode/aingle/smartcontract/service/native/global_params"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
	" github.com/Daironode/aingle/smartcontract/storage"
	"github.com/stretchr/testify/assert"
)

var (
	acct *account.Account

	setAcct = func() {
		acct = account.NewAccount("")
	}

	generateSomeAcct = func() *account.Account {
		return account.NewAccount("")
	}

	getNativeFunc = func(args []byte, db *storage.CacheDB) *native.NativeService {
		store, _ := leveldbstore.NewMemLevelDBStore()
		if db == nil {
			db = storage.NewCacheDB(overlaydb.NewOverlayDB(store))
		}

		return &native.NativeService{
			CacheDB: db,
			Input:   args,
			ContextRef: &smartcontract.SmartContract{
				Config: &smartcontract.Config{
					Tx: &types.Transaction{
						SignedAddr: []common.Address{acct.Address},
					},
				},
			},
		}
	}

	getGenesisHeader = func() []byte {
		blkInfo := &vconfig.VbftBlockInfo{
			NewChainConfig: &vconfig.ChainConfig{
				Peers: []*vconfig.PeerConfig{
					{Index: 0, ID: hex.EncodeToString(keypair.SerializePublicKey(acct.PublicKey))},
				},
			},
		}
		payload, _ := json.Marshal(blkInfo)
		bd := &cccom.Header{
			Version:          0,
			Height:           0,
			ChainID:          0,
			Bookkeepers:      []keypair.PublicKey{acct.PublicKey},
			ConsensusPayload: payload,
			NextBookkeeper:   acct.Address,
		}
		hash := bd.Hash()
		sig, _ := signature.Sign(acct, hash[:])
		bd.SigData = [][]byte{sig}
		sink := common.NewZeroCopySink(nil)
		bd.Serialization(sink)

		return sink.Bytes()
	}

	getHeaders = func(n uint32) [][]byte {
		hdrs := make([][]byte, 0)

		blkInfo := &vconfig.VbftBlockInfo{
			NewChainConfig: &vconfig.ChainConfig{
				Peers: []*vconfig.PeerConfig{
					{Index: 0, ID: vconfig.PubkeyID(acct.PublicKey)},
				},
			},
		}
		payload, _ := json.Marshal(blkInfo)
		for i := uint32(1); i <= n; i++ {
			bd := &cccom.Header{
				Version:          0,
				Height:           i,
				ChainID:          0,
				Bookkeepers:      []keypair.PublicKey{acct.PublicKey},
				ConsensusPayload: payload,
				NextBookkeeper:   acct.Address,
			}

			hash := bd.Hash()
			sig, _ := signature.Sign(acct, hash[:])
			bd.SigData = [][]byte{sig}
			sink := common.NewZeroCopySink(nil)
			bd.Serialization(sink)
			hdrs = append(hdrs, sink.Bytes())
		}

		return hdrs
	}
)

func init() {
	setAcct()
}

func TestSyncGenesisHeader(t *testing.T) {
	// normal case: with peers
	sink := common.NewZeroCopySink(nil)
	p := &header_sync.SyncGenesisHeaderParam{
		GenesisHeader: getGenesisHeader(),
	}
	p.Serialization(sink)

	bf := common.NewZeroCopySink(nil)
	utils.EncodeAddress(bf, acct.Address)
	si := &states.StorageItem{Value: bf.Bytes()}

	ns := getNativeFunc(sink.Bytes(), nil)
	ns.CacheDB.Put(global_params.GenerateOperatorKey(utils.ParamContractAddress), si.ToArray())

	ok, err := header_sync.SyncGenesisHeader(ns)
	assert.NoError(t, err)
	assert.Equal(t, utils.BYTE_TRUE, ok, "wrong result")

	// wrong owner
	ns.ContextRef.(*smartcontract.SmartContract).Config.Tx.SignedAddr = []common.Address{generateSomeAcct().Address}
	ok, err = header_sync.SyncGenesisHeader(ns)
	assert.EqualError(t, err, "SyncGenesisHeader, checkWitness error: validateOwner, authentication failed!",
		"not the right error")
}

func TestSyncBlockHeader(t *testing.T) {
	// first, we need to sync genesis header
	sink := common.NewZeroCopySink(nil)
	p := &header_sync.SyncGenesisHeaderParam{
		GenesisHeader: getGenesisHeader(),
	}
	p.Serialization(sink)

	bf := common.NewZeroCopySink(nil)
	utils.EncodeAddress(bf, acct.Address)
	si := &states.StorageItem{Value: bf.Bytes()}

	ns := getNativeFunc(sink.Bytes(), nil)
	ns.CacheDB.Put(global_params.GenerateOperatorKey(utils.ParamContractAddress), si.ToArray())

	_, _ = header_sync.SyncGenesisHeader(ns)

	// 1. next to check normal case
	sink = common.NewZeroCopySink(nil)
	param := &header_sync.SyncBlockHeaderParam{
		Address: acct.Address,
		Headers: getHeaders(3),
	}
	param.Serialization(sink)

	ns.Input = sink.Bytes()
	ok, err := header_sync.SyncBlockHeader(ns)
	assert.NoError(t, err)
	assert.Equal(t, utils.BYTE_TRUE, ok, "wrong result")

	// 2.more case?
}
