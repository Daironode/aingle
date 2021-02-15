
package handlers

import (
	"encoding/hex"
	"encoding/json"

	"github.com/Daironode/aingle-crypto/keypair"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	cliutil " github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/constants"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
)

type SigMutilRawTransactionReq struct {
	RawTx   string   `json:"raw_tx"`
	M       int      `json:"m"`
	PubKeys []string `json:"pub_keys"`
}

type SigMutilRawTransactionRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigMutilRawTransaction(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigMutilRawTransactionReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	numkeys := len(rawReq.PubKeys)
	if rawReq.M <= 0 || numkeys < rawReq.M || numkeys <= 1 || numkeys > constants.MULTI_SIG_MAX_PUBKEY_SIZE {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	rawTxData, err := hex.DecodeString(rawReq.RawTx)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction hex.DecodeString error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}

	tmpTx, err := types.TransactionFromRawBytes(rawTxData)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction TransactionFromRawBytes error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}
	mutTx, err := tmpTx.IntoMutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction IntoMutable error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_TX
		return
	}

	pubKeys := make([]keypair.PublicKey, 0, len(rawReq.PubKeys))
	for _, pkStr := range rawReq.PubKeys {
		pkData, err := hex.DecodeString(pkStr)
		if err != nil {
			log.Info("Cli Qid:%s SigMutilRawTransaction pk hex.DecodeString error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		pk, err := keypair.DeserializePublicKey(pkData)
		if err != nil {
			log.Info("Cli Qid:%s SigMutilRawTransaction keypair.DeserializePublicKey error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		pubKeys = append(pubKeys, pk)
	}

	signer, err := req.GetAccount()
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction GetAccount:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_ACCOUNT_UNLOCK
		return
	}
	err = cliutil.MultiSigTransaction(mutTx, uint16(rawReq.M), pubKeys, signer)
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction MultiSigTransaction error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	tmpTx, err = mutTx.IntoImmutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigMutilRawTransaction tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	sink := common.ZeroCopySink{}
	tmpTx.Serialization(&sink)
	resp.Result = &SigRawTransactionRsp{
		SignedTx: hex.EncodeToString(sink.Bytes()),
	}
}
