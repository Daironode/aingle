
package states

import (
	"io"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/event"
)

// Invoke smart contract struct
// Param Version: invoke smart contract version, default 0
// Param Address: invoke on blockchain smart contract by address
// Param Method: invoke smart contract method, default ""
// Param Args: invoke smart contract arguments
type ContractInvokeParam struct {
	Version byte
	Address common.Address
	Method  string
	Args    []byte
}

func (this *ContractInvokeParam) Serialization(sink *common.ZeroCopySink) {
	sink.WriteByte(this.Version)
	sink.WriteAddress(this.Address)
	sink.WriteVarBytes([]byte(this.Method))
	sink.WriteVarBytes([]byte(this.Args))
}

// `ContractInvokeParam.Args` has reference of `source`
func (this *ContractInvokeParam) Deserialization(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.Version, eof = source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Address, eof = source.NextAddress()
	if eof {
		return io.ErrUnexpectedEOF
	}
	var method []byte
	method, _, irregular, eof = source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}
	this.Method = string(method)

	this.Args, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}

type PreExecResult struct {
	State  byte
	Gas    uint64
	Result interface{}
	Notify []*event.NotifyEventInfo
}
