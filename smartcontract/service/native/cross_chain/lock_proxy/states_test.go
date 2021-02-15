 
package lock_proxy

import (
	"encoding/hex"
	"math/big"
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/constants"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestLockParam_Serialize(t *testing.T) {
	fromAddr, _ := common.AddressFromBase58("709c937270e1d5a490718a2b4a230186bdd06a01")
	toAddrBs, _ := hex.DecodeString("709c937270e1d5a490718a2b4a230186bdd06a02")
	param := LockParam{
		SourceAssetHash: utils.OntContractAddress,
		ToChainID:       0,
		FromAddress:     fromAddr,
		ToAddress:       toAddrBs,
		Value:           1,
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	param2 := LockParam{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := param2.Deserialization(source); err != nil {
		t.Fatal("LockParam deserialize fail!")
	}
	assert.Equal(t, param, param2)
}

func TestUnlockParam_Serialize(t *testing.T) {
	param := UnlockParam{
		ArgsBs:             []byte{1, 2, 3, 0, 100},
		FromContractHashBs: utils.OntContractAddress[:],
		FromChainId:        2,
	}
	sink := common.NewZeroCopySink(nil)
	param.Serialization(sink)

	param2 := UnlockParam{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := param2.Deserialization(source); err != nil {
		t.Fatal("LockParam deserialize fail!")
	}
	assert.Equal(t, param, param2)
}

func TestArgs_Serialize(t *testing.T) {
	toAddr, _ := hex.DecodeString("709c937270e1d5a490718a2b4a230186bdd06a02")
	args := Args{
		TargetAssetHash: utils.OntContractAddress[:],
		ToAddress:       toAddr,
		Value:           100,
	}
	sink := common.NewZeroCopySink(nil)
	args.Serialization(sink)

	args2 := Args{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := args2.Deserialization(source); err != nil {
		t.Fatal("Args deserialize fail!")
	}
	assert.Equal(t, args, args2)
}

func TestBindAssetParam_Serialize(t *testing.T) {
	bindAssetParam := BindAssetParam{
		SourceAssetHash:    utils.OntContractAddress,
		TargetChainId:      uint64(0),
		TargetAssetHash:    utils.OntContractAddress[:],
		Limit:              big.NewInt(int64(constants.ONT_TOTAL_SUPPLY)),
		IsTargetChainAsset: false,
	}
	sink := common.NewZeroCopySink(nil)
	bindAssetParam.Serialization(sink)

	bindAssetParam2 := BindAssetParam{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := bindAssetParam2.Deserialization(source); err != nil {
		t.Fatal("BindAssetParam deserialize fail!")
	}
	assert.Equal(t, bindAssetParam, bindAssetParam2)
}
