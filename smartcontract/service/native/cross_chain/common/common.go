 package common

import (
	"fmt"
	"math/big"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
	ctypes " github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/neovm"
	ntypes " github.com/Daironode/aingle/vm/neovm/types"
)

func CrossChainNeoVMCall(this *native.NativeService, address common.Address, method string, args []byte,
	fromContractAddress []byte, fromChainID uint64) (interface{}, error) {
	dep, err := this.CacheDB.GetContract(address)
	if err != nil {
		return nil, errors.NewErr("[NeoVMCall] Get contract context error!")
	}
	log.Debugf("[NeoVMCall] native invoke neovm contract address:%s", address.ToHexString())
	if dep == nil {
		return nil, errors.NewErr("[NeoVMCall] native invoke neovm contract is nil")
	}
	m, err := ntypes.VmValueFromBytes([]byte(method))
	if err != nil {
		return nil, err
	}
	array := ntypes.NewArrayValue()
	a, err := ntypes.VmValueFromBytes(args)
	if err != nil {
		return nil, err
	}
	if err := array.Append(a); err != nil {
		return nil, err
	}
	fca, err := ntypes.VmValueFromBytes(fromContractAddress)
	if err != nil {
		return nil, err
	}
	if err := array.Append(fca); err != nil {
		return nil, err
	}
	fci, err := ntypes.VmValueFromBigInt(new(big.Int).SetUint64(fromChainID))
	if err != nil {
		return nil, err
	}
	if err := array.Append(fci); err != nil {
		return nil, err
	}
	if !this.ContextRef.CheckUseGas(neovm.NATIVE_INVOKE_GAS) {
		return nil, fmt.Errorf("[CrossChainNeoVMCall], check use gaslimit insufficient！")
	}
	engine, err := this.ContextRef.NewExecuteEngine(dep.GetRawCode(), ctypes.InvokeNeo)
	if err != nil {
		return nil, err
	}
	evalStack := engine.(*neovm.NeoVmService).Engine.EvalStack
	if err := evalStack.Push(ntypes.VmValueFromArrayVal(array)); err != nil {
		return nil, err
	}
	if err := evalStack.Push(m); err != nil {
		return nil, err
	}
	return engine.Invoke()
}
