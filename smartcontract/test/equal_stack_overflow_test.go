
package test

import (
	"os"
	"testing"

	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
	. " github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestEqualStackOverflow(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("./Log")
	}()

	code := []byte{
		byte(neovm.PUSH1),    // {1}
		byte(neovm.NEWARRAY), // {[]}
		byte(neovm.DUP),      // {[],[]}
		byte(neovm.DUP),      // {[],[],[]}
		byte(neovm.PUSH0),    // {[],[],[],0}
		byte(neovm.ROT),      // {[],[],0,[]}
		byte(neovm.SETITEM),  // {[[]]}
		byte(neovm.DUP),      // {[[]],[[]]}
		byte(neovm.EQUAL),
	}

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	sc := SmartContract{
		Config:  config,
		Gas:     10000,
		CacheDB: nil,
	}
	engine, _ := sc.NewExecuteEngine(code, types.InvokeNeo)
	_, err := engine.Invoke()

	assert.Nil(t, err)
}
