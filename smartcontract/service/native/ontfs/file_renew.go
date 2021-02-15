
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type FileReNew struct {
	FileHash       []byte
	FileOwner      common.Address
	Payer          common.Address
	NewTimeExpired uint64
}

type FileReNewList struct {
	FilesReNew []FileReNew
}

func (this *FileReNew) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
	utils.EncodeAddress(sink, this.FileOwner)
	utils.EncodeAddress(sink, this.Payer)
	utils.EncodeVarUint(sink, this.NewTimeExpired)
}

func (this *FileReNew) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.FileOwner, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.Payer, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.NewTimeExpired, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}

func (this *FileReNewList) Serialization(sink *common.ZeroCopySink) {
	fileReNewCount := uint64(len(this.FilesReNew))
	utils.EncodeVarUint(sink, fileReNewCount)

	for _, fileReNew := range this.FilesReNew {
		sinkTmp := common.NewZeroCopySink(nil)
		fileReNew.Serialization(sinkTmp)
		sink.WriteVarBytes(sinkTmp.Bytes())
	}
}

func (this *FileReNewList) Deserialization(source *common.ZeroCopySource) error {
	fileReNewCount, err := utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	for i := uint64(0); i < fileReNewCount; i++ {
		fileReNewTmp, err := DecodeVarBytes(source)
		if err != nil {
			return err
		}

		var fileReNew FileReNew
		src := common.NewZeroCopySource(fileReNewTmp)
		if err = fileReNew.Deserialization(src); err != nil {
			return err
		}
		this.FilesReNew = append(this.FilesReNew, fileReNew)
	}
	return nil
}
