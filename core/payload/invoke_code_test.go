 package payload

import (
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestInvokeCode_Serialize(t *testing.T) {
	code := InvokeCode{
		Code: []byte{1, 2, 3},
	}

	sink := common.NewZeroCopySink(nil)
	code.Serialization(sink)
	bs := sink.Bytes()
	var code2 InvokeCode
	source := common.NewZeroCopySource(bs)
	code2.Deserialization(source)
	assert.Equal(t, code, code2)

	source = common.NewZeroCopySource(bs[:len(bs)-2])
	err := code.Deserialization(source)

	assert.NotNil(t, err)
}
