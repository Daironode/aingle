
package states

import (
	"io"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
)

type StorageItem struct {
	StateBase
	Value []byte
}

func (this *StorageItem) Serialization(sink *common.ZeroCopySink) {
	this.StateBase.Serialization(sink)
	sink.WriteVarBytes(this.Value)
}

func (this *StorageItem) Deserialization(source *common.ZeroCopySource) error {
	err := this.StateBase.Deserialization(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageItem], StateBase Deserialize failed.")
	}
	value, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return errors.NewDetailErr(common.ErrIrregularData, errors.ErrNoCode, "[StorageItem], Value Deserialize failed.")
	}
	if eof {
		return errors.NewDetailErr(io.ErrUnexpectedEOF, errors.ErrNoCode, "[StorageItem], Value Deserialize failed.")
	}
	this.Value = value
	return nil
}

func (storageItem *StorageItem) ToArray() []byte {
	return common.SerializeToBytes(storageItem)
}

func GetValueFromRawStorageItem(raw []byte) ([]byte, error) {
	item := StorageItem{}
	err := item.Deserialization(common.NewZeroCopySource(raw))
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func GenRawStorageItem(value []byte) []byte {
	item := StorageItem{}
	item.Value = value
	return item.ToArray()
}
