
package handlers

import (
	"encoding/json"
	"os"
	"testing"

	" github.com/Daironode/aingle/account"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
)

func TestExportWallet(t *testing.T) {
	exportReq := &ExportAccountReq{}
	data, _ := json.Marshal(exportReq)
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "exportaccount",
		Pwd:    string(pwd),
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	ExportAccount(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("ExportAccount failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
	exportRsp, ok := resp.Result.(*ExportAccountResp)
	if !ok {
		t.Errorf("TestExportWallet resp asset to ExportAccountResp failed")
		return
	}

	wallet, err := account.Open(exportRsp.WalletFile)
	if err != nil {
		t.Errorf("TestExportWallet failed, OpenWallet error:%s", err)
		return
	}
	if wallet.GetAccountNum() != exportRsp.AccountNumber {
		t.Errorf("TestExportWallet failed, account number %d != %d", wallet.GetAccountNum(), exportRsp.AccountNumber)
		return
	}
	os.Remove(exportRsp.WalletFile)
}
