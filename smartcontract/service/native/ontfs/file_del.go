
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type FileDel struct {
	FileHash []byte
}

type FileDelList struct {
	FilesDel []FileDel
}

func (this *FileDel) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
}

func (this *FileDel) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	return nil
}

func (this *FileDelList) Serialization(sink *common.ZeroCopySink) {
	fileDelCount := uint64(len(this.FilesDel))
	utils.EncodeVarUint(sink, fileDelCount)

	for _, fileDel := range this.FilesDel {
		sinkTmp := common.NewZeroCopySink(nil)
		fileDel.Serialization(sinkTmp)
		sink.WriteVarBytes(sinkTmp.Bytes())
	}
}

func (this *FileDelList) Deserialization(source *common.ZeroCopySource) error {
	fileDelCount, err := utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	for i := uint64(0); i < fileDelCount; i++ {
		fileDelTmp, err := DecodeVarBytes(source)
		if err != nil {
			return err
		}

		var fileDel FileDel
		src := common.NewZeroCopySource(fileDelTmp)
		if err = fileDel.Deserialization(src); err != nil {
			return err
		}
		this.FilesDel = append(this.FilesDel, fileDel)
	}
	return nil
}
