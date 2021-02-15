 
package test

import (
	"os"
	"testing"

	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/store/ledgerstore"
	"github.com/stretchr/testify/assert"
)

func TestPreExecuteContractWasmDeploy(t *testing.T) {
	acct := account.NewAccount("")
	testLedgerStore, err := ledgerstore.NewLedgerStore("test/ledgerfortmp", 0)
	/** file: test_create.wat
		(module
		  (type (;0;) (func))
		  (type (;1;) (func (param i32 i32)))
		  (import "env" "aingle_return" (func (;0;) (type 1)))
		  (func (;1;) (type 0)
			i32.const 0
		    i64.const 2222
			i64.store
			i32.const 0
			i32.const 8
			call 0
			)
		  (memory (;0;) 1)
		  (export "invoke" (func 1)))
	 **/
	code := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00, 0x01, 0x09, 0x02, 0x60, 0x00, 0x00, 0x60, 0x02, 0x7f, 0x7f, 0x00, 0x02, 0x14, 0x01, 0x03, 0x65, 0x6e, 0x76, 0x0c, 0x6f, 0x6e, 0x74, 0x69, 0x6f, 0x5f, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x00, 0x01, 0x03, 0x02, 0x01, 0x00, 0x05, 0x03, 0x01, 0x00, 0x01, 0x07, 0x0a, 0x01, 0x06, 0x69, 0x6e, 0x76, 0x6f, 0x6b, 0x65, 0x00, 0x01, 0x0a, 0x12, 0x01, 0x10, 0x00, 0x41, 0x00, 0x42, 0xae, 0x11, 0x37, 0x03, 0x00, 0x41, 0x00, 0x41, 0x08, 0x10, 0x00, 0x0b}
	mutable, _ := utils.NewDeployCodeTransaction(0, 100000000, code, payload.NEOVM_TYPE, "name", "version",
		"author", "email", "desc")
	_ = utils.SignTransaction(acct, mutable)
	tx, err := mutable.IntoImmutable()
	_, err = testLedgerStore.PreExecuteContract(tx)
	assert.EqualError(t, err, "this code is wasm binary. can not deployed as neo contract")

	mutable, _ = utils.NewDeployCodeTransaction(0, 100000000, code, payload.WASMVM_TYPE, "name", "version",
		"author", "email", "desc")
	_ = utils.SignTransaction(acct, mutable)
	tx, err = mutable.IntoImmutable()
	_, err = testLedgerStore.PreExecuteContract(tx)
	assert.Nil(t, err)

	_ = os.RemoveAll("./test")
}
