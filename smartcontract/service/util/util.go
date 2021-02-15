 package util

import (
	"bytes"
	"errors"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/utils"
	" github.com/Daironode/aingle/smartcontract/context"
	neovms " github.com/Daironode/aingle/smartcontract/service/neovm"
	" github.com/Daironode/aingle/vm/crossvm_codec"
	" github.com/Daironode/aingle/vm/neovm"
)

func BuildNeoVMParamEvalStack(params []interface{}) (*neovm.ValueStack, error) {
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := utils.BuildNeoVMParam(builder, params)
	if err != nil {
		return nil, err
	}

	exec := neovm.NewExecutor(builder.ToArray(), neovm.VmFeatureFlag{true, true})
	err = exec.Execute()
	if err != nil {
		return nil, err
	}
	return exec.EvalStack, nil
}

//create paramters for neovm contract
func GenerateNeoVMParamEvalStack(input []byte) (*neovm.ValueStack, error) {
	params, err := crossvm_codec.DeserializeCallParam(input)
	if err != nil {
		return nil, err
	}

	list, ok := params.([]interface{})
	if !ok {
		return nil, errors.New("invoke neovm param is not list type")
	}

	stack, err := BuildNeoVMParamEvalStack(list)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func SetNeoServiceParamAndEngine(addr common.Address, engine context.Engine, stack *neovm.ValueStack) error {
	service, ok := engine.(*neovms.NeoVmService)
	if !ok {
		return errors.New("engine should be NeoVmService")
	}

	code, err := service.GetNeoContract(addr)
	if err != nil {
		return err
	}

	feature := service.Engine.Features
	service.Engine = neovm.NewExecutor(code, feature)
	service.Code = code

	service.Engine.EvalStack = stack

	return nil
}
