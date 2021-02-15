 
package states

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type StateBase struct {
	StateVersion byte
}

func (this *StateBase) Serialization(sink *common.ZeroCopySink) {
	sink.WriteByte(this.StateVersion)
}

func (this *StateBase) Deserialization(source *common.ZeroCopySource) error {
	stateVersion, eof := source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.StateVersion = stateVersion
	return nil
}
