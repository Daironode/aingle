
package test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
	. " github.com/Daironode/aingle/smartcontract"
	neovm2 " github.com/Daironode/aingle/smartcontract/service/neovm"
	" github.com/Daironode/aingle/vm/neovm"
	"github.com/stretchr/testify/assert"
)

func TestRandomCodeCrash(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("code %x \n", code)
		}
	}()

	for i := 1; i < 10; i++ {
		fmt.Printf("test round:%d \n", i)
		code := make([]byte, i)
		for j := 0; j < 10; j++ {
			rand.Read(code)

			sc := SmartContract{
				Config:  config,
				Gas:     10000,
				CacheDB: nil,
			}
			engine, _ := sc.NewExecuteEngine(code, types.InvokeNeo)
			engine.Invoke()
		}
	}
}

func TestOpCodeDUP(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code = []byte{byte(neovm.DUP)}

	sc := SmartContract{
		Config:  config,
		Gas:     10000,
		CacheDB: nil,
	}
	engine, _ := sc.NewExecuteEngine(code, types.InvokeNeo)
	_, err := engine.Invoke()

	assert.NotNil(t, err)
}

func TestOpReadMemAttack(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	bf := new(bytes.Buffer)
	builder := neovm.NewParamsBuilder(bf)
	builder.Emit(neovm.SYSCALL)
	sink := common.NewZeroCopySink(builder.ToArray())
	builder.EmitPushByteArray([]byte(neovm2.NATIVE_INVOKE_NAME))
	l := 0x7fffffc7 - 1
	sink.WriteVarUint(uint64(l))
	b := make([]byte, 4)
	sink.WriteBytes(b)

	sc := SmartContract{
		Config:  config,
		Gas:     100000,
		CacheDB: nil,
	}
	engine, _ := sc.NewExecuteEngine(sink.Bytes(), types.InvokeNeo)
	_, err := engine.Invoke()

	assert.NotNil(t, err)

}
