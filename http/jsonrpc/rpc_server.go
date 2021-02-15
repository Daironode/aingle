
// Package jsonrpc privides a function to start json rpc server
package jsonrpc

import (
	"fmt"
	"net/http"
	"strconv"

	cfg " github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/http/base/rpc"
)

func StartRPCServer() error {
	log.Debug()
	http.HandleFunc("/", rpc.Handle)
	rpc.HandleFunc("getbestblockhash", GetBestBlockHash)
	rpc.HandleFunc("getblock", GetBlock)
	rpc.HandleFunc("getblockcount", GetBlockCount)
	rpc.HandleFunc("getblockhash", GetBlockHash)
	rpc.HandleFunc("getconnectioncount", GetConnectionCount)
	rpc.HandleFunc("getsyncstatus", GetSyncStatus)
	//HandleFunc("getrawmempool", GetRawMemPool)

	rpc.HandleFunc("getrawtransaction", GetRawTransaction)
	rpc.HandleFunc("sendrawtransaction", SendRawTransaction)
	rpc.HandleFunc("getstorage", GetStorage)
	rpc.HandleFunc("getversion", GetNodeVersion)
	rpc.HandleFunc("getnetworkid", GetNetworkId)

	rpc.HandleFunc("getcontractstate", GetContractState)
	rpc.HandleFunc("getmempooltxcount", GetMemPoolTxCount)
	rpc.HandleFunc("getmempooltxstate", GetMemPoolTxState)
	rpc.HandleFunc("getmempooltxhashlist", GetMemPoolTxHashList)
	rpc.HandleFunc("getsmartcodeevent", GetSmartCodeEvent)
	rpc.HandleFunc("getblockheightbytxhash", GetBlockHeightByTxHash)

	rpc.HandleFunc("getbalance", GetBalance)
	rpc.HandleFunc("getoep4balance", GetOep4Balance)
	rpc.HandleFunc("getallowance", GetAllowance)
	rpc.HandleFunc("getmerkleproof", GetMerkleProof)
	rpc.HandleFunc("getblocktxsbyheight", GetBlockTxsByHeight)
	rpc.HandleFunc("getgasprice", GetGasPrice)
	rpc.HandleFunc("getunboundong", GetUnboundOng)
	rpc.HandleFunc("getgrantong", GetGrantOng)

	rpc.HandleFunc("getcrosschainmsg", GetCrossChainMsg)
	rpc.HandleFunc("getcrossstatesproof", GetCrossStatesProof)

	err := http.ListenAndServe(":"+strconv.Itoa(int(cfg.DefConfig.Rpc.HttpJsonPort)), nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe error:%s", err)
	}
	return nil
}
