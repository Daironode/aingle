 package stateless

import (
	"testing"
	"time"

	"github.com/Daironode/aingle-crypto/keypair"
	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/signature"
	ctypes " github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/core/utils"
	" github.com/Daironode/aingle/errors"
	types2 " github.com/Daironode/aingle/validator/types"
	"github.com/stretchr/testify/assert"
)

func signTransaction(signer *account.Account, tx *ctypes.MutableTransaction) error {
	hash := tx.Hash()
	sign, _ := signature.Sign(signer, hash[:])
	tx.Sigs = append(tx.Sigs, ctypes.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sign},
	})
	return nil
}

func TestStatelessValidator(t *testing.T) {
	log.InitLog(log.InfoLog, log.Stdout)
	acc := account.NewAccount("")

	code := []byte{1, 2, 3}

	mutable, err := utils.NewDeployTransaction(code, "test", "1", "author", "author@123.com", "test desp", payload.NEOVM_TYPE)
	assert.Nil(t, err)
	mutable.Payer = acc.Address

	signTransaction(acc, mutable)

	tx, err := mutable.IntoImmutable()
	assert.Nil(t, err)

	validator := &validator{id: "test"}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	pid, err := actor.SpawnNamed(props, validator.id)
	assert.Nil(t, err)

	msg := &types2.CheckTx{WorkerId: 1, Tx: tx}
	fut := pid.RequestFuture(msg, time.Second)

	res, err := fut.Result()
	assert.Nil(t, err)

	result := res.(*types2.CheckResponse)
	assert.Equal(t, result.ErrCode, errors.ErrNoError)
	assert.Equal(t, mutable.Hash(), result.Hash)
}
