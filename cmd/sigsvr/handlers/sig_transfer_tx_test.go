
package handlers

import (
	"encoding/json"
	"testing"

	" github.com/Daironode/aingle/account"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
)

func TestSigTransferTransaction(t *testing.T) {
	acc := account.NewAccount("")
	defAcc, err := testWallet.GetDefaultAccount(pwd)
	if err != nil {
		t.Errorf("GetDefaultAccount error:%s", err)
		return
	}
	sigReq := &SigTransferTransactionReq{
		GasLimit: 0,
		GasPrice: 0,
		Asset:    "ont",
		From:     defAcc.Address.ToBase58(),
		To:       acc.Address.ToBase58(),
		Amount:   "10",
	}
	data, err := json.Marshal(sigReq)
	if err != nil {
		t.Errorf("json.Marshal SigTransferTransactionReq error:%s", err)
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:     "t",
		Method:  "sigtransfertx",
		Params:  data,
		Account: defAcc.Address.ToBase58(),
		Pwd:     string(pwd),
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigTransferTransaction(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigTransferTransaction failed. ErrorCode:%d", rsp.ErrorCode)
		return
	}
}
