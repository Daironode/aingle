
package actor

import (
	"errors"
	"time"

	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/core/types"
	ontErrors " github.com/Daironode/aingle/errors"
	txpool " github.com/Daironode/aingle/txnpool/common"
)

type TxPoolActor struct {
	Pool *actor.PID
}

func (self *TxPoolActor) GetTxnPool(byCount bool, height uint32) []*txpool.TXEntry {
	poolmsg := &txpool.GetTxnPoolReq{ByCount: byCount, Height: height}
	future := self.Pool.RequestFuture(poolmsg, time.Second*10)
	entry, err := future.Result()
	if err != nil {
		return nil
	}

	txs := entry.(*txpool.GetTxnPoolRsp).TxnPool
	return txs
}

func (self *TxPoolActor) VerifyBlock(txs []*types.Transaction, height uint32) error {
	poolmsg := &txpool.VerifyBlockReq{Txs: txs, Height: height}
	future := self.Pool.RequestFuture(poolmsg, time.Second*10)
	entry, err := future.Result()
	if err != nil {
		return err
	}

	txentry := entry.(*txpool.VerifyBlockRsp).TxnPool
	for _, entry := range txentry {
		if entry.ErrCode != ontErrors.ErrNoError {
			return errors.New(entry.ErrCode.Error())
		}
	}

	return nil
}
