 package connect_controller

import (
	"net"

	" github.com/Daironode/aingle/p2pserver/common"
)

// Conn is a net.Conn wrapper to do some clean up when Close.
type Conn struct {
	net.Conn
	addr       string
	listenAddr string
	kid        common.PeerId
	boundIndex int
	connectId  uint64
	controller *ConnectController
}

// Close overwrite net.Conn
// warning: this method will try to lock the controller, be carefull to avoid deadlock
func (self *Conn) Close() error {
	self.controller.logger.Infof("closing connection: peer %s, address: %s", self.kid.ToHexString(), self.addr)

	self.controller.removePeer(self)

	return self.Conn.Close()
}
