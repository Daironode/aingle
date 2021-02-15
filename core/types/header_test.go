 package types

import (
	"fmt"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Serialize(t *testing.T) {
	header := Header{}
	header.Height = 321
	header.Bookkeepers = make([]keypair.PublicKey, 0)
	header.SigData = make([][]byte, 0)
	sink := common.NewZeroCopySink(nil)
	header.Serialization(sink)
	bs := sink.Bytes()

	var h2 Header
	source := common.NewZeroCopySource(bs)
	err := h2.Deserialization(source)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprint(header), fmt.Sprint(h2))

	var h3 RawHeader
	source = common.NewZeroCopySource(bs)
	err = h3.Deserialization(source)
	assert.Nil(t, err)
	assert.Equal(t, header.Height, h3.Height)
	assert.Equal(t, bs, h3.Payload)

}
