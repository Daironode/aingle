 package handlers

import (
	"testing"

	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
)

func TestCreateAccount(t *testing.T) {
	walletStore := clisvrcom.DefWalletStore
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "createaccount",
		Pwd:    string(pwd),
	}
	resp := &clisvrcom.CliRpcResponse{}
	CreateAccount(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("CreateAccount failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
	createRsp, ok := resp.Result.(*CreateAccountRsp)
	if !ok {
		t.Errorf("CreateAccount resp asset to CreateAccountRsp failed")
		return
	}
	_, err := walletStore.GetAccountByAddress(createRsp.Account, pwd)
	if err != nil {
		t.Errorf("Test CreateAccount failed GetAccountByAddress error:%s", err)
		return
	}
}
