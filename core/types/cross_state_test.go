 package types

import (
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestCrossState(t *testing.T) {
	sigData := make([][]byte, 3)
	sigData[0] = []byte{1, 2, 3, 4, 5}
	sigData[1] = []byte{2, 3, 4, 5, 6}
	sigData[2] = []byte{3, 4, 5, 6, 7}

	msg := &CrossChainMsg{
		Version:    CURR_CROSS_STATES_VERSION,
		Height:     1,
		StatesRoot: common.UINT256_EMPTY,
		SigData:    sigData,
	}
	sink := common.NewZeroCopySink(nil)
	msg.Serialization(sink)

	source := common.NewZeroCopySource(sink.Bytes())

	var msg1 CrossChainMsg
	err := msg1.Deserialization(source)

	assert.NoError(t, err)
	assert.Equal(t, *msg, msg1)
}
