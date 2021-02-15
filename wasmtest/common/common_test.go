 package common

import (
	"encoding/json"
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestTestEnv(t *testing.T) {
	a := TestEnv{Witness: []common.Address{common.ADDRESS_EMPTY}}

	encoded, _ := json.Marshal(&a)
	assert.Equal(t, string(encoded), `{"witness":["AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM"]}`)

	var b TestEnv
	err := json.Unmarshal(encoded, &b)
	assert.Nil(t, err)
	assert.Equal(t, a, b)
}

func TestTestCase(t *testing.T) {
	a := TestEnv{Witness: []common.Address{common.ADDRESS_EMPTY}}
	ts := TestCase{Env: a, Method: "func1", Param: "int:100, bool:true", Expect: "int:10"}

	encoded, _ := json.Marshal(ts)

	assert.Equal(t, string(encoded), `{"env":{"witness":["AFmseVrdL9f9oyCzZefL9tG6UbvhPbdYzM"]},"needcontext":false,"method":"func1","param":"int:100, bool:true","expected":"int:10","notify":""}`)
}
