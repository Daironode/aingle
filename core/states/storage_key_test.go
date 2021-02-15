 package states

import (
	"testing"

	"crypto/rand"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestStorageKey_Deserialize_Serialize(t *testing.T) {
	var addr common.Address
	rand.Read(addr[:])

	storage := StorageKey{
		ContractAddress: addr,
		Key:             []byte{1, 2, 3},
	}

	sink := common.NewZeroCopySink(nil)
	storage.Serialization(sink)
	bs := sink.Bytes()

	var storage2 StorageKey
	source := common.NewZeroCopySource(sink.Bytes())
	storage2.Deserialization(source)
	assert.Equal(t, storage, storage2)

	buf := common.NewZeroCopySource(bs[:len(bs)-1])
	err := storage2.Deserialization(buf)
	assert.NotNil(t, err)
}
