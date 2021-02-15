
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	" github.com/Daironode/aingle/common/config"
	rpcerr " github.com/Daironode/aingle/http/base/error"
)

//JsonRpc version
const JSON_RPC_VERSION = "2.0"

const (
	ERROR_INVALID_PARAMS   = rpcerr.INVALID_PARAMS
	ERROR_AINGLE_COMMON  = 10000
	ERROR_AINGLE_SUCCESS = 0
)

type AIngleError struct {
	ErrorCode int64
	Error     error
}

func NewAIngleError(err error, errCode ...int64) *AIngleError {
	ontErr := &AIngleError{Error: err}
	if len(errCode) > 0 {
		ontErr.ErrorCode = errCode[0]
	} else {
		ontErr.ErrorCode = ERROR_AINGLE_COMMON
	}
	if err == nil {
		ontErr.ErrorCode = ERROR_AINGLE_SUCCESS
	}
	return ontErr
}

//JsonRpcRequest object in rpc
type JsonRpcRequest struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//JsonRpcResponse object response for JsonRpcRequest
type JsonRpcResponse struct {
	Error  int64           `json:"error"`
	Desc   string          `json:"desc"`
	Result json.RawMessage `json:"result"`
}

func sendRpcRequest(method string, params []interface{}) ([]byte, *AIngleError) {
	rpcReq := &JsonRpcRequest{
		Version: JSON_RPC_VERSION,
		Id:      "cli",
		Method:  method,
		Params:  params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, NewAIngleError(fmt.Errorf("JsonRpcRequest json.Marshal error:%s", err))
	}

	addr := fmt.Sprintf("http://localhost:%d", config.DefConfig.Rpc.HttpJsonPort)
	resp, err := http.Post(addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, NewAIngleError(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, NewAIngleError(fmt.Errorf("read rpc response body error:%s", err))
	}
	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, NewAIngleError(fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err))
	}
	if rpcRsp.Error != 0 {
		return nil, NewAIngleError(fmt.Errorf("\n %s ", string(body)), rpcRsp.Error)
	}
	return rpcRsp.Result, nil
}
