
package actor

import (
	" github.com/Daironode/aingle/p2pserver/common"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
)

var netServer p2p.P2P

func SetNetServer(p2p p2p.P2P) {
	netServer = p2p
}

//GetConnectionCnt from netSever actor
func GetConnectionCnt() uint32 {
	if netServer == nil {
		return 1
	}

	return netServer.GetConnectionCnt()
}

//GetMaxPeerBlockHeight from netSever actor
func GetMaxPeerBlockHeight() uint64 {
	if netServer == nil {
		return 1
	}
	return netServer.GetMaxPeerBlockHeight()
}

//GetNeighborAddrs from netSever actor
func GetNeighborAddrs() []common.PeerAddr {
	if netServer == nil {
		return []common.PeerAddr{}
	}
	return netServer.GetNeighborAddrs()
}

//GetNodePort from netSever actor
func GetNodePort() uint16 {
	if netServer == nil {
		return 0
	}
	return netServer.GetHostInfo().Port
}

//GetID from netSever actor
func GetID() common.PeerId {
	if netServer == nil {
		return common.PeerId{}
	}
	return netServer.GetID()
}

//GetRelayState from netSever actor
func GetRelayState() bool {
	if netServer == nil {
		return false
	}
	return netServer.GetHostInfo().Relay
}

//GetVersion from netSever actor
func GetVersion() uint32 {
	if netServer == nil {
		return 0
	}
	return netServer.GetHostInfo().Version
}

//GetNodeType from netSever actor
func GetNodeType() uint64 {
	if netServer == nil {
		return 0
	}
	return netServer.GetHostInfo().Services
}
