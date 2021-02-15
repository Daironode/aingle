
package utils

import (
	"bytes"
	"encoding/hex"
	"testing"

	" github.com/Daironode/aingle/cmd/abi"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/vm/neovm"
)

func TestParseNativeParam(t *testing.T) {
	paramAbi := []*abi.NativeContractParamAbi{
		{
			Name: "Param1",
			Type: "String",
		},
		{
			Name: "Param2",
			Type: "Int",
		},
		{
			Name: "Param3",
			Type: "Bool",
		},
		{
			Name: "Param4",
			Type: "Address",
		},
		{
			Name: "Param5",
			Type: "Uint256",
		},
		{
			Name: "Param6",
			Type: "Byte",
		},
		{
			Name: "Param7",
			Type: "ByteArray",
		},
		{
			Name: "Param8",
			Type: "Array",
			SubType: []*abi.NativeContractParamAbi{
				{
					Name: "",
					Type: "Int",
				},
			},
		},
		{
			Name: "Param9",
			Type: "Struct",
			SubType: []*abi.NativeContractParamAbi{
				{
					Name: "Param9_0",
					Type: "String",
				},
				{
					Name: "Param9_1",
					Type: "Int",
				},
			},
		},
	}
	addr := common.Address([20]byte{})
	address := addr.ToBase58()

	params := []interface{}{
		"Hello, World",
		"12",
		"true",
		address,
		"a757b22282b43e0852c48feae0892af19e48da8627296ef7a051993afb316b9b",
		"128",
		hex.EncodeToString([]byte("foo")),
		[]interface{}{"1", "2", "3", "4", "5", "6"},
		[]interface{}{"bar", "10"},
	}
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := ParseNativeFuncParam(builder, "", params, paramAbi)
	if err != nil {
		t.Errorf("ParseNativeParam error:%s", err)
		return
	}
}
