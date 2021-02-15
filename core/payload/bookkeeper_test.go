 
package payload

import (
	"encoding/hex"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestBookkeeper_Serialization(t *testing.T) {
	pubkey, err := hex.DecodeString("039af138392513408f9d1509c651c60066c05b2305de17e44f68088510563e2279")
	assert.Nil(t, err)
	pub, err := keypair.DeserializePublicKey(pubkey)
	assert.Nil(t, err)
	bookkeeper := &Bookkeeper{
		PubKey: pub,
		Action: BookkeeperAction(1),
		Cert:   pubkey,
		Issuer: pub,
	}
	sink := common.NewZeroCopySink(nil)
	bookkeeper.Serialization(sink)
	bookkeeper2 := &Bookkeeper{}
	source := common.NewZeroCopySource(sink.Bytes())
	err = bookkeeper2.Deserialization(source)
	assert.Nil(t, err)

	assert.Equal(t, bookkeeper, bookkeeper2)
}
