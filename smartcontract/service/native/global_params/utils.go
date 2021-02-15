
package global_params

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/config"
	cstates " github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/smartcontract/event"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

const (
	PARAM    = "param"
	TRANSFER = "transfer"
	ADMIN    = "admin"
	OPERATOR = "operator"
)

func getRoleStorageItem(role common.Address) *cstates.StorageItem {
	bf := common.NewZeroCopySink(nil)
	utils.EncodeAddress(bf, role)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func getParamStorageItem(params Params) *cstates.StorageItem {
	return &cstates.StorageItem{Value: common.SerializeToBytes(&params)}
}

func generateParamKey(contract common.Address, valueType paramType) []byte {
	key := append(contract[:], PARAM...)
	key = append(key[:], byte(valueType))
	return key
}

func generateAdminKey(contract common.Address, isTransferAdmin bool) []byte {
	if isTransferAdmin {
		return append(contract[:], TRANSFER...)
	} else {
		return append(contract[:], ADMIN...)
	}
}

func GenerateOperatorKey(contract common.Address) []byte {
	return append(contract[:], OPERATOR...)
}

func getStorageParam(native *native.NativeService, key []byte) (Params, error) {
	item, err := utils.GetStorageItem(native, key)
	params := Params{}
	if err != nil || item == nil {
		return params, err
	}
	err = params.Deserialization(common.NewZeroCopySource(item.Value))
	return params, err
}

func GetStorageRole(native *native.NativeService, key []byte) (common.Address, error) {
	item, err := utils.GetStorageItem(native, key)
	var role common.Address
	if err != nil || item == nil {
		return role, err
	}
	bf := common.NewZeroCopySource(item.Value)
	role, err = utils.DecodeAddress(bf)
	return role, err
}

func NotifyRoleChange(native *native.NativeService, contract common.Address, functionName string,
	newAddr common.Address) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, newAddr.ToBase58()},
		})
}

func NotifyTransferAdmin(native *native.NativeService, contract common.Address, functionName string,
	originAdmin, newAdmin common.Address) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, originAdmin.ToBase58(), newAdmin.ToBase58()},
		})
}

func NotifyParamChange(native *native.NativeService, contract common.Address, functionName string, params Params) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	paramsString := ""
	for _, param := range params {
		paramsString += param.Key + "," + param.Value + ";"
	}
	paramsString = paramsString[:len(paramsString)-1]
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{functionName, paramsString},
		})
}
