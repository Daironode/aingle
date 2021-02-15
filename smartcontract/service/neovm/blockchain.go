 
package neovm

import (
	"math/big"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// BlockChainGetHeight put blockchain's height to vm stack
func BlockChainGetHeight(service *NeoVmService, engine *vm.Executor) error {
	return engine.EvalStack.PushUint32(service.Height - 1)
}

func BlockChainGetHeightNew(service *NeoVmService, engine *vm.Executor) error {
	return engine.EvalStack.PushUint32(service.Height)
}

func BlockChainGetHeaderNew(service *NeoVmService, engine *vm.Executor) error {
	data, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	b := common.BigIntFromNeoBytes(data)
	if b.Cmp(big.NewInt(int64(service.Height))) != 0 {
		return errors.NewErr("can only get current block header")
	}

	header := &HeaderValue{Height: service.Height, Timestamp: service.Time, Hash: service.BlockHash}
	return engine.EvalStack.PushAsInteropValue(header)
}

// BlockChainGetContract put blockchain's contract to vm stack
func BlockChainGetContract(service *NeoVmService, engine *vm.Executor) error {
	b, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	address, err := common.AddressParseFromBytes(b)
	if err != nil {
		return err
	}
	item, err := service.Store.GetContractState(address)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetContract] GetContract error!")
	}
	err = engine.EvalStack.PushAsInteropValue(item)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetContract] PushAsInteropValue error!")
	}
	return nil
}
