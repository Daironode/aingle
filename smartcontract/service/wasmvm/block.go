 package wasmvm

import (
	"github.com/Daironode/aingle-wagon/exec"
)

func GetCurrentBlockHash(proc *exec.Process, ptr uint32) uint32 {
	self := proc.HostData().(*Runtime)
	self.checkGas(CURRENT_BLOCK_HASH_GAS)
	blockhash := self.Service.BlockHash

	length, err := proc.WriteAt(blockhash[:], int64(ptr))
	if err != nil {
		panic(err)
	}
	return uint32(length)
}
