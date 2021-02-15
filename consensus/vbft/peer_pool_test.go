
package vbft

import (
	"testing"
	"time"

	" github.com/Daironode/aingle/common"
	vconfig " github.com/Daironode/aingle/consensus/vbft/config"
)

func constructPeerPool(connect bool) *PeerPool {
	peer := &Peer{
		Index:          uint32(1),
		LastUpdateTime: time.Unix(0, 0),
		connected:      connect,
	}
	peers := make(map[uint32]*Peer)
	peers[1] = peer
	peerpool := &PeerPool{
		maxSize: int(3),
		configs: make(map[uint32]*vconfig.PeerConfig),
		IDMap:   make(map[string]uint32),
		peers:   peers,
	}
	return peerpool
}
func TestIsNewPeer(t *testing.T) {
	peerpool := constructPeerPool(false)
	isnew := peerpool.isNewPeer(uint32(2))
	t.Logf("TestIsNewPeer: %v", isnew)
}

func TestAddPeer(t *testing.T) {
	nodeId := "120202c924ed1a67fd1719020ce599d723d09d48362376836e04b0be72dfe825e24d81"
	peerconfig := &vconfig.PeerConfig{
		Index: uint32(1),
		ID:    nodeId,
	}
	peerpool := constructPeerPool(false)
	res := peerpool.addPeer(peerconfig)
	t.Logf("TestAddPeer : %v", res)
}

func TestGetActivePeerCount(t *testing.T) {
	peerpool := constructPeerPool(true)
	count := peerpool.getActivePeerCount()
	t.Logf("TestGetActivePeerCount count:%v", count)
}

func TestPeerHandshake(t *testing.T) {
	nodeId := "120202c924ed1a67fd1719020ce599d723d09d48362376836e04b0be72dfe825e24d81"
	peerconfig := &vconfig.PeerConfig{
		Index: uint32(1),
		ID:    nodeId,
	}
	peerpool := constructPeerPool(false)
	peerpool.addPeer(peerconfig)
	handshakemsg := &peerHandshakeMsg{
		CommittedBlockNumber: uint32(2),
		CommittedBlockHash:   common.Uint256{},
		CommittedBlockLeader: uint32(1),
	}
	peerpool.peerHandshake(uint32(1), handshakemsg)
}

func TestPeerHeartbeat(t *testing.T) {
	nodeId := "120202c924ed1a67fd1719020ce599d723d09d48362376836e04b0be72dfe825e24d81"
	peerconfig := &vconfig.PeerConfig{
		Index: uint32(1),
		ID:    nodeId,
	}
	peerpool := constructPeerPool(false)
	peerpool.addPeer(peerconfig)
	heartbeatmsg := &peerHeartbeatMsg{
		CommittedBlockNumber: uint32(2),
		CommittedBlockHash:   common.Uint256{},
		CommittedBlockLeader: uint32(1),
		ChainConfigView:      uint32(1),
	}
	peerpool.peerHeartbeat(uint32(1), heartbeatmsg)
}

func TestGetNeighbours(t *testing.T) {
	peerpool := constructPeerPool(true)
	peers := peerpool.getNeighbours()
	t.Logf("TestGetNeighbours: %d", len(peers))
}

func TestGetPeerIndex(t *testing.T) {
	nodeId := "12020298fe9f22e9df64f6bfcc1c2a14418846cffdbbf510d261bbc3fa6d47073df9a2"
	peerconfig := &vconfig.PeerConfig{
		Index: uint32(1),
		ID:    nodeId,
	}
	peerpool := constructPeerPool(false)
	peerpool.addPeer(peerconfig)
	idx, present := peerpool.GetPeerIndex(nodeId)
	if !present {
		t.Errorf("TestGetPeerIndex is not exist: %d", idx)
		return
	}
	t.Logf("TestGetPeerIndex: %d,%v", idx, present)
}

func TestGetPeer(t *testing.T) {
	nodeId := "12020298fe9f22e9df64f6bfcc1c2a14418846cffdbbf510d261bbc3fa6d47073df9a2"
	peerconfig := &vconfig.PeerConfig{
		Index: uint32(1),
		ID:    nodeId,
	}
	peerpool := constructPeerPool(false)
	peerpool.addPeer(peerconfig)
	peer := peerpool.getPeer(uint32(1))
	if peer == nil {
		t.Errorf("TestGetPeer failed peer is nil")
		return
	}
	t.Logf("TestGetPeer: %v", peer.Index)
}
