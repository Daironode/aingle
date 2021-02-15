
package handlers

import (
	"encoding/json"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	"github.com/Daironode/aingle-crypto/signature"
	" github.com/Daironode/aingle/cmd/abi"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	nutils " github.com/Daironode/aingle/smartcontract/service/native/utils"
)

func TestSigNativeInvokeTx(t *testing.T) {
	defAcc, err := testWallet.GetDefaultAccount(pwd)
	if err != nil {
		t.Errorf("GetDefaultAccount error:%s", err)
		return
	}
	acc1, err := clisvrcom.DefWalletStore.NewAccountData(keypair.PK_ECDSA, keypair.P256, signature.SHA256withECDSA, pwd)
	if err != nil {
		t.Errorf("wallet.NewAccount error:%s", err)
		return
	}
	clisvrcom.DefWalletStore.AddAccountData(acc1)
	invokeReq := &SigNativeInvokeTxReq{
		GasPrice: 0,
		GasLimit: 40000,
		Address:  nutils.OntContractAddress.ToHexString(),
		Method:   "transfer",
		Version:  0,
		Params: []interface{}{
			[]interface{}{
				[]interface{}{
					defAcc.Address.ToBase58(),
					acc1.Address,
					"10000000000",
				},
			},
		},
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigNativeInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:     "t",
		Method:  "signativeinvoketx",
		Params:  data,
		Account: acc1.Address,
		Pwd:     string(pwd),
	}
	rsp := &clisvrcom.CliRpcResponse{}
	abiPath := "../../abi/native_abi_script"
	abi.DefAbiMgr.Init(abiPath)
	SigNativeInvokeTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigNativeInvokeTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
