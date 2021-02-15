 
// Package localrpc privides a function to start local rpc server
package localrpc

import (
	"fmt"
	"net/http"
	"strconv"

	cfg " github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/http/base/rpc"
)

const (
	LOCAL_HOST string = "127.0.0.1"
	LOCAL_DIR  string = "/local"
)

func StartLocalServer() error {
	log.Debug()
	http.HandleFunc(LOCAL_DIR, rpc.Handle)

	rpc.HandleFunc("getneighbor", GetNeighbor)
	rpc.HandleFunc("getnodestate", GetNodeState)
	rpc.HandleFunc("startconsensus", StartConsensus)
	rpc.HandleFunc("stopconsensus", StopConsensus)
	rpc.HandleFunc("setdebuginfo", SetDebugInfo)

	// TODO: only listen to local host
	err := http.ListenAndServe(LOCAL_HOST+":"+strconv.Itoa(int(cfg.DefConfig.Rpc.HttpLocalPort)), nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe error:%s", err)
	}
	return nil
}
