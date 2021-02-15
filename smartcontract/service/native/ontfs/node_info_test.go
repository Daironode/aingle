
package ontfs

import (
	"fmt"
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestFsNodeInfo_Serialization(t *testing.T) {
	nodeInfo := FsNodeInfo{
		Pledge:      uint64(10),
		Profit:      uint64(20),
		Volume:      uint64(30),
		RestVol:     uint64(40),
		ServiceTime: uint64(50),
		NodeAddr: common.Address{0x01, 0x02, 0x03, 0x04, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05},
		NodeNetAddr: []byte("111.111.111.111：111"),
	}
	sink := common.NewZeroCopySink(nil)
	nodeInfo.Serialization(sink)

	fmt.Printf("%v", sink.Bytes())

	nodeInfo2 := FsNodeInfo{}
	src := common.NewZeroCopySource(sink.Bytes())
	if err := nodeInfo2.Deserialization(src); err != nil {
		t.Fatal("nodeInfo2 deserialize fail!", err.Error())
	}

	assert.Equal(t, nodeInfo, nodeInfo2)
}
