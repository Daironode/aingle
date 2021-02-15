 package states

import (
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestBookkeeper_Deserialize_Serialize(t *testing.T) {
	_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey3, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey4, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	bk := BookkeeperState{
		StateBase:      StateBase{(byte)(1)},
		CurrBookkeeper: []keypair.PublicKey{pubKey1, pubKey2},
		NextBookkeeper: []keypair.PublicKey{pubKey3, pubKey4},
	}

	sink := common.NewZeroCopySink(nil)
	bk.Serialization(sink)
	bs := sink.Bytes()

	var bk2 BookkeeperState
	source := common.NewZeroCopySource(bs)
	bk2.Deserialization(source)
	assert.Equal(t, bk, bk2)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err := bk2.Deserialization(source)
	assert.NotNil(t, err)
}
