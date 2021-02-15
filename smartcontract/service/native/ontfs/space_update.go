
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type SpaceUpdate struct {
	SpaceOwner     common.Address
	Payer          common.Address
	NewVolume      uint64
	NewTimeExpired uint64
}

func (this *SpaceUpdate) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeAddress(sink, this.SpaceOwner)
	utils.EncodeAddress(sink, this.Payer)
	utils.EncodeVarUint(sink, this.NewVolume)
	utils.EncodeVarUint(sink, this.NewTimeExpired)
}

func (this *SpaceUpdate) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.SpaceOwner, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.Payer, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.NewVolume, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.NewTimeExpired, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}
