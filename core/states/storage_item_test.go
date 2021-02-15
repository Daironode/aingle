 package states

import (
	"testing"

	" github.com/Daironode/aingle/common"
)

func TestStorageItem_Serialize_Deserialize(t *testing.T) {

	item := &StorageItem{
		StateBase: StateBase{StateVersion: 1},
		Value:     []byte{1},
	}

	bf := common.NewZeroCopySink(nil)
	item.Serialization(bf)

	var storage = new(StorageItem)
	source := common.NewZeroCopySource(bf.Bytes())
	if err := storage.Deserialization(source); err != nil {
		t.Fatalf("StorageItem deserialize error: %v", err)
	}
}
