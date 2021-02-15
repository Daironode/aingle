
package handlers

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/Daironode/aingle-crypto/keypair"
	"github.com/Daironode/aingle-crypto/signature"
	clisvrcom " github.com/Daironode/aingle/cmd/sigsvr/common"
	" github.com/Daironode/aingle/cmd/utils"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
	"github.com/stretchr/testify/assert"
)

func TestSigMutilRawTransaction(t *testing.T) {
	acc1, err := clisvrcom.DefWalletStore.NewAccountData(keypair.PK_ECDSA, keypair.P256, signature.SHA256withECDSA, pwd)
	if err != nil {
		t.Errorf("wallet.NewAccount error:%s", err)
		return
	}
	clisvrcom.DefWalletStore.AddAccountData(acc1)
	acc2, err := clisvrcom.DefWalletStore.NewAccountData(keypair.PK_ECDSA, keypair.P256, signature.SHA256withECDSA, pwd)
	if err != nil {
		t.Errorf("wallet.NewAccount error:%s", err)
		return
	}
	clisvrcom.DefWalletStore.AddAccountData(acc2)

	pkData, _ := hex.DecodeString(acc1.PubKey)
	acc1PubKey, _ := keypair.DeserializePublicKey(pkData)
	pkData, _ = hex.DecodeString(acc2.PubKey)
	acc2PubKey, _ := keypair.DeserializePublicKey(pkData)

	pubKeys := []keypair.PublicKey{acc1PubKey, acc2PubKey}
	m := 2
	fromAddr, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		t.Errorf("TestSigMutilRawTransaction AddressFromMultiPubKeys error:%s", err)
		return
	}
	tx, err := utils.TransferTx(0, 0, "ont", fromAddr.ToBase58(), acc1.Address, 10)
	if err != nil {
		t.Errorf("TransferTx error:%s", err)
		return
	}
	immut, err := tx.IntoImmutable()
	assert.Nil(t, err)
	sink := common.ZeroCopySink{}
	immut.Serialization(&sink)

	rawReq := &SigMutilRawTransactionReq{
		RawTx:   hex.EncodeToString(sink.Bytes()),
		M:       m,
		PubKeys: []string{acc1.PubKey, acc2.PubKey},
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigRawTransactionReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:     "t",
		Method:  "sigmutilrawtx",
		Params:  data,
		Account: acc1.Address,
		Pwd:     string(pwd),
	}
	resp := &clisvrcom.CliRpcResponse{}
	SigMutilRawTransaction(req, resp)
	if resp.ErrorCode != clisvrcom.CLIERR_OK {
		t.Errorf("SigMutilRawTransaction failed,ErrorCode:%d ErrorString:%s", resp.ErrorCode, resp.ErrorInfo)
		return
	}

	req.Account = acc2.Address
	SigMutilRawTransaction(req, resp)
	if resp.ErrorCode != clisvrcom.CLIERR_OK {
		t.Errorf("SigMutilRawTransaction failed,ErrorCode:%d ErrorString:%s", resp.ErrorCode, resp.ErrorInfo)
		return
	}
}
