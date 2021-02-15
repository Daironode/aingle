 package p2p

import (
	" github.com/Daironode/aingle/p2pserver/message/types"
	" github.com/Daironode/aingle/p2pserver/peer"
)

type Context struct {
	sender  *peer.Peer
	net     P2P
	MsgSize uint32
}

func NewContext(sender *peer.Peer, net P2P, msgSize uint32) *Context {
	return &Context{sender, net, msgSize}
}

func (self *Context) Sender() *peer.Peer {
	return self.sender
}

func (self *Context) Network() P2P {
	return self.net
}

type Protocol interface {
	HandlePeerMessage(ctx *Context, msg types.Message)
	HandleSystemMessage(net P2P, msg SystemMessage)
}

type SystemMessage interface {
	systemMessage()
}

type implSystemMessage struct{}

func (self implSystemMessage) systemMessage() {}

type PeerConnected struct {
	Info *peer.PeerInfo
	implSystemMessage
}

type PeerDisConnected struct {
	Info *peer.PeerInfo
	implSystemMessage
}

type NetworkStart struct {
	implSystemMessage
}

type NetworkStop struct {
	implSystemMessage
}

type HostAddrDetected struct {
	implSystemMessage
	ListenAddr string
}
