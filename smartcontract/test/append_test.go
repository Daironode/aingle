 
package test

import (
	"testing"

	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/vm/neovm"
	" github.com/Daironode/aingle/vm/neovm/errors"
	"github.com/stretchr/testify/assert"
)

func TestAppendOverFlow(t *testing.T) {
	// define 1024 len array
	byteCode := []byte{
		byte(0x04), //neovm.PUSHBYTES4
		byte(0x00),
		byte(0x04),
		byte(0x00),
		byte(0x00),
		byte(neovm.NEWARRAY),
		byte(neovm.PUSH2),
		byte(neovm.APPEND),
	}

	config := &smartcontract.Config{
		Time:   10,
		Height: 10,
	}
	sc := smartcontract.SmartContract{
		Config:  config,
		Gas:     200,
		CacheDB: nil,
	}
	engine, _ := sc.NewExecuteEngine(byteCode, types.InvokeNeo)
	_, err := engine.Invoke()
	assert.EqualError(t, err, "[NeoVmService] vm execution error!: "+errors.ERR_OVER_MAX_ARRAY_SIZE.Error())
}
