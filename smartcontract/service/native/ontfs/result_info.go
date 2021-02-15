
package ontfs

import (
	"fmt"

	" github.com/Daironode/aingle/common"
)

type RetInfo struct {
	Ret  bool
	Info []byte
}

func (this *RetInfo) Serialization(sink *common.ZeroCopySink) {
	sink.WriteBool(this.Ret)
	sink.WriteVarBytes(this.Info)
}

func (this *RetInfo) Deserialization(source *common.ZeroCopySource) error {
	var err error
	if this.Ret, err = DecodeBool(source); err != nil {
		return fmt.Errorf("[RetInfo] [Ret] Deserialization from error:%v", err)
	}
	if this.Info, err = DecodeVarBytes(source); err != nil {
		return fmt.Errorf("[RetInfo] [Info] Deserialization from error:%v", err)
	}
	return nil
}

func EncRet(ret bool, info []byte) []byte {
	retInfo := RetInfo{ret, info}
	sink := common.NewZeroCopySink(nil)
	retInfo.Serialization(sink)
	return sink.Bytes()
}

func DecRet(ret []byte) (*RetInfo, error) {
	var retInfo RetInfo
	source := common.NewZeroCopySource(ret)
	err := retInfo.Deserialization(source)
	return &retInfo, err
}
