 package mock

import (
	"fmt"
	"strings"
	"testing"
	"time"

	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/p2pserver/common"
	" github.com/Daironode/aingle/p2pserver/net/netserver"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/p2pserver/peer"
	"github.com/stretchr/testify/assert"
)

func TestReserved(t *testing.T) {
	log.Info("test reserved start")
	//topo
	/**
	normal —————— normal
	  \  reserved  /
	   \    |     /
	    \  seed  /
	        |
	        |
	      normal
	*/
	N := 4
	net := NewNetwork()
	seedNode := NewReservedNode(nil, net, nil, "seed")

	var nodes []*netserver.NetServer
	go seedNode.Start()
	seedAddr := seedNode.GetHostInfo().Addr
	seedIP := strings.Split(seedAddr, ":")[0]
	for i := 0; i < N; i++ {
		var node *netserver.NetServer
		var reserved []string
		prefix := "norm"
		if i == 0 {
			reserved = []string{seedIP}
			prefix = "resv"
		}
		node = NewReservedNode([]string{seedAddr}, net, reserved, prefix)
		net.AllowConnect(seedNode.GetHostInfo().Id, node.GetHostInfo().Id)
		go node.Start()
		nodes = append(nodes, node)
	}

	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			net.AllowConnect(nodes[i].GetHostInfo().Id, nodes[j].GetHostInfo().Id)
		}
	}

	time.Sleep(time.Second * 10)
	assert.Equal(t, uint32(N), seedNode.GetConnectionCnt())
	assert.Equal(t, uint32(1), nodes[0].GetConnectionCnt())
	for i := 1; i < N; i++ {
		assert.Equal(t, uint32(N-1), nodes[i].GetConnectionCnt(), i)
		assert.False(t, hasPeerId(nodes[i].GetNeighborAddrs(), nodes[0].GetID()))
	}
}

func hasPeerId(pas []common.PeerAddr, id common.PeerId) bool {
	for _, pa := range pas {
		if pa.ID == id {
			return true
		}
	}
	return false
}

func NewReservedNode(seeds []string, net Network, reservedPeers []string, logPrefix string) *netserver.NetServer {
	seedId := common.RandPeerKeyId()
	info := peer.NewPeerInfo(seedId.Id, 0, 0, true, 0,
		0, 0, "1.10", "")
	dis := NewDiscoveryProtocol(seeds, nil)
	dis.RefleshInterval = time.Millisecond * 1000
	context := fmt.Sprintf("peer %s-%s:, ", logPrefix, seedId.Id.ToHexString()[:6])
	logger := common.LoggerWithContext(common.NewGlobalLoggerWrapper(), context)

	rsvFilter := p2p.NoneAddrFilter()
	if len(reservedPeers) == 0 {
		rsvFilter = p2p.AllAddrFilter()
	}
	return NewNode(seedId, "", info, dis, net, reservedPeers, rsvFilter, logger)
}
