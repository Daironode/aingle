
package handlers

import (
	"encoding/json"
	"testing"

	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	" github.com/Daironode/aingle/cmd/utils"
)

func TestSigNeoVMInvokeTx(t *testing.T) {
	defAcc, err := testWallet.GetDefaultAccount(pwd)
	if err != nil {
		t.Errorf("GetDefaultAccount error:%s", err)
		return
	}

	address1 := defAcc.Address.ToHexString()
	invokeReq := &SigNeoVMInvokeTxReq{
		GasPrice: 0,
		GasLimit: 0,
		Address:  address1,
		Params: []interface{}{
			&utils.NeoVMInvokeParam{
				Type:  "string",
				Value: "foo",
			},
			&utils.NeoVMInvokeParam{
				Type: "array",
				Value: []interface{}{
					&utils.NeoVMInvokeParam{
						Type:  "int",
						Value: "0",
					},
					&utils.NeoVMInvokeParam{
						Type:  "bool",
						Value: "true",
					},
				},
			},
		},
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigNeoVMInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:     "t",
		Method:  "signeovminvoketx",
		Params:  data,
		Account: defAcc.Address.ToBase58(),
		Pwd:     string(pwd),
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigNeoVMInvokeTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigNeoVMInvokeTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
