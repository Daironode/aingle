
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type ReadPlan struct {
	NodeAddr         common.Address
	MaxReadBlockNum  uint64
	HaveReadBlockNum uint64
	NumOfSettlements uint64
}

type ReadPledge struct {
	FileHash   []byte
	Downloader common.Address
	RestMoney  uint64
	ReadPlans  []ReadPlan
}

func (this *ReadPlan) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeAddress(sink, this.NodeAddr)
	utils.EncodeVarUint(sink, this.MaxReadBlockNum)
	utils.EncodeVarUint(sink, this.HaveReadBlockNum)
	utils.EncodeVarUint(sink, this.NumOfSettlements)
}

func (this *ReadPlan) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.NodeAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.MaxReadBlockNum, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.HaveReadBlockNum, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.NumOfSettlements, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}

func (this *ReadPledge) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
	utils.EncodeAddress(sink, this.Downloader)
	utils.EncodeVarUint(sink, this.RestMoney)

	planCount := uint64(len(this.ReadPlans))
	utils.EncodeVarUint(sink, planCount)

	for _, readPlan := range this.ReadPlans {
		sinkTmp := common.NewZeroCopySink(nil)
		readPlan.Serialization(sinkTmp)
		sink.WriteVarBytes(sinkTmp.Bytes())
	}
}

func (this *ReadPledge) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.Downloader, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.RestMoney, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	planCount, err := utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	var readPlan ReadPlan
	for i := uint64(0); i < planCount; i++ {
		readPlanTmp, err := DecodeVarBytes(source)
		if err != nil {
			return err
		}
		src := common.NewZeroCopySource(readPlanTmp)
		if err = readPlan.Deserialization(src); err != nil {
			return err
		}
		this.ReadPlans = append(this.ReadPlans, readPlan)
	}
	return nil
}

func addReadPledge(native *native.NativeService, readPledge *ReadPledge) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	key := GenFsReadPledgeKey(contract, readPledge.Downloader, readPledge.FileHash)
	sink := common.NewZeroCopySink(nil)
	readPledge.Serialization(sink)
	utils.PutBytes(native, key, sink.Bytes())
}

func getRawReadPledge(native *native.NativeService, downLoader common.Address, fileHash []byte) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	key := GenFsReadPledgeKey(contract, downLoader, fileHash)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getRawReadPledge GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("getRawReadPledge not found!")
	}

	return item.Value, nil
}

func getReadPledge(native *native.NativeService, downLoader common.Address, fileHash []byte) (*ReadPledge, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	key := GenFsReadPledgeKey(contract, downLoader, fileHash)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getReadPledge GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("getReadPledge not found!")
	}

	var readPledge ReadPledge
	source := common.NewZeroCopySource(item.Value)
	err = readPledge.Deserialization(source)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getReadPledge Deserialization error!")
	}
	return &readPledge, nil
}

func delReadPledge(native *native.NativeService, downloader common.Address, fileHash []byte) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	key := GenFsReadPledgeKey(contract, downloader, fileHash)
	native.CacheDB.Delete(key)
}
