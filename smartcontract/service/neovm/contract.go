
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// ContractCreate create a new smart contract on blockchain, and put it to vm stack
func ContractCreate(service *NeoVmService, engine *vm.Executor) error {
	contract, err := isContractParamValid(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractCreate] contract parameters invalid!")
	}
	contractAddress := contract.Address()
	dep, err := service.CacheDB.GetContract(contractAddress)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractCreate] GetOrAdd error!")
	}
	if dep == nil {
		service.CacheDB.PutContract(contract)
		dep = contract
	}
	return engine.EvalStack.PushAsInteropValue(dep)
}

// ContractMigrate migrate old smart contract to a new contract, and destroy old contract
func ContractMigrate(service *NeoVmService, engine *vm.Executor) error {
	contract, err := isContractParamValid(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractMigrate] contract parameters invalid!")
	}
	newAddr := contract.Address()

	if err := isContractExist(service, newAddr); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ContractMigrate] contract invalid!")
	}
	context := service.ContextRef.CurrentContext()
	oldAddr := context.ContractAddress

	service.CacheDB.PutContract(contract)
	service.CacheDB.DeleteContract(oldAddr)

	iter := service.CacheDB.NewIterator(oldAddr[:])
	for has := iter.First(); has; has = iter.Next() {
		key := iter.Key()
		val := iter.Value()

		newKey := genStorageKey(newAddr, key[20:])
		service.CacheDB.Put(newKey, val)
		service.CacheDB.Delete(key)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}
	return engine.EvalStack.PushAsInteropValue(contract)
}

// ContractDestory destroy a contract
func ContractDestory(service *NeoVmService, engine *vm.Executor) error {
	context := service.ContextRef.CurrentContext()
	if context == nil {
		return errors.NewErr("[ContractDestory] current contract context invalid!")
	}
	addr := context.ContractAddress
	contract, err := service.CacheDB.GetContract(addr)
	if err != nil || contract == nil {
		return errors.NewErr("[ContractDestory] get current contract fail!")
	}

	service.CacheDB.DeleteContract(addr)

	iter := service.CacheDB.NewIterator(addr[:])
	for has := iter.First(); has; has = iter.Next() {
		key := iter.Key()
		service.CacheDB.Delete(key)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}

	return nil
}

// ContractGetStorageContext put contract storage context to vm stack
func ContractGetStorageContext(service *NeoVmService, engine *vm.Executor) error {
	opInterface, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if opInterface.Data == nil {
		return errors.NewErr("[GetStorageContext] Pop data nil!")
	}
	contractState, ok := opInterface.Data.(*payload.DeployCode)
	if !ok {
		return errors.NewErr("[GetStorageContext] Pop data not contract!")
	}
	address := contractState.Address()
	item, err := service.CacheDB.GetContract(address)
	if err != nil || item == nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[GetStorageContext] Get StorageContext nil")
	}
	if address != service.ContextRef.CurrentContext().ContractAddress {
		return errors.NewErr("[GetStorageContext] CodeHash not equal!")
	}
	return engine.EvalStack.PushAsInteropValue(NewStorageContext(address))
}

// ContractGetCode put contract to vm stack
func ContractGetCode(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if d, ok := i.Data.(*payload.DeployCode); ok {
		return engine.EvalStack.PushBytes(d.GetRawCode())
	}
	return fmt.Errorf("[ContractGetCode] Type error ")
}

func isContractParamValid(engine *vm.Executor) (*payload.DeployCode, error) {
	if engine.EvalStack.Count() < 7 {
		return nil, errors.NewErr("[Contract] Too few input parameters")
	}
	code, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	vmType, err := engine.EvalStack.PopAsInt64()
	if err != nil {
		return nil, err
	}
	name, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	version, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	author, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	email, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	desc, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return nil, err
	}

	contract, err := payload.CreateDeployCode(code, uint32(vmType), name, version, author, email, desc)
	if err != nil {
		return nil, err
	}

	if contract.VmType() != payload.NEOVM_TYPE {
		return nil, fmt.Errorf("[Contract] expect NEOVM_TYPE. get WASMVM_TYPE")
	}

	return contract, nil
}

func isContractExist(service *NeoVmService, contractAddress common.Address) error {
	item, err := service.CacheDB.GetContract(contractAddress)

	if err != nil || item != nil {
		return fmt.Errorf("[Contract] Get contract %x error or contract exist", contractAddress)
	}
	return nil
}
