
package mock

import (
	"errors"
	"net"

	" github.com/Daironode/aingle/p2pserver/common"
	" github.com/Daironode/aingle/p2pserver/connect_controller"
)

type dialer struct {
	id      common.PeerId
	address string
	network *network
}

var _ connect_controller.Dialer = &dialer{}

func (d *dialer) Dial(nodeAddr string) (net.Conn, error) {
	d.network.Lock()
	defer d.network.Unlock()
	l, exist := d.network.listeners[nodeAddr]

	if !exist {
		return nil, errors.New("can not be reached")
	}

	if _, allow := d.network.canEstablish[combineKey(d.id, l.id)]; !allow {
		return nil, errors.New("can not be reached")
	}

	c, s := net.Pipe()

	cw := connWraper{c, d.address, d.network, l.address}
	sw := connWraper{s, l.address, d.network, d.address}
	l.PushToAccept(sw)

	return cw, nil
}

func (n *network) NewDialer(pid common.PeerId) connect_controller.Dialer {
	host := n.nextFakeIP()
	return n.NewDialerWithHost(pid, host)
}

func (n *network) NewDialerWithHost(pid common.PeerId, host string) connect_controller.Dialer {
	addr := host + ":" + n.nextPortString()

	d := &dialer{
		id:      pid,
		address: addr,
		network: n,
	}

	return d
}

func (d *dialer) ID() common.PeerId {
	return d.id
}
