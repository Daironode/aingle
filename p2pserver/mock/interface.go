 
package mock

import (
	"net"
	"strconv"

	" github.com/Daironode/aingle/p2pserver/common"
	" github.com/Daironode/aingle/p2pserver/connect_controller"
	" github.com/Daironode/aingle/p2pserver/net/netserver"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/p2pserver/peer"
)

type Network interface {
	// NewListener will gen random ip to listen
	NewListener(id common.PeerId) (string, net.Listener)
	// addr: ip:port
	NewListenerWithAddr(id common.PeerId, addr string) net.Listener

	// NewDialer will gen random source IP
	NewDialer(id common.PeerId) connect_controller.Dialer
	NewDialerWithHost(id common.PeerId, host string) connect_controller.Dialer
	AllowConnect(id1, id2 common.PeerId)
	DeliverRate(percent uint)
}

func NewNode(keyId *common.PeerKeyId, listenAddr string, localInfo *peer.PeerInfo, proto p2p.Protocol, nw Network,
	reservedPeers []string, reserveAddrFilter p2p.AddressFilter, logger common.Logger) *netserver.NetServer {
	var listener net.Listener
	if listenAddr != "" {
		listener = nw.NewListenerWithAddr(keyId.Id, listenAddr)
	} else {
		listenAddr, listener = nw.NewListener(keyId.Id)
	}
	host, port, _ := net.SplitHostPort(listenAddr)
	dialer := nw.NewDialerWithHost(keyId.Id, host)
	localInfo.Addr = listenAddr
	iport, _ := strconv.Atoi(port)
	localInfo.Port = uint16(iport)
	opt := connect_controller.NewConnCtrlOption().MaxInBoundPerIp(10).
		MaxInBound(20).MaxOutBound(20).WithDialer(dialer).ReservedOnly(reservedPeers)
	opt.ReservedPeers = p2p.CombineAddrFilter(opt.ReservedPeers, reserveAddrFilter)
	return netserver.NewCustomNetServer(keyId, localInfo, proto, listener, opt, logger)
}
