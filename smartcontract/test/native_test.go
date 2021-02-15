
package test

import (
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	"github.com/stretchr/testify/assert"
)

func TestBuildParamToNative(t *testing.T) {
	code := `00c57676c84c0500000000004c1400000000000000000000000000000000000000060068164f6e746f6c6f67792e4e61746976652e496e766f6b65`

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

	assert.Error(t, err, "invoke smart contract err: [NeoVmService] service system call error!: [SystemCall] service execute error!: invoke native circular reference!")
}
