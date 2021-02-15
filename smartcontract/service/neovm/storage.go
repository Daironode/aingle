
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// StoragePut put smart contract storage item to cache
func StoragePut(service *NeoVmService, engine *vm.Executor) error {
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StoragePut] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] check context error!")
	}

	key, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	if len(key) > 1024 {
		return errors.NewErr("[StoragePut] Storage key to long")
	}

	value, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}

	service.CacheDB.Put(genStorageKey(context.Address, key), states.GenRawStorageItem(value))
	return nil
}

// StorageDelete delete smart contract storage item from cache
func StorageDelete(service *NeoVmService, engine *vm.Executor) error {
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StorageDelete] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] check context error!")
	}
	ba, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	service.CacheDB.Delete(genStorageKey(context.Address, ba))

	return nil
}

// StorageGet push smart contract storage item from cache to vm stack
func StorageGet(service *NeoVmService, engine *vm.Executor) error {

	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageGet] get pop context error!")
	}
	ba, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}

	raw, err := service.CacheDB.Get(genStorageKey(context.Address, ba))
	if err != nil {
		return err
	}

	if len(raw) == 0 {
		return engine.EvalStack.PushBytes([]byte{})
	}
	value, err := states.GetValueFromRawStorageItem(raw)
	if err != nil {
		return err
	}
	return engine.EvalStack.PushBytes(value)
}

// StorageGetContext push smart contract storage context to vm stack
func StorageGetContext(service *NeoVmService, engine *vm.Executor) error {
	return engine.EvalStack.PushAsInteropValue(NewStorageContext(service.ContextRef.CurrentContext().ContractAddress))
}

func StorageGetReadOnlyContext(service *NeoVmService, engine *vm.Executor) error {
	context := NewStorageContext(service.ContextRef.CurrentContext().ContractAddress)
	context.IsReadOnly = true
	return engine.EvalStack.PushAsInteropValue(context)
}

func checkStorageContext(service *NeoVmService, context *StorageContext) error {
	item, err := service.CacheDB.GetContract(context.Address)
	if err != nil || item == nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[CheckStorageContext] get context fail!")
	}
	return nil
}

func getContext(engine *vm.Executor) (*StorageContext, error) {
	opInterface, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return nil, err
	}
	if opInterface.Data == nil {
		return nil, errors.NewErr("[Context] Get storageContext nil")
	}
	context, ok := opInterface.Data.(*StorageContext)
	if !ok {
		return nil, errors.NewErr("[Context] Get storageContext invalid")
	}
	return context, nil
}

func genStorageKey(address common.Address, key []byte) []byte {
	res := make([]byte, 0, len(address[:])+len(key))
	res = append(res, address[:]...)
	res = append(res, key...)
	return res
}