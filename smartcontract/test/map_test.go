
package test

import (
	"fmt"
	"testing"

	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	byteCode := []byte{
		byte(neovm.NEWMAP),
		byte(neovm.DUP),   // dup map
		byte(neovm.PUSH0), // key (index)
		byte(neovm.PUSH0), // key (index)
		byte(neovm.SETITEM),

		byte(neovm.DUP),   // dup map
		byte(neovm.PUSH0), // key (index)
		byte(neovm.PUSH1), // value (newItem)
		byte(neovm.SETITEM),
	}

	// pick a value out
	byteCode = append(byteCode,
		[]byte{ // extract element
			byte(neovm.DUP),   // dup map (items)
			byte(neovm.PUSH0), // key (index)

			byte(neovm.PICKITEM),
			byte(neovm.JMPIF), // dup map (items)
			0x04, 0x00,        // skip a drop?
			byte(neovm.DROP),
		}...)

	// count faults vs successful executions
	N := 1024
	faults := 0

	//dbFile := "/tmp/test"
	//os.RemoveAll(dbFile)
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	panic(err)
	//}

	for n := 0; n < N; n++ {
		// Setup Execution Environment
		//store := statestore.NewMemDatabase()
		//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
		config := &smartcontract.Config{
			Time:   10,
			Height: 10,
			Tx:     &types.Transaction{},
		}
		sc := smartcontract.SmartContract{
			Config:  config,
			Gas:     100,
			CacheDB: nil,
		}
		engine, err := sc.NewExecuteEngine(byteCode, types.InvokeNeo)

		_, err = engine.Invoke()
		if err != nil {
			fmt.Println("err:", err)
			faults += 1
		}
	}
	assert.Equal(t, faults, 0)

}
