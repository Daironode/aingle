 package states

import (
	"testing"

	" github.com/Daironode/aingle/common"
)

func TestStateBase_Serialize_Deserialize(t *testing.T) {

	st := &StateBase{byte(1)}

	bf := common.NewZeroCopySink(nil)
	st.Serialization(bf)

	var st2 = new(StateBase)
	source := common.NewZeroCopySource(bf.Bytes())
	if err := st2.Deserialization(source); err != nil {
		t.Fatalf("StateBase deserialize error: %v", err)
	}
}
