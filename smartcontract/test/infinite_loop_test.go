
package test

import (
	"os"
	"testing"

	" github.com/Daironode/aingle/core/types"
	. " github.com/Daironode/aingle/smartcontract"
	"github.com/stretchr/testify/assert"
)

func TestInfiniteLoopCrash(t *testing.T) {
	evilBytecode := []byte(" e\xff\u007f\xffhm\xb7%\xa7AAAAAAAAAAAAAAAC\xef\xed\x04INVERT\x95ve")
	dbFile := "test"
	defer func() {
		os.RemoveAll(dbFile)
	}()
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//store := statestore.NewMemDatabase()
	//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
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
	engine, err := sc.NewExecuteEngine(evilBytecode, types.InvokeNeo)
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Invoke()

	assert.NotNil(t, err)
}
