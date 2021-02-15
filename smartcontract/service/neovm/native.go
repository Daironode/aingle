
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/states"
	vm " github.com/Daironode/aingle/vm/neovm"
)

func NativeInvoke(service *NeoVmService, engine *vm.Executor) error {
	version, err := engine.EvalStack.PopAsInt64()
	if err != nil {
		return err
	}
	address, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("invoke native contract:%s, address invalid", address)
	}
	method, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	if len(method) > METHOD_LENGTH_LIMIT {
		return fmt.Errorf("invoke native contract:%s method:%s too long, over max length 1024 limit", address, method)
	}
	args, err := engine.EvalStack.Pop()
	if err != nil {
		return err
	}
	sink := new(common.ZeroCopySink)
	if err := args.BuildParamToNative(sink); err != nil {
		return err
	}

	contract := states.ContractInvokeParam{
		Version: byte(version),
		Address: addr,
		Method:  string(method),
		Args:    sink.Bytes(),
	}

	nat := &native.NativeService{
		Store:       service.Store,
		CacheDB:     service.CacheDB,
		InvokeParam: contract,
		Tx:          service.Tx,
		Height:      service.Height,
		Time:        service.Time,
		ContextRef:  service.ContextRef,
		ServiceMap:  make(map[string]native.Handler),
		PreExec:     service.PreExec,
	}

	result, err := nat.Invoke()
	if err != nil {
		return err
	}
	return engine.EvalStack.PushBytes(result)
}
