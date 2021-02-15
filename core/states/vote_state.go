
package states

import (
	"io"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
)

type VoteState struct {
	StateBase
	PublicKeys []keypair.PublicKey
	Count      common.Fixed64
}

func (this *VoteState) Serialization(sink *common.ZeroCopySink) {
	this.StateBase.Serialization(sink)
	sink.WriteUint32(uint32(len(this.PublicKeys)))
	for _, v := range this.PublicKeys {
		buf := keypair.SerializePublicKey(v)
		sink.WriteVarBytes(buf)
	}
	sink.WriteUint64(uint64(this.Count))
}

func (this *VoteState) Deserialization(source *common.ZeroCopySource) error {
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
		pk, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.PublicKeys = append(this.PublicKeys, pk)
	}
	c, eof := source.NextUint64()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Count = common.Fixed64(int64(c))
	return nil
}
