
package ontfs

import (
	"encoding/binary"

	" github.com/Daironode/aingle/common"
	"golang.org/x/crypto/sha3"
)

func genRandSlice(listLen uint64, seed []byte, addr [20]byte) []uint16 {
	uint16Slice := make([]uint16, listLen)

	h := sha3.New512()
	h.Write(addr[:])
	h.Write(seed)
	finalRandData := h.Sum(nil)

	for {
		if uint64(len(finalRandData)) >= 2*listLen {
			break
		}
		h.Reset()
		h.Write(addr[:])
		h.Write(finalRandData)
		randData := h.Sum(nil)

		finalRandData = append(finalRandData, randData...)
	}
	finalRandData = finalRandData[0 : 2*listLen]

	for i := uint64(0); i < listLen; i++ {
		uint16Slice[i] = binary.LittleEndian.Uint16(finalRandData[i*2 : (i+1)*2])
	}
	return uint16Slice
}

func sortByRandSlice(values []uint16, nodeAddrList []common.Address) []common.Address {
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if values[i] > values[j] {
				values[i], values[j] = values[j], values[i]
				nodeAddrList[i], nodeAddrList[j] = nodeAddrList[j], nodeAddrList[i]
			}
		}
	}
	return nodeAddrList
}
