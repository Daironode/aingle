 package states

import (
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestValidatorState_Deserialize_Serialize(t *testing.T) {
	_, pubKey, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	vs := ValidatorState{
		StateBase: StateBase{(byte)(1)},
		PublicKey: pubKey,
	}

	sink := common.NewZeroCopySink(nil)
	vs.Serialization(sink)
	bs := sink.Bytes()

	var vs2 ValidatorState
	source := common.NewZeroCopySource(sink.Bytes())
	vs2.Deserialization(source)
	assert.Equal(t, vs, vs2)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err := vs2.Deserialization(source)
	assert.NotNil(t, err)
}
