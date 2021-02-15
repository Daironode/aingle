
package utils

import (
	"bytes"
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/serialization"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
)

func ConcatKey(contract common.Address, args ...[]byte) []byte {
	temp := contract[:]
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func ConcatBytes(args ...[]byte) []byte {
	temp := []byte{}
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func ValidateOwner(native *native.NativeService, address common.Address) error {
	if !native.ContextRef.CheckWitness(address) {
		return errors.NewErr("validateOwner, authentication failed!")
	}
	return nil
}

func GetUint32Bytes(num uint32) ([]byte, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteUint32(bf, num); err != nil {
		return nil, fmt.Errorf("serialization.WriteUint32, serialize uint32 error: %v", err)
	}
	return bf.Bytes(), nil
}

func GetBytesUint32(b []byte) (uint32, error) {
	num, err := serialization.ReadUint32(bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("serialization.ReadUint32, deserialize uint32 error: %v", err)
	}
	return num, nil
}

func GetUint64Bytes(num uint64) ([]byte, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteUint64(bf, num); err != nil {
		return nil, fmt.Errorf("serialization.WriteUint64, serialize uint64 error: %v", err)
	}
	return bf.Bytes(), nil
}

func GetBytesUint64(b []byte) (uint64, error) {
	num, err := serialization.ReadUint64(bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("serialization.ReadUint64, deserialize uint64 error: %v", err)
	}
	return num, nil
}
