 package neovm

import (
	"fmt"
	"reflect"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/core/utils"
	" github.com/Daironode/aingle/vm/crossvm_codec"
	vm " github.com/Daironode/aingle/vm/neovm"
)

//neovm contract call wasmvm contract
func WASMInvoke(service *NeoVmService, engine *vm.Executor) error {
	address, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}

	contractAddress, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("invoke wasm contract:%s, address invalid", address)
	}

	dp, err := service.CacheDB.GetContract(contractAddress)
	if err != nil {
		return err
	}
	if dp == nil {
		return fmt.Errorf("wasm contract does not exist")
	}

	if dp.VmType() != payload.WASMVM_TYPE {
		return fmt.Errorf("not a wasm contract")
	}

	parambytes, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	list, err := crossvm_codec.DeserializeCallParam(parambytes)
	if err != nil {
		return err
	}

	params, ok := list.([]interface{})
	if !ok {
		return fmt.Errorf("wasm invoke error: wrong param type:%s", reflect.TypeOf(list).String())
	}

	inputs, err := utils.BuildWasmVMInvokeCode(contractAddress, params)
	if err != nil {
		return err
	}

	newservice, err := service.ContextRef.NewExecuteEngine(inputs, types.InvokeWasm)
	if err != nil {
		return err
	}

	tmpRes, err := newservice.Invoke()
	if err != nil {
		return err
	}

	return engine.EvalStack.PushBytes(tmpRes.([]byte))
}
