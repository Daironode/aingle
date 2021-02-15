
package payload

import (
	"io"

	" github.com/Daironode/aingle/common"
)

// InvokeCode is an implementation of transaction payload for invoke smartcontract
type InvokeCode struct {
	Code []byte
}

//note: InvokeCode.Code has data reference of param source
func (self *InvokeCode) Deserialization(source *common.ZeroCopySource) error {
	code, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	self.Code = code
	return nil
}

func (self *InvokeCode) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(self.Code)
}
