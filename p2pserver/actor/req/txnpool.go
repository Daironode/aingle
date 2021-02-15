 
package req

import (
	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
	tc " github.com/Daironode/aingle/txnpool/common"
)

var txnPoolPid *actor.PID

func SetTxnPoolPid(txnPid *actor.PID) {
	txnPoolPid = txnPid
}

//add txn to txnpool
func AddTransaction(transaction *types.Transaction) {
	if txnPoolPid == nil {
		log.Error("[p2p]net_server AddTransaction(): txnpool pid is nil")
		return
	}
	txReq := &tc.TxReq{
		Tx:         transaction,
		Sender:     tc.NetSender,
		TxResultCh: nil,
	}
	txnPoolPid.Tell(txReq)
}
