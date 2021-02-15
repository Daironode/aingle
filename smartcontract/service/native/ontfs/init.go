 
package ontfs

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

func InitFs() {
	native.Contracts[utils.OntFSContractAddress] = RegisterFsContract
}

func RegisterFsContract(native *native.NativeService) {
	//native.Register(FS_SET_GLOBAL_PARAM, FsSetGlobalParam)
	native.Register(FS_GET_GLOBAL_PARAM, FsGetGlobalParam)

	native.Register(FS_NODE_REGISTER, FsNodeRegister)
	native.Register(FS_NODE_QUERY, FsNodeQuery)
	native.Register(FS_NODE_UPDATE, FsNodeUpdate)
	native.Register(FS_NODE_CANCEL, FsNodeCancel)
	native.Register(FS_FILE_PROVE, FsFileProve)
	native.Register(FS_NODE_WITHDRAW_PROFIT, FsNodeWithdrawProfit)

	native.Register(FS_GET_NODE_LIST, FsGetNodeInfoList)
	native.Register(FS_GET_PDP_INFO_LIST, FsGetPdpInfoList)

	native.Register(FS_CHALLENGE, FsChallenge)
	native.Register(FS_RESPONSE, FsResponse)
	native.Register(FS_JUDGE, FsJudge)
	native.Register(FS_GET_CHALLENGE, FsGetChallenge)
	native.Register(FS_GET_FILE_CHALLENGE_LIST, FsGetFileChallengeList)
	native.Register(FS_GET_NODE_CHALLENGE_LIST, FsGetNodeChallengeList)

	native.Register(FS_STORE_FILES, FsStoreFiles)
	native.Register(FS_RENEW_FILES, FsRenewFiles)
	native.Register(FS_DELETE_FILES, FsDeleteFiles)
	native.Register(FS_TRANSFER_FILES, FsTransferFiles)

	native.Register(FS_GET_FILE_INFO, FsGetFileInfo)
	native.Register(FS_GET_FILE_LIST, FsGetFileHashList)

	native.Register(FS_READ_FILE_PLEDGE, FsReadFilePledge)
	native.Register(FS_READ_FILE_SETTLE, FsReadFileSettle)
	native.Register(FS_GET_READ_PLEDGE, FsGetReadPledge)

	native.Register(FS_CREATE_SPACE, FsCreateSpace)
	native.Register(FS_DELETE_SPACE, FsDeleteSpace)
	native.Register(FS_UPDATE_SPACE, FsUpdateSpace)
	native.Register(FS_GET_SPACE_INFO, FsGetSpaceInfo)
}

//To enable administrators to adjust global parameters
func FsSetGlobalParam(native *native.NativeService) ([]byte, error) {
	var globalParam FsGlobalParam
	if err := CheckOntFsAvailability(native); err != nil {
		return utils.BYTE_FALSE, err
	}

	infoSource := common.NewZeroCopySource(native.Input)
	if err := globalParam.Deserialization(infoSource); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Init] FsSetGlobalParam Deserialization error!")
	}
	setGlobalParam(native, &globalParam)
	return utils.BYTE_TRUE, nil
}

func FsGetGlobalParam(native *native.NativeService) ([]byte, error) {
	if err := CheckOntFsAvailability(native); err != nil {
		return utils.BYTE_FALSE, err
	}
	globalParam, err := getGlobalParam(native)
	if err != nil || globalParam == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Init] FsGetGlobalParam error!")
	}
	sink := common.NewZeroCopySink(nil)
	globalParam.Serialization(sink)

	return EncRet(true, sink.Bytes()), nil
}
