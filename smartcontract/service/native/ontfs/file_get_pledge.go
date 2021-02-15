
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type GetReadPledge struct {
	FileHash   []byte
	Downloader common.Address
}

func (this *GetReadPledge) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
	utils.EncodeAddress(sink, this.Downloader)
}

func (this *GetReadPledge) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.Downloader, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	return nil
}
