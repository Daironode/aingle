
package neovm

import (
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
	vmtypes " github.com/Daironode/aingle/vm/neovm/types"
)

// AttributeGetUsage put attribute's usage to vm stack
func AttributeGetUsage(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if a, ok := i.Data.(*types.TxAttribute); ok {
		return engine.EvalStack.PushInt64(int64(a.Usage))
	}
	return errors.NewErr("[AttributeGetUsage] Wrong type!")
}

// AttributeGetData put attribute's data to vm stack
func AttributeGetData(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if a, ok := i.Data.(*types.TxAttribute); ok {
		val, err := vmtypes.VmValueFromBytes(a.Data)
		if err != nil {
			return err
		}
		return engine.EvalStack.Push(val)
	}
	return errors.NewErr("[AttributeGetData] Wrong type!")
}
