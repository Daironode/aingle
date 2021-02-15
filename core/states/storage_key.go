 
package states

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type StorageKey struct {
	ContractAddress common.Address
	Key             []byte
}

func (this *StorageKey) Serialization(sink *common.ZeroCopySink) {
	this.ContractAddress.Serialization(sink)
	sink.WriteVarBytes(this.Key)
}

func (this *StorageKey) Deserialization(source *common.ZeroCopySource) error {
	if err := this.ContractAddress.Deserialization(source); err != nil {
		return err
	}
	key, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Key = key
	return nil
}

func (this *StorageKey) ToArray() []byte {
	return common.SerializeToBytes(this)
}
