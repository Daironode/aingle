
package states

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type WasmContractParam struct {
	Address common.Address
	Args    []byte
}

func (this *WasmContractParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteAddress(this.Address)
	sink.WriteVarBytes([]byte(this.Args))
}

// `ContractInvokeParam.Args` has reference of `source`
func (this *WasmContractParam) Deserialization(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.Address, eof = source.NextAddress()
	if eof {
		return io.ErrUnexpectedEOF
	}

	this.Args, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}
