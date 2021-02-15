 package wasmvm

import (
	"errors"
	"math"

	" github.com/Daironode/aingle/core/states"
	"github.com/Daironode/aingle-wagon/exec"
)

func storageRead(service *WasmVmService, keybytes []byte, klen uint32, vlen uint32, offset uint32) ([]byte, uint32, error) {
	key := serializeStorageKey(service.ContextRef.CurrentContext().ContractAddress, keybytes)

	raw, err := service.CacheDB.Get(key)
	if err != nil {
		return []byte{}, 0, err
	}

	if raw == nil {
		return []byte{}, math.MaxUint32, nil
	}

	item, err := states.GetValueFromRawStorageItem(raw)
	if err != nil {
		return []byte{}, 0, err
	}

	length := vlen
	itemlen := uint32(len(item))
	if itemlen < vlen {
		length = itemlen
	}

	if uint32(len(item)) < offset {
		return []byte{}, 0, errors.New("offset is invalid")
	}

	return item[offset : offset+length], uint32(len(item)), nil
}

func StorageRead(proc *exec.Process, keyPtr uint32, klen uint32, val uint32, vlen uint32, offset uint32) uint32 {
	self := proc.HostData().(*Runtime)
	self.checkGas(STORAGE_GET_GAS)
	keybytes, err := ReadWasmMemory(proc, keyPtr, klen)
	if err != nil {
		panic(err)
	}

	itemWrite, originLen, err := storageRead(self.Service, keybytes, klen, vlen, offset)
	if err != nil {
		panic(err)
	}

	if originLen != math.MaxUint32 {
		_, err = proc.WriteAt(itemWrite[:], int64(val))

		if err != nil {
			panic(err)
		}
	}

	return originLen
}

func StorageWrite(proc *exec.Process, keyPtr uint32, keyLen uint32, valPtr uint32, valLen uint32) {
	self := proc.HostData().(*Runtime)
	keybytes, err := ReadWasmMemory(proc, keyPtr, keyLen)
	if err != nil {
		panic(err)
	}

	valbytes, err := ReadWasmMemory(proc, valPtr, valLen)
	if err != nil {
		panic(err)
	}

	cost := uint64((len(keybytes)+len(valbytes)-1)/1024+1) * STORAGE_PUT_GAS
	self.checkGas(cost)

	key := serializeStorageKey(self.Service.ContextRef.CurrentContext().ContractAddress, keybytes)

	self.Service.CacheDB.Put(key, states.GenRawStorageItem(valbytes))
}

func StorageDelete(proc *exec.Process, keyPtr uint32, keyLen uint32) {
	self := proc.HostData().(*Runtime)
	self.checkGas(STORAGE_DELETE_GAS)
	keybytes, err := ReadWasmMemory(proc, keyPtr, keyLen)
	if err != nil {
		panic(err)
	}
	key := serializeStorageKey(self.Service.ContextRef.CurrentContext().ContractAddress, keybytes)

	self.Service.CacheDB.Delete(key)
}
