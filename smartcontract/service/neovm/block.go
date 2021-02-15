 package neovm

import (
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
	vmtypes " github.com/Daironode/aingle/vm/neovm/types"
)

// BlockGetTransactionCount put block's transactions count to vm stack
func BlockGetTransactionCount(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		val := vmtypes.VmValueFromInt64(int64(len(block.Transactions)))
		return engine.EvalStack.Push(val)
	}
	return errors.NewErr("[BlockGetTransactionCount] Wrong type ")
}

// BlockGetTransactions put block's transactions to vm stack
func BlockGetTransactions(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		transactions := block.Transactions
		transactionList := make([]vmtypes.VmValue, 0)

		for _, v := range transactions {
			transactionList = append(transactionList, vmtypes.VmValueFromInteropValue(vmtypes.NewInteropValue(v)))
		}

		return engine.EvalStack.PushAsArray(transactionList)
	}
	return errors.NewErr("[BlockGetTransactions] Wrong type ")
}

// BlockGetTransaction put block's transaction to vm stack
func BlockGetTransaction(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	index, err := engine.EvalStack.PopAsInt64()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		if index < 0 || int(index) >= len(block.Transactions) {
			return errors.NewErr("[BlockGetTransaction] index out of bounds")
		}
		return engine.EvalStack.PushAsInteropValue(block.Transactions[index])

	}
	return errors.NewErr("[BlockGetTransaction] Wrong type ")
}
