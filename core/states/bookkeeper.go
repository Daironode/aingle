
package states

import (
	"io"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
)

type BookkeeperState struct {
	StateBase
	CurrBookkeeper []keypair.PublicKey
	NextBookkeeper []keypair.PublicKey
}

func (this *BookkeeperState) Serialization(sink *common.ZeroCopySink) {
	this.StateBase.Serialization(sink)
	sink.WriteUint32(uint32(len(this.CurrBookkeeper)))
	for _, v := range this.CurrBookkeeper {
		buf := keypair.SerializePublicKey(v)
		sink.WriteVarBytes(buf)
	}
	sink.WriteUint32(uint32(len(this.NextBookkeeper)))
	for _, v := range this.NextBookkeeper {
		buf := keypair.SerializePublicKey(v)
		sink.WriteVarBytes(buf)
	}
}

func (this *BookkeeperState) Deserialization(source *common.ZeroCopySource) error {
	err := this.StateBase.Deserialization(source)
	if err != nil {
		return err
	}
	n, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < int(n); i++ {
		buf, _, irregular, eof := source.NextVarBytes()
		if irregular {
			return common.ErrIrregularData
		}
		if eof {
			return io.ErrUnexpectedEOF
		}
		key, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.CurrBookkeeper = append(this.CurrBookkeeper, key)
	}
	n, eof = source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < int(n); i++ {
		buf, _, irregular, eof := source.NextVarBytes()
		if irregular {
			return common.ErrIrregularData
		}
		if eof {
			return io.ErrUnexpectedEOF
		}
		key, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.NextBookkeeper = append(this.NextBookkeeper, key)
	}
	return nil
}

func (v *BookkeeperState) ToArray() []byte {
	return common.SerializeToBytes(v)
}
