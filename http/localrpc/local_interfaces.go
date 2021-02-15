
package localrpc

import (
	"time"

	" github.com/Daironode/aingle/common/log"
	bactor " github.com/Daironode/aingle/http/base/actor"
	" github.com/Daironode/aingle/http/base/common"
	berr " github.com/Daironode/aingle/http/base/error"
	" github.com/Daironode/aingle/http/base/rpc"
)

func GetNeighbor(params []interface{}) map[string]interface{} {
	addr := bactor.GetNeighborAddrs()
	return rpc.ResponseSuccess(addr)
}

func GetNodeState(params []interface{}) map[string]interface{} {
	t := time.Now().UnixNano()
	port := bactor.GetNodePort()
	id := bactor.GetID()
	ver := bactor.GetVersion()
	tpe := bactor.GetNodeType()
	relay := bactor.GetRelayState()
	height := bactor.GetCurrentBlockHeight()
	txnCnt, err := bactor.GetTxnCount()
	if err != nil {
		return rpc.ResponsePack(berr.INTERNAL_ERROR, false)
	}
	n := common.NodeInfo{
		NodeTime:    t,
		NodePort:    port,
		ID:          id,
		NodeVersion: ver,
		NodeType:    tpe,
		Relay:       relay,
		Height:      height,
		TxnCnt:      txnCnt,
	}
	return rpc.ResponseSuccess(n)
}

func StartConsensus(params []interface{}) map[string]interface{} {
	if err := bactor.ConsensusSrvStart(); err != nil {
		return rpc.ResponsePack(berr.INTERNAL_ERROR, false)
	}
	return rpc.ResponsePack(berr.SUCCESS, true)
}

func StopConsensus(params []interface{}) map[string]interface{} {
	if err := bactor.ConsensusSrvHalt(); err != nil {
		return rpc.ResponsePack(berr.INTERNAL_ERROR, false)
	}
	return rpc.ResponsePack(berr.SUCCESS, true)
}

func SetDebugInfo(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return rpc.ResponsePack(berr.INVALID_PARAMS, "")
	}
	switch params[0].(type) {
	case float64:
		level := params[0].(float64)
		if err := log.Log().SetDebugLevel(int(level)); err != nil {
			return rpc.ResponsePack(berr.INVALID_PARAMS, "")
		}
	default:
		return rpc.ResponsePack(berr.INVALID_PARAMS, "")
	}
	return rpc.ResponsePack(berr.SUCCESS, true)
}
