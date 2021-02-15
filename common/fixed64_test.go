 package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixed64_Serialize(t *testing.T) {
	val := Fixed64(10)
	buf := NewZeroCopySink(nil)
	val.Serialization(buf)
	val2 := Fixed64(0)
	val2.Deserialization(NewZeroCopySource(buf.Bytes()))

	assert.Equal(t, val, val2)
}

func TestFixed64_Deserialize(t *testing.T) {
	val := Fixed64(0)
	err := val.Deserialization(NewZeroCopySource([]byte{1, 2, 3}))

	assert.NotNil(t, err)

}
