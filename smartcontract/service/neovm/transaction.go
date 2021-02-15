
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/core/types"
	vm " github.com/Daironode/aingle/vm/neovm"
	vmtypes " github.com/Daironode/aingle/vm/neovm/types"
)

// GetExecutingAddress push transaction's hash to vm stack
func TransactionGetHash(service *NeoVmService, engine *vm.Executor) error {
	txn, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return fmt.Errorf("[TransactionGetHash] PopAsInteropValue error:%s", err)
	}
	if tx, ok := txn.Data.(*types.Transaction); ok {
		txHash := tx.Hash()
		return engine.EvalStack.PushBytes(txHash.ToArray())
	}
	return fmt.Errorf("[TransactionGetHash] Type error")
}

// TransactionGetType push transaction's type to vm stack
func TransactionGetType(service *NeoVmService, engine *vm.Executor) error {
	txn, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return fmt.Errorf("[TransactionGetType] PopAsInteropValue error:%s", err)
	}
	if tx, ok := txn.Data.(*types.Transaction); ok {
		return engine.EvalStack.PushInt64(int64(tx.TxType))
	}
	return fmt.Errorf("[TransactionGetType] Type error")
}

// TransactionGetAttributes push transaction's attributes to vm stack
func TransactionGetAttributes(service *NeoVmService, engine *vm.Executor) error {
	_, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return fmt.Errorf("[TransactionGetAttributes] PopAsInteropValue error: %s", err)
	}
	attributList := make([]vmtypes.VmValue, 0)
	return engine.EvalStack.PushAsArray(attributList)
}
