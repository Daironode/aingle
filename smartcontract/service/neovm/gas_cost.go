
package neovm

import (
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

func StoreGasCost(gasTable map[string]uint64, engine *vm.Executor) (uint64, error) {
	key, err := engine.EvalStack.PeekAsBytes(1)
	if err != nil {
		return 0, err
	}
	value, err := engine.EvalStack.PeekAsBytes(2)
	if err != nil {
		return 0, err
	}
	if putCost, ok := gasTable[STORAGE_PUT_NAME]; ok {
		return uint64((len(key)+len(value)-1)/1024+1) * putCost, nil
	} else {
		return uint64(0), errors.NewErr("[StoreGasCost] get STORAGE_PUT_NAME gas failed")
	}
}

func GasPrice(gasTable map[string]uint64, engine *vm.Executor, name string) (uint64, error) {
	switch name {
	case STORAGE_PUT_NAME:
		return StoreGasCost(gasTable, engine)
	default:
		if value, ok := gasTable[name]; ok {
			return value, nil
		}
		return OPCODE_GAS, nil
	}
}
