
package handlers

import (
	"encoding/json"

	"github.com/Daironode/aingle-crypto/keypair"
	s "github.com/Daironode/aingle-crypto/signature"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	" github.com/Daironode/aingle/common/log"
)

type CreateAccountReq struct {
}

type CreateAccountRsp struct {
	Account string `json:"account"`
}

func CreateAccount(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	pwd := req.Pwd
	if pwd == "" {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		resp.ErrorInfo = "pwd cannot empty"
		return
	}
	accData, err := clisvrcom.DefWalletStore.NewAccountData(keypair.PK_ECDSA, keypair.P256, s.SHA256withECDSA, []byte(pwd))
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		resp.ErrorInfo = "create wallet failed"
		log.Errorf("CreateAccount Qid:%s NewAccountData error:%s", req.Qid, err)
		return
	}
	_, err = clisvrcom.DefWalletStore.AddAccountData(accData)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		resp.ErrorInfo = "create wallet failed"
		log.Errorf("CreateAccount Qid:%s AddAccountData error:%s", req.Qid, err)
		return
	}
	resp.Result = &CreateAccountRsp{
		Account: accData.Address,
	}

	data, _ := json.Marshal(accData)
	log.Infof("[CreateAccount]%s", data)
}
