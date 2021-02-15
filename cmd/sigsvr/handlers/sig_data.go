
package handlers

import (
	"encoding/hex"
	"encoding/json"

	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	cliutil " github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common/log"
)

type SigDataReq struct {
	RawData string `json:"raw_data"`
}

type SigDataRsp struct {
	SignedData string `json:"signed_data"`
}

func SigData(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigDataReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawData, err := hex.DecodeString(rawReq.RawData)
	if err != nil {
		log.Infof("Cli Qid:%s SigData hex.DecodeString error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	signer, err := req.GetAccount()
	if err != nil {
		log.Infof("Cli Qid:%s SigData GetAccount:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_ACCOUNT_UNLOCK
		return
	}
	sigData, err := cliutil.Sign(rawData, signer)
	if err != nil {
		log.Infof("Cli Qid:%s SigData Sign error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigDataRsp{
		SignedData: hex.EncodeToString(sigData),
	}
}
