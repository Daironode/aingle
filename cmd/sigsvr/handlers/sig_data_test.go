
package handlers

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
)

func TestSigData(t *testing.T) {
	defAcc, err := testWallet.GetDefaultAccount(pwd)
	if err != nil {
		t.Errorf("GetDefaultAccount error:%s", err)
		return
	}

	rawData := []byte("HelloWorld")
	rawReq := &SigDataReq{
		RawData: hex.EncodeToString(rawData),
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigDataReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:     "t",
		Method:  "sigdata",
		Params:  data,
		Account: defAcc.Address.ToBase58(),
		Pwd:     string(pwd),
	}
	resp := &clisvrcom.CliRpcResponse{}
	SigData(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("SigData failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
}
