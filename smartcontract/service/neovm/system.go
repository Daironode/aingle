
package neovm

import (
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// GetCodeContainer push current transaction to vm stack
func GetCodeContainer(service *NeoVmService, engine *vm.Executor) error {
	return engine.EvalStack.PushAsInteropValue(service.Tx)
}

// GetExecutingAddress push current context to vm stack
func GetExecutingAddress(service *NeoVmService, engine *vm.Executor) error {
	context := service.ContextRef.CurrentContext()
	if context == nil {
		return errors.NewErr("Current context invalid")
	}
	return engine.EvalStack.PushBytes(context.ContractAddress[:])
}

// GetExecutingAddress push previous context to vm stack
func GetCallingAddress(service *NeoVmService, engine *vm.Executor) error {
	context := service.ContextRef.CallingContext()
	if context == nil {
		return errors.NewErr("Calling context invalid")
	}
	return engine.EvalStack.PushBytes(context.ContractAddress[:])
}

// GetExecutingAddress push entry call context to vm stack
func GetEntryAddress(service *NeoVmService, engine *vm.Executor) error {
	context := service.ContextRef.EntryContext()
	if context == nil {
		return errors.NewErr("Entry context invalid")
	}
	return engine.EvalStack.PushBytes(context.ContractAddress[:])
}
