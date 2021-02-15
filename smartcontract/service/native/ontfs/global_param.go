
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type FsGlobalParam struct {
	MinTimeForFileStorage uint64
	ContractInvokeGasFee  uint64
	ChallengeReward       uint64
	FilePerServerPdpTimes uint64
	PassportExpire        uint64
	ChallengeInterval     uint64
	NodeMinVolume         uint64 //min total volume with fsNode
	NodePerKbPledge       uint64 //fsNode's pledge for participant
	FeePerBlockForRead    uint64 //cost for ontfs-sdk read from fsNode
	FilePerBlockFeeRate   uint64 //cost for ontfs-sdk save from fsNode
	SpacePerBlockFeeRate  uint64 //cost for ontfs-sdk save from fsNode
}

func (this *FsGlobalParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.MinTimeForFileStorage)
	utils.EncodeVarUint(sink, this.ContractInvokeGasFee)
	utils.EncodeVarUint(sink, this.ChallengeReward)
	utils.EncodeVarUint(sink, this.FilePerServerPdpTimes)
	utils.EncodeVarUint(sink, this.PassportExpire)
	utils.EncodeVarUint(sink, this.ChallengeInterval)
	utils.EncodeVarUint(sink, this.NodeMinVolume)
	utils.EncodeVarUint(sink, this.NodePerKbPledge)
	utils.EncodeVarUint(sink, this.FeePerBlockForRead)
	utils.EncodeVarUint(sink, this.FilePerBlockFeeRate)
	utils.EncodeVarUint(sink, this.SpacePerBlockFeeRate)
}

func (this *FsGlobalParam) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.MinTimeForFileStorage, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.ContractInvokeGasFee, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.ChallengeReward, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FilePerServerPdpTimes, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.PassportExpire, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.NodeMinVolume, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.NodePerKbPledge, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FeePerBlockForRead, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FilePerBlockFeeRate, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.SpacePerBlockFeeRate, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return err
}

func setGlobalParam(native *native.NativeService, globalParam *FsGlobalParam) {
	globalParamKey := GenGlobalParamKey(native.ContextRef.CurrentContext().ContractAddress)
	sink := common.NewZeroCopySink(nil)
	globalParam.Serialization(sink)
	utils.PutBytes(native, globalParamKey, sink.Bytes())
}

func getGlobalParam(native *native.NativeService) (*FsGlobalParam, error) {
	var globalParam FsGlobalParam

	globalParamKey := GenGlobalParamKey(native.ContextRef.CurrentContext().ContractAddress)
	item, err := utils.GetStorageItem(native, globalParamKey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getGlobalParam GetStorageItem error!")
	}
	if item == nil {
		globalParam = FsGlobalParam{
			MinTimeForFileStorage: DefaultMinTimeForFileStorage,
			ContractInvokeGasFee:  DefaultContractInvokeGasFee,
			ChallengeReward:       DefaultChallengeReward,
			FilePerServerPdpTimes: DefaultFilePerServerPdpTimes,
			PassportExpire:        DefaultPassportExpire,
			ChallengeInterval:     DefaultChallengeInterval,
			NodeMinVolume:         DefaultNodeMinVolume,
			NodePerKbPledge:       DefaultNodePerKbPledge,
			FeePerBlockForRead:    DefaultGasPerBlockForRead,
			FilePerBlockFeeRate:   DefaultFilePerBlockFeeRate,
			SpacePerBlockFeeRate:  DefaultSpacePerBlockFeeRate,
		}
		return &globalParam, nil
	}

	source := common.NewZeroCopySource(item.Value)
	if err := globalParam.Deserialization(source); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getGlobalParam Deserialization error!")
	}
	return &globalParam, nil
}
