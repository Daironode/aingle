
package utils

import (
	"bytes"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/serialization"
	cstates " github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
)

func GetStorageItem(native *native.NativeService, key []byte) (*cstates.StorageItem, error) {
	store, err := native.CacheDB.Get(key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[GetStorageItem] storage error!")
	}
	if store == nil {
		return nil, nil
	}
	item := new(cstates.StorageItem)
	err = item.Deserialization(common.NewZeroCopySource(store))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[GetStorageItem] instance doesn't StorageItem!")
	}
	return item, nil
}

func GetStorageUInt64(native *native.NativeService, key []byte) (uint64, error) {
	item, err := GetStorageItem(native, key)
	if err != nil {
		return 0, err
	}
	if item == nil {
		return 0, nil
	}
	v, err := serialization.ReadUint64(bytes.NewBuffer(item.Value))
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetStorageUInt32(native *native.NativeService, key []byte) (uint32, error) {
	item, err := GetStorageItem(native, key)
	if err != nil {
		return 0, err
	}
	if item == nil {
		return 0, nil
	}
	v, err := serialization.ReadUint32(bytes.NewBuffer(item.Value))
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GenUInt64StorageItem(value uint64) *cstates.StorageItem {
	sink := common.NewZeroCopySink(nil)
	sink.WriteUint64(value)
	return &cstates.StorageItem{Value: sink.Bytes()}
}

func GenUInt32StorageItem(value uint32) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	serialization.WriteUint32(bf, value)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func PutBytes(native *native.NativeService, key []byte, value []byte) {
	native.CacheDB.Put(key, cstates.GenRawStorageItem(value))
}

func GetStorageVarBytes(native *native.NativeService, key []byte) ([]byte, error) {
	item, err := GetStorageItem(native, key)
	if err != nil {
		return []byte{}, err
	}
	if item == nil {
		return []byte{}, nil
	}
	v, err := serialization.ReadVarBytes(bytes.NewBuffer(item.Value))
	if err != nil {
		return []byte{}, err
	}
	return v, nil
}

func GenVarBytesStorageItem(value []byte) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	serialization.WriteVarBytes(bf, value)
	return &cstates.StorageItem{Value: bf.Bytes()}
}
