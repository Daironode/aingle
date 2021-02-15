
package handlers

import (
	"encoding/hex"
	"encoding/json"

	"github.com/Daironode/aingle-crypto/keypair"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	cliutil " github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
)

type SigRawTransactionReq struct {
	RawTx string `json:"raw_tx"`
}

type SigRawTransactionRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigRawTransaction(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigRawTransactionReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTxData, err := hex.DecodeString(rawReq.RawTx)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction hex.DecodeString error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	tmpTx, err := types.TransactionFromRawBytes(rawTxData)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction tx Deserialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}
	mutable, err := tmpTx.IntoMutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction tx IntoMutable error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}
	signer, err := req.GetAccount()
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction GetAccount:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_ACCOUNT_UNLOCK
		return
	}
	var emptyAddress = common.Address{}
	if mutable.Payer == emptyAddress {
		mutable.Payer = signer.Address
	}

	txHash := mutable.Hash()
	sigData, err := cliutil.Sign(txHash.ToArray(), signer)
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction Sign error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	if len(mutable.Sigs) == 0 {
		mutable.Sigs = make([]types.Sig, 0)
	}
	mutable.Sigs = append(mutable.Sigs, types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sigData},
	})

	rawTx, err := mutable.IntoImmutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigRawTransaction tx IntoImmutable error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	sink := common.ZeroCopySink{}
	rawTx.Serialization(&sink)
	resp.Result = &SigRawTransactionRsp{
		SignedTx: hex.EncodeToString(sink.Bytes()),
	}
}
