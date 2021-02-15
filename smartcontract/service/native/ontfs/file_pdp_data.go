 
package ontfs

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/ontfs/pdp"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type PdpData struct {
	NodeAddr        common.Address
	FileHash        []byte
	ProveData       []byte
	ChallengeHeight uint64
}

func (this *PdpData) Serialization(sink *common.ZeroCopySink) error {
	if len(this.ProveData) < pdp.VersionLength {
		return fmt.Errorf("PdpData Serialization error: ProveData length shorter than 8")
	}
	utils.EncodeAddress(sink, this.NodeAddr)
	sink.WriteVarBytes(this.FileHash)
	sink.WriteVarBytes(this.ProveData)
	utils.EncodeVarUint(sink, this.ChallengeHeight)
	return nil
}

func (this *PdpData) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.NodeAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.ProveData, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.ChallengeHeight, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}
