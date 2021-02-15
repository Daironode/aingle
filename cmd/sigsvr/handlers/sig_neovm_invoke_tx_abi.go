
package handlers

import (
	"encoding/hex"
	"encoding/json"

	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	cliutil " github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
	httpcom " github.com/Daironode/aingle/http/base/common"
)

type SigNeoVMInvokeTxAbiReq struct {
	GasPrice    uint64          `json:"gas_price"`
	GasLimit    uint64          `json:"gas_limit"`
	Address     string          `json:"address"`
	Method      string          `json:"method"`
	Params      []string        `json:"params"`
	Payer       string          `json:"payer"`
	ContractAbi json.RawMessage `json:"contract_abi"`
}

type SigNeoVMInvokeTxAbiRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigNeoVMInvokeAbiTx(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigNeoVMInvokeTxAbiReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		log.Infof("SigNeoVMInvokeAbiTx json.Unmarshal SigNeoVMInvokeTxAbiReq:%s error:%s", req.Params, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	contractAbi, err := cliutil.NewNeovmContractAbi(rawReq.ContractAbi)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_ABI_UNMATCH
		resp.ErrorInfo = err.Error()
		return
	}
	funcAbi := contractAbi.GetFunc(rawReq.Method)
	if funcAbi == nil {
		resp.ErrorCode = clisvrcom.CLIERR_ABI_NOT_FOUND
		return
	}
	invokParams, err := cliutil.ParseNeovmFunc(rawReq.Params, funcAbi)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_ABI_UNMATCH
		resp.ErrorInfo = err.Error()
		return
	}
	contAddr, err := common.AddressFromHexString(rawReq.Address)
	if err != nil {
		log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx AddressParseFromBytes:%s error:%s", req.Qid, rawReq.Address, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	mutable, err := httpcom.NewNeovmInvokeTransaction(rawReq.GasPrice, rawReq.GasLimit, contAddr, invokParams)
	if err != nil {
		log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx InvokeNeoVMContractTx error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	if rawReq.Payer != "" {
		payerAddress, err := common.AddressFromBase58(rawReq.Payer)
		if err != nil {
			log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx AddressFromBase58 error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		mutable.Payer = payerAddress
	}
	signer, err := req.GetAccount()
	if err != nil {
		log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx GetAccount:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_ACCOUNT_UNLOCK
		return
	}
	err = cliutil.SignTransaction(signer, mutable)
	if err != nil {
		log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx SignTransaction error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}

	tx, err := mutable.IntoImmutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigNeoVMInvokeAbiTx tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigNeoVMInvokeTxAbiRsp{
		SignedTx: hex.EncodeToString(common.SerializeToBytes(tx)),
	}
}
