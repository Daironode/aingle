 package crossvm_codec

import (
	"encoding/hex"
	"math/big"
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestDe1(t *testing.T) {
	h, _ := hex.DecodeString("657674001001000000010500000068656c6c6f")

	_, err := parseNotify(h)
	assert.Nil(t, err)
}

func EncodeNotify(t *testing.T, value interface{}) []byte {
	val, err := EncodeValue(value)
	assert.Nil(t, err)

	return append([]byte("evt\x00"), val...)
}

func TestDeserializeNotify(t *testing.T) {
	addr := common.AddressFromVmCode([]byte("123"))
	value := []interface{}{"helloworld", []byte("1234"), 123, -1, -128, -260, true, big.NewInt(100), addr, common.UINT256_EMPTY}
	expected := []interface{}{"helloworld", hex.EncodeToString([]byte("1234")), "123", "-1", "-128", "-260", true, "100", addr.ToBase58(), common.UINT256_EMPTY.ToHexString()}
	for i, val := range value {
		assert.Equal(t, DeserializeNotify(EncodeNotify(t, val)), interface{}(expected[i]))
	}

	assert.Equal(t, DeserializeNotify(EncodeNotify(t, value)), interface{}(expected))
}
