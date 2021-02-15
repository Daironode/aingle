
package stateless

import (
	"reflect"

	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/validation"
	vatypes " github.com/Daironode/aingle/validator/types"
)

// Validator wraps validator actor's pid
type Validator interface {
	// Register send a register message to poolId
	Register(poolId *actor.PID)
	// UnRegister send an unregister message to poolId
	UnRegister(poolId *actor.PID)
	// VerifyType returns the type of validator
	VerifyType() vatypes.VerifyType
}

type validator struct {
	pid *actor.PID
	id  string
}

// NewValidator spawns a validator actor and return its pid wraped in Validator
func NewValidator(id string) (Validator, error) {
	validator := &validator{id: id}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	pid, err := actor.SpawnNamed(props, id)
	validator.pid = pid
	return validator, err
}

func (self *validator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("stateless-validator: started and be ready to receive txn")
	case *actor.Stopping:
		log.Info("stateless-validator: stopping")
	case *actor.Restarting:
		log.Info("stateless-validator: restarting")
	case *actor.Stopped:
		log.Info("stateless-validator: stopped")
	case *vatypes.CheckTx:
		log.Debugf("stateless-validator receive tx %x", msg.Tx.Hash())
		sender := context.Sender()
		errCode := validation.VerifyTransaction(msg.Tx)

		response := &vatypes.CheckResponse{
			WorkerId: msg.WorkerId,
			ErrCode:  errCode,
			Hash:     msg.Tx.Hash(),
			Type:     self.VerifyType(),
			Height:   0,
		}

		sender.Tell(response)
	case *vatypes.UnRegisterAck:
		context.Self().Stop()
	default:
		log.Info("stateless-validator: unknown msg ", msg, "type", reflect.TypeOf(msg))
	}

}

func (self *validator) VerifyType() vatypes.VerifyType {
	return vatypes.Stateless
}

// Register send RegisterValidator message to txpool
func (self *validator) Register(poolId *actor.PID) {
	poolId.Tell(&vatypes.RegisterValidator{
		Sender: self.pid,
		Type:   self.VerifyType(),
		Id:     self.id,
	})
}

// UnRegister send UnRegisterValidator message to txpool
func (self *validator) UnRegister(poolId *actor.PID) {
	poolId.Tell(&vatypes.UnRegisterValidator{
		Id: self.id,
	})

}
