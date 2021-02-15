 package ontid

import (
	"errors"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

func addProof(srvc *native.NativeService) ([]byte, error) {
	return utils.BYTE_FALSE, errors.New("property \"proof\" in ONT ID document is not supported yet")
}

func getProof(srvc *native.NativeService, encId []byte) (string, error) {
	key := append(encId, FIELD_PROOF)
	proofStore, err := utils.GetStorageItem(srvc, key)
	if err != nil {
		return "", errors.New("getProof error:" + err.Error())
	}
	if proofStore == nil {
		return "", nil
	}
	source := common.NewZeroCopySource(proofStore.Value)
	proof, err := utils.DecodeVarBytes(source)
	if err != nil {
		return "", errors.New("DecodeVarBytes error:" + err.Error())
	}
	return string(proof), nil
}
