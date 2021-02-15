
package cross_chain_manager

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type CreateCrossChainTxParam struct {
	ToChainID         uint64
	ToContractAddress []byte
	Method            string
	Args              []byte
}

func (this *CreateCrossChainTxParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.ToChainID)
	utils.EncodeVarBytes(sink, this.ToContractAddress)
	utils.EncodeString(sink, this.Method)
	utils.EncodeVarBytes(sink, this.Args)
}

func (this *CreateCrossChainTxParam) Deserialization(source *common.ZeroCopySource) error {
	toChainID, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("CreateCrossChainTxParam deserialize toChainID error:%s", err)
	}
	toContractAddress, err := utils.DecodeVarBytes(source)
	if err != nil {
		return fmt.Errorf("CreateCrossChainTxParam deserialize toContractAddress error:%s", err)
	}
	method, err := utils.DecodeString(source)
	if err != nil {
		return fmt.Errorf("CreateCrossChainTxParam deserialize method error:%s", err)
	}
	args, err := utils.DecodeVarBytes(source)
	if err != nil {
		return fmt.Errorf("CreateCrossChainTxParam deserialize args error:%s", err)
	}

	this.ToChainID = toChainID
	this.ToContractAddress = toContractAddress
	this.Method = method
	this.Args = args
	return nil
}

type ProcessCrossChainTxParam struct {
	Address     common.Address
	FromChainID uint64
	Height      uint32
	Proof       string
	Header      []byte
}

func (this *ProcessCrossChainTxParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeAddress(sink, this.Address)
	utils.EncodeVarUint(sink, this.FromChainID)
	utils.EncodeVarUint(sink, uint64(this.Height))
	utils.EncodeString(sink, this.Proof)
	utils.EncodeVarBytes(sink, this.Header)
}

func (this *ProcessCrossChainTxParam) Deserialization(source *common.ZeroCopySource) error {
	address, err := utils.DecodeAddress(source)
	if err != nil {
		return fmt.Errorf("ProcessCrossChainTxParam deserialize address error:%s", err)
	}
	fromChainID, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("ProcessCrossChainTxParam deserialize fromChainID error:%s", err)
	}
	height, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("ProcessCrossChainTxParam deserialize height error:%s", err)
	}
	proof, err := utils.DecodeString(source)
	if err != nil {
		return fmt.Errorf("ProcessCrossChainTxParam deserialize proof error:%s", err)
	}
	header, err := utils.DecodeVarBytes(source)
	if err != nil {
		return fmt.Errorf("ProcessCrossChainTxParam deserialize header error:%s", err)
	}
	this.Address = address
	this.FromChainID = fromChainID
	this.Height = uint32(height)
	this.Proof = proof
	this.Header = header
	return nil
}

type OngUnlockParam struct {
	FromChainID uint64
	Address     common.Address
	Amount      uint64
}

func (this *OngUnlockParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.FromChainID)
	utils.EncodeAddress(sink, this.Address)
	utils.EncodeVarUint(sink, this.Amount)
}

func (this *OngUnlockParam) Deserialization(source *common.ZeroCopySource) error {
	fromChainID, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("OngLockParam deserialize fromChainID error:%s", err)
	}
	address, err := utils.DecodeAddress(source)
	if err != nil {
		return fmt.Errorf("OngLockParam deserialize address error:%s", err)
	}
	amount, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("OngLockParam deserialize amount error:%s", err)
	}
	this.FromChainID = fromChainID
	this.Address = address
	this.Amount = amount
	return nil
}
