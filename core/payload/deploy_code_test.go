 package payload

import (
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestDeployCode_Serialize(t *testing.T) {
	ty, _ := VmTypeFromByte(1)
	deploy, err := NewDeployCode([]byte{1, 2, 3}, ty, "", "", "", "", "")
	assert.Nil(t, err)
	sink := common.NewZeroCopySink(nil)
	deploy.Serialization(sink)
	bs := sink.Bytes()
	var deploy2 DeployCode

	source := common.NewZeroCopySource(bs)
	deploy2.Deserialization(source)
	assert.Equal(t, &deploy2, deploy)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err = deploy2.Deserialization(source)
	assert.NotNil(t, err)
}
