
// Package p2p provides an network interface
package p2p

import (
	" github.com/Daironode/aingle/p2pserver/common"
	" github.com/Daironode/aingle/p2pserver/message/types"
	" github.com/Daironode/aingle/p2pserver/peer"
)

//P2P represent the net interface of p2p package
type P2P interface {
	Connect(addr string)
	GetHostInfo() *peer.PeerInfo
	GetID() common.PeerId
	GetNeighbors() []*peer.Peer
	GetNeighborAddrs() []common.PeerAddr
	GetConnectionCnt() uint32
	GetMaxPeerBlockHeight() uint64
	GetPeer(id common.PeerId) *peer.Peer
	SetHeight(uint64)
	Send(p *peer.Peer, msg types.Message) error
	SendTo(p common.PeerId, msg types.Message)
	GetOutConnRecordLen() uint
	Broadcast(msg types.Message)
	IsOwnAddress(addr string) bool
}
