
package states

import (
	"io"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
)

type ValidatorState struct {
	StateBase
	PublicKey keypair.PublicKey
}

func (this *ValidatorState) Serialization(sink *common.ZeroCopySink) {
	this.StateBase.Serialization(sink)
	buf := keypair.SerializePublicKey(this.PublicKey)
	sink.WriteVarBytes(buf)
}

func (this *ValidatorState) Deserialization(source *common.ZeroCopySource) error {
	err := this.StateBase.Deserialization(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ValidatorState], StateBase Deserialize failed.")
	}
	buf, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return errors.NewDetailErr(common.ErrIrregularData, errors.ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	if eof {
		return errors.NewDetailErr(io.ErrUnexpectedEOF, errors.ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	pk, err := keypair.DeserializePublicKey(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	this.PublicKey = pk
	return nil
}
