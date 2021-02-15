
package test

import (
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/cross_chain_manager"
	"github.com/stretchr/testify/assert"
)

func TestCreateCrossChainTxParam(t *testing.T) {
	param := cross_chain_manager.CreateCrossChainTxParam{
		ToChainID:         1,
		ToContractAddress: []byte{1, 2, 3, 4},
		Method:            "test",
		Args:              []byte{1, 2, 3, 4},
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	var p cross_chain_manager.CreateCrossChainTxParam
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, p, param)
}

func TestProcessCrossChainTxParam(t *testing.T) {
	param := cross_chain_manager.ProcessCrossChainTxParam{
		Address:     common.ADDRESS_EMPTY,
		FromChainID: 1,
		Height:      2,
		Proof:       "test",
		Header:      []byte{1, 2, 3, 4},
	}

	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	var p cross_chain_manager.ProcessCrossChainTxParam
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, param, p)
}

func TestOngUnlockParam(t *testing.T) {
	param := cross_chain_manager.OngUnlockParam{
		FromChainID: 1,
		Address:     common.ADDRESS_EMPTY,
		Amount:      1,
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	var p cross_chain_manager.OngUnlockParam
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, param, p)
}
