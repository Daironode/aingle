
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type FileReadSettleSlice struct {
	FileHash     []byte
	PayFrom      common.Address
	PayTo        common.Address
	SliceId      uint64
	PledgeHeight uint64
	Sig          []byte
	PubKey       []byte
}

func (this *FileReadSettleSlice) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
	utils.EncodeAddress(sink, this.PayFrom)
	utils.EncodeAddress(sink, this.PayTo)
	utils.EncodeVarUint(sink, this.SliceId)
	utils.EncodeVarUint(sink, this.PledgeHeight)
	sink.WriteVarBytes(this.Sig)
	sink.WriteVarBytes(this.PubKey)
}

func (this *FileReadSettleSlice) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.PayFrom, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.PayTo, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.SliceId, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.PledgeHeight, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.Sig, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.PubKey, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	return nil
}
