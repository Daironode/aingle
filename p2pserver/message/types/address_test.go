
package types

import (
	"bytes"
	"net"
	"testing"

	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
	"github.com/stretchr/testify/assert"
)

func MessageTest(t *testing.T, msg Message) {
	sink := common.NewZeroCopySink(nil)
	WriteMessage(sink, msg)

	demsg, _, err := ReadMessage(bytes.NewBuffer(sink.Bytes()))
	assert.Nil(t, err)

	assert.Equal(t, msg, demsg)
}

func TestAddressSerializationDeserialization(t *testing.T) {
	var msg Addr
	var addr [16]byte
	ip := net.ParseIP("192.168.0.1")
	ip.To16()
	copy(addr[:], ip[:16])
	id := comm.PseudoPeerIdFromUint64(987654321)
	nodeAddr := comm.PeerAddr{
		Time:          12345678,
		Services:      100,
		IpAddr:        addr,
		Port:          8080,
		ConsensusPort: 8081,
		ID:            id,
	}
	msg.NodeAddrs = append(msg.NodeAddrs, nodeAddr)

	MessageTest(t, &msg)
}
