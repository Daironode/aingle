 
package ont

import (
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestState_Serialize(t *testing.T) {
	state := State{
		From:  common.AddressFromVmCode([]byte{1, 2, 3}),
		To:    common.AddressFromVmCode([]byte{4, 5, 6}),
		Value: 1,
	}
	sink := common.NewZeroCopySink(nil)
	state.Serialization(sink)

	state2 := State{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := state2.Deserialization(source); err != nil {
		t.Fatal("state deserialize fail!")
	}

	assert.Equal(t, state, state2)
}
