
// Package websocket privides a function to start websocket server
package websocket

import (
	" github.com/Daironode/aingle/common"
	cfg " github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
	" github.com/Daironode/aingle/events/message"
	bactor " github.com/Daironode/aingle/http/base/actor"
	bcomn " github.com/Daironode/aingle/http/base/common"
	Err " github.com/Daironode/aingle/http/base/error"
	" github.com/Daironode/aingle/http/base/rest"
	" github.com/Daironode/aingle/http/websocket/websocket"
	" github.com/Daironode/aingle/smartcontract/event"
)

var ws *websocket.WsServer

func StartServer() {
	bactor.SubscribeEvent(message.TOPIC_SAVE_BLOCK_COMPLETE, sendBlock2WSclient)
	bactor.SubscribeEvent(message.TOPIC_SMART_CODE_EVENT, pushSmartCodeEvent)
	go func() {
		ws = websocket.InitWsServer()
		ws.Start()
	}()
}
func sendBlock2WSclient(v interface{}) {
	if cfg.DefConfig.Ws.HttpWsPort != 0 {
		go func() {
			pushBlock(v)
			pushBlockTransactions(v)
		}()
	}
}
func Stop() {
	if ws == nil {
		return
	}
	ws.Stop()
}
func ReStartServer() {
	if ws == nil {
		ws = websocket.InitWsServer()
		ws.Start()
		return
	}
	ws.Restart()
}

func pushSmartCodeEvent(v interface{}) {
	if ws == nil {
		return
	}
	rs, ok := v.(types.SmartCodeEvent)
	if !ok {
		log.Errorf("[PushSmartCodeEvent]", "SmartCodeEvent err")
		return
	}
	go func() {
		switch object := rs.Result.(type) {
		case *event.LogEventArgs:
			contractAddrs, evts := bcomn.GetLogEvent(object)
			pushEvent(contractAddrs, rs.TxHash.ToHexString(), rs.Error, rs.Action, evts)
		case *event.ExecuteNotify:
			contractAddrs, notify := bcomn.GetExecuteNotify(object)
			pushEvent(contractAddrs, rs.TxHash.ToHexString(), rs.Error, rs.Action, notify)
		default:
		}
	}()
}

func pushEvent(contractAddrs map[string]bool, txHash string, errcode int64, action string, result interface{}) {
	if ws != nil {
		resp := rest.ResponsePack(Err.SUCCESS)
		resp["Result"] = result
		resp["Error"] = errcode
		resp["Action"] = action
		resp["Desc"] = Err.ErrMap[resp["Error"].(int64)]
		ws.PushTxResult(contractAddrs, txHash, resp)
		ws.BroadcastToSubscribers(contractAddrs, websocket.WSTOPIC_EVENT, resp)
	}
}

func pushBlock(v interface{}) {
	if ws == nil {
		return
	}
	resp := rest.ResponsePack(Err.SUCCESS)
	if block, ok := v.(types.Block); ok {
		resp["Action"] = "sendrawblock"
		resp["Result"] = common.ToHexString(block.ToArray())
		ws.BroadcastToSubscribers(nil, websocket.WSTOPIC_RAW_BLOCK, resp)

		resp["Action"] = "sendjsonblock"
		resp["Result"] = bcomn.GetBlockInfo(&block)
		ws.BroadcastToSubscribers(nil, websocket.WSTOPIC_JSON_BLOCK, resp)
	}
}
func pushBlockTransactions(v interface{}) {
	if ws == nil {
		return
	}
	resp := rest.ResponsePack(Err.SUCCESS)
	if block, ok := v.(types.Block); ok {
		resp["Result"] = bcomn.GetBlockTransactions(&block)
		resp["Action"] = "sendblocktxhashs"
		ws.BroadcastToSubscribers(nil, websocket.WSTOPIC_TXHASHS, resp)
	}
}
