 package common

import (
	"fmt"

	" github.com/Daironode/aingle/common"
)

type ToMerkleValue struct {
	TxHash      []byte
	FromChainID uint64
	MakeTxParam *MakeTxParam
}

func (this *ToMerkleValue) Deserialization(source *common.ZeroCopySource) error {
	txHash, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MerkleValue deserialize txHash error")
	}
	fromChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MerkleValue deserialize fromChainID error")
	}

	makeTxParam := new(MakeTxParam)
	err := makeTxParam.Deserialization(source)
	if err != nil {
		return fmt.Errorf("MerkleValue deserialize makeTxParam error:%s", err)
	}

	this.TxHash = txHash
	this.FromChainID = fromChainID
	this.MakeTxParam = makeTxParam
	return nil
}

type MakeTxParam struct {
	TxHash              []byte
	CrossChainID        []byte
	FromContractAddress []byte
	ToChainID           uint64
	ToContractAddress   []byte
	Method              string
	Args                []byte
}

func (this *MakeTxParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.TxHash)
	sink.WriteVarBytes(this.CrossChainID)
	sink.WriteVarBytes(this.FromContractAddress)
	sink.WriteUint64(this.ToChainID)
	sink.WriteVarBytes(this.ToContractAddress)
	sink.WriteVarBytes([]byte(this.Method))
	sink.WriteVarBytes(this.Args)
}

func (this *MakeTxParam) Deserialization(source *common.ZeroCopySource) error {
	txHash, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize txHash error")
	}
	crossChainID, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize crossChainID error")
	}
	fromContractAddress, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize fromContractAddress error")
	}
	toChainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("MakeTxParam deserialize toChainID error")
	}
	toContractAddress, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize toContractAddress error")
	}
	method, _, irr, eof := source.NextString()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize method error")
	}
	args, _, irr, eof := source.NextVarBytes()
	if eof || irr {
		return fmt.Errorf("MakeTxParam deserialize args error")
	}

	this.TxHash = txHash
	this.CrossChainID = crossChainID
	this.FromContractAddress = fromContractAddress
	this.ToChainID = toChainID
	this.ToContractAddress = toContractAddress
	this.Method = method
	this.Args = args
	return nil
}
