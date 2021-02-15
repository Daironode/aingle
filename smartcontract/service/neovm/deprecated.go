 
package neovm

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/errors"
	vm " github.com/Daironode/aingle/vm/neovm"
)

// BlockChainGetHeader put blockchain's header to vm stack
func BlockChainGetHeader(service *NeoVmService, engine *vm.Executor) error {
	var (
		header *types.Header
		err    error
	)
	data, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	l := len(data)
	if l <= 5 {
		b := common.BigIntFromNeoBytes(data)
		height := uint32(b.Int64())
		hash := service.Store.GetBlockHash(height)
		header, err = service.Store.GetHeaderByHash(hash)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetHeader] GetHeader error!.")
		}
	} else if l == 32 {
		hash, _ := common.Uint256ParseFromBytes(data)
		header, err = service.Store.GetHeaderByHash(hash)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetHeader] GetHeader error!.")
		}
	} else {
		return errors.NewErr("[BlockChainGetHeader] data invalid.")
	}

	err = engine.EvalStack.PushAsInteropValue(header)
	if err != nil {
		return errors.NewErr("[BlockChainGetHeader] PushAsInteropValue error.")
	}
	return nil
}

// BlockChainGetBlock put blockchain's block to vm stack
func BlockChainGetBlock(service *NeoVmService, engine *vm.Executor) error {
	data, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	var block *types.Block
	l := len(data)
	if l <= 5 {
		b := common.BigIntFromNeoBytes(data)
		height := uint32(b.Int64())
		var err error
		block, err = service.Store.GetBlockByHeight(height)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] GetBlock error!.")
		}
	} else if l == 32 {
		hash, err := common.Uint256ParseFromBytes(data)
		if err != nil {
			return err
		}
		block, err = service.Store.GetBlockByHash(hash)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] GetBlock error!.")
		}
	} else {
		return errors.NewErr("[BlockChainGetBlock] data invalid.")
	}
	err = engine.EvalStack.PushAsInteropValue(block)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] PushAsInteropValue error!.")
	}
	return nil
}

// BlockChainGetTransaction put blockchain's transaction to vm stack
func BlockChainGetTransaction(service *NeoVmService, engine *vm.Executor) error {
	d, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil {
		return err
	}
	t, _, err := service.Store.GetTransaction(hash)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetTransaction] GetTransaction error!")
	}
	err = engine.EvalStack.PushAsInteropValue(t)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetTransaction] PushAsInteropValue error!")
	}
	return nil
}

// BlockChainGetTransactionHeight put transaction in block height to vm stack
func BlockChainGetTransactionHeight(service *NeoVmService, engine *vm.Executor) error {
	d, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil {
		return err
	}
	_, h, err := service.Store.GetTransaction(hash)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetTransactionHeight] GetTransaction error!")
	}
	err = engine.EvalStack.PushUint32(h)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetTransactionHeight] PushInt64 error!")
	}
	return nil
}
