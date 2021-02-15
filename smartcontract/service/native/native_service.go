
package native

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/store"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/merkle"
	" github.com/Daironode/aingle/smartcontract/context"
	" github.com/Daironode/aingle/smartcontract/event"
	" github.com/Daironode/aingle/smartcontract/states"
	" github.com/Daironode/aingle/smartcontract/storage"
)

type (
	Handler         func(native *NativeService) ([]byte, error)
	RegisterService func(native *NativeService)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// Native service struct
// Invoke a native smart contract, new a native service
type NativeService struct {
	Store         store.LedgerStore
	CacheDB       *storage.CacheDB
	ServiceMap    map[string]Handler
	Notifications []*event.NotifyEventInfo
	InvokeParam   states.ContractInvokeParam
	Input         []byte
	Tx            *types.Transaction
	Height        uint32
	Time          uint32
	BlockHash     common.Uint256
	ContextRef    context.ContextRef
	PreExec       bool
	CrossHashes   []common.Uint256
}

func (this *NativeService) Register(methodName string, handler Handler) {
	this.ServiceMap[methodName] = handler
}

func (this *NativeService) Invoke() ([]byte, error) {
	contract := this.InvokeParam
	services, ok := Contracts[contract.Address]
	if !ok {
		return BYTE_FALSE, fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}
	services(this)
	service, ok := this.ServiceMap[contract.Method]
	if !ok {
		return BYTE_FALSE, fmt.Errorf("Native contract %x doesn't support this function %s.",
			contract.Address, contract.Method)
	}
	args := this.Input
	this.Input = contract.Args
	this.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	notifications := this.Notifications
	this.Notifications = []*event.NotifyEventInfo{}
	hashes := this.CrossHashes
	this.CrossHashes = []common.Uint256{}
	result, err := service(this)
	if err != nil {
		return result, errors.NewDetailErr(err, errors.ErrNoCode, "[Invoke] Native serivce function execute error!")
	}
	this.ContextRef.PopContext()
	this.ContextRef.PushNotifications(this.Notifications)
	this.ContextRef.PutCrossStateHashes(this.CrossHashes)
	this.Notifications = notifications
	this.Input = args
	this.CrossHashes = hashes
	return result, nil
}

func (this *NativeService) NativeCall(address common.Address, method string, args []byte) ([]byte, error) {
	c := states.ContractInvokeParam{
		Address: address,
		Method:  method,
		Args:    args,
	}
	this.InvokeParam = c
	return this.Invoke()
}

func (this *NativeService) PushCrossState(data []byte) {
	this.CrossHashes = append(this.CrossHashes, merkle.HashLeaf(data))
}
