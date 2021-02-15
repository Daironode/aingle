
package test

import (
	"testing"

	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/smartcontract"
	" github.com/Daironode/aingle/vm/neovm"
	" github.com/Daironode/aingle/vm/neovm/errors"
	"github.com/stretchr/testify/assert"
)

func TestHeight(t *testing.T) {
	byteCode0 := []byte{
		byte(neovm.NEWMAP),
		byte(neovm.PUSH0),
		byte(neovm.HASKEY),
	}

	byteCode1 := []byte{
		byte(neovm.NEWMAP),
		byte(neovm.KEYS),
	}

	byteCode2 := []byte{
		byte(neovm.NEWMAP),
		byte(neovm.VALUES),
	}

	bytecode := [...][]byte{byteCode0, byteCode1, byteCode2}

	disableHeight := config.GetOpcodeUpdateCheckHeight(config.DefConfig.P2PNode.NetworkId)
	heights := []uint32{10, disableHeight, disableHeight + 1}

	for _, height := range heights {
		config := &smartcontract.Config{Time: 10, Height: height}
		sc := smartcontract.SmartContract{Config: config, Gas: 100}
		expected := "[NeoVmService] vm execution error!: " + errors.ERR_NOT_SUPPORT_OPCODE.Error()
		if height > disableHeight {
			expected = ""
		}
		for i := 0; i < 3; i++ {
			engine, err := sc.NewExecuteEngine(bytecode[i], types.InvokeNeo)
			assert.Nil(t, err)

			_, err = engine.Invoke()
			if len(expected) > 0 {
				assert.EqualError(t, err, expected)
			} else {
				assert.Nil(t, err)
			}
		}
	}
}
