
package neovm

import (
	"io"

	" github.com/Daironode/aingle/vm/neovm/utils"
)

type ExecutionContext struct {
	Code               []byte
	OpReader           *utils.VmReader
	InstructionPointer int
	vmFlags            VmFeatureFlag
}

func NewExecutionContext(code []byte, flag VmFeatureFlag) *ExecutionContext {
	var context ExecutionContext
	context.Code = code
	context.OpReader = utils.NewVmReader(code)
	context.OpReader.AllowEOF = flag.AllowReaderEOF
	context.vmFlags = flag

	context.InstructionPointer = 0
	return &context
}

func (ec *ExecutionContext) GetInstructionPointer() int {
	return ec.OpReader.Position()
}

func (ec *ExecutionContext) SetInstructionPointer(offset int64) error {
	_, err := ec.OpReader.Seek(offset, io.SeekStart)
	return err
}

func (ec *ExecutionContext) NextInstruction() OpCode {
	return OpCode(ec.Code[ec.OpReader.Position()])
}

func (self *ExecutionContext) ReadOpCode() (val OpCode, eof bool) {
	code, err := self.OpReader.ReadByte()
	if err != nil {
		eof = true
		return
	}
	val = OpCode(code)
	return val, false
}

func (ec *ExecutionContext) Clone() *ExecutionContext {
	context := NewExecutionContext(ec.Code, ec.vmFlags)
	context.InstructionPointer = ec.InstructionPointer
	_ = context.SetInstructionPointer(int64(ec.GetInstructionPointer()))
	return context
}
