 
package test

import (
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestConvertNeoVmTypeHexString(t *testing.T) {
	code := `00c57676c8681553797374656d2e52756e74696d652e4e6f74696679`

	hex, err := common.HexToBytes(code)

	if err != nil {
		t.Fatal("hex to byte error:", err)
	}

	config := &smartcontract.Config{
		Time:   10,
		Height: 10,
		Tx:     nil,
	}
	sc := smartcontract.SmartContract{
		Config: config,
		Gas:    100000,
	}
	engine, err := sc.NewExecuteEngine(hex, types.InvokeNeo)

	_, err = engine.Invoke()

	assert.Error(t, err, "over max parameters convert length")
}

func BenchmarkExecuteAdd(b *testing.B) {
	code := []byte{byte(neovm.PUSH1)}

	N := 50000
	for i := 0; i < N; i++ {
		code = append(code, byte(neovm.PUSH1), byte(neovm.ADD))
	}
	code = append(code, byte(neovm.RET))

	config := &smartcontract.Config{
		Time:   10,
		Height: 10,
		Tx:     nil,
	}

	for i := 0; i < b.N; i++ {
		sc := smartcontract.SmartContract{
			Config: config,
			Gas:    1000000,
		}
		engine, err := sc.NewExecuteEngine(code, types.InvokeNeo)
		if err != nil {
			panic(err)
		}
		_, err = engine.Invoke()
		if err != nil {
			panic(err)
		}
	}

}
