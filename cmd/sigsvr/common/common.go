
package common

import (
	"encoding/json"
	"fmt"

	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/cmd/sigsvr/store"
)

var DefWalletStore *store.WalletStore

type CliRpcRequest struct {
	Qid     string          `json:"qid"`
	Params  json.RawMessage `json:"params"`
	Account string          `json:"account"`
	Pwd     string          `json:"pwd"`
	Method  string          `json:"method"`
}

func (this *CliRpcRequest) GetAccount() (*account.Account, error) {
	var acc *account.Account
	var err error

	pwd := []byte(this.Pwd)
	if this.Pwd == "" {
		return nil, fmt.Errorf("pwd cannot empty")
	}
	if this.Account == "" {
		return nil, fmt.Errorf("account cannot empty")
	}
	acc, err = DefWalletStore.GetAccountByAddress(this.Account, pwd)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, fmt.Errorf("cannot find account by address: %s", this.Account)
	}
	return acc, nil
}

type CliRpcResponse struct {
	Qid       string      `json:"qid"`
	Method    string      `json:"method"`
	Result    interface{} `json:"result"`
	ErrorCode int         `json:"error_code"`
	ErrorInfo string      `json:"error_info"`
}
