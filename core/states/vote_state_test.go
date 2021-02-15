 package states

import (
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestVoteState_Deserialize_Serialize(t *testing.T) {
	_, pubKey1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	_, pubKey2, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	vs := VoteState{
		StateBase:  StateBase{(byte)(1)},
		PublicKeys: []keypair.PublicKey{pubKey1, pubKey2},
		Count:      10,
	}

	sink := common.NewZeroCopySink(nil)
	vs.Serialization(sink)
	bs := sink.Bytes()

	var vs2 VoteState
	source := common.NewZeroCopySource(bs)
	vs2.Deserialization(source)
	assert.Equal(t, vs, vs2)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err := vs2.Deserialization(source)
	assert.NotNil(t, err)
}
