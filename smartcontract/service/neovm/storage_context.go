 
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// StorageContext store smart contract address
type StorageContext struct {
	Address    common.Address
	IsReadOnly bool
}

// NewStorageContext return a new smart contract storage context
func NewStorageContext(address common.Address) *StorageContext {
	var storageContext StorageContext
	storageContext.Address = address
	storageContext.IsReadOnly = false
	return &storageContext
}

// ToArray return address byte array
func (this *StorageContext) ToArray() []byte {
	return this.Address[:]
}

func StorageContextAsReadOnly(service *NeoVmService, engine *vm.Executor) error {
	data, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	context, ok := data.Data.(*StorageContext)
	if !ok {
		return fmt.Errorf("%s", "pop storage context type invalid")
	}
	if !context.IsReadOnly {
		context = NewStorageContext(context.Address)
		context.IsReadOnly = true
	}
	return engine.EvalStack.PushAsInteropValue(context)
}
