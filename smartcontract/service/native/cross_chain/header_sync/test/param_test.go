
package test

import (
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/header_sync"
	"github.com/stretchr/testify/assert"
)

func TestSyncBlockHeaderParam(t *testing.T) {
	param := header_sync.SyncBlockHeaderParam{
		Address: common.ADDRESS_EMPTY,
		Headers: [][]byte{{1}, {2}, {3}},
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	var p header_sync.SyncBlockHeaderParam
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, p, param)
}

func TestSyncGenesisHeaderParam(t *testing.T) {
	param := header_sync.SyncGenesisHeaderParam{
		GenesisHeader: []byte{1, 2, 3},
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	var p header_sync.SyncGenesisHeaderParam
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, p, param)
}
