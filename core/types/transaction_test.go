
package types

import (
	"math"
	"testing"

	" github.com/Daironode/aingle/core/payload"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_SigHashForChain(t *testing.T) {
	mutable := &MutableTransaction{
		TxType:  InvokeNeo,
		Payload: &payload.InvokeCode{},
	}

	tx, err := mutable.IntoImmutable()
	assert.Nil(t, err)

	assert.Equal(t, tx.Hash(), tx.SigHashForChain(0))
	assert.NotEqual(t, tx.Hash(), tx.SigHashForChain(1))
	assert.NotEqual(t, tx.Hash(), tx.SigHashForChain(math.MaxUint32))
}
