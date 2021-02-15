 package connect_controller

import (
	" github.com/Daironode/aingle/common/config"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
)

type ConnCtrlOption struct {
	MaxConnOutBound     uint
	MaxConnInBound      uint
	MaxConnInBoundPerIP uint
	ReservedPeers       p2p.AddressFilter // enabled if not empty
	dialer              Dialer
}

func NewConnCtrlOption() ConnCtrlOption {
	return ConnCtrlOption{
		MaxConnInBound:      config.DEFAULT_MAX_CONN_IN_BOUND,
		MaxConnOutBound:     config.DEFAULT_MAX_CONN_OUT_BOUND,
		MaxConnInBoundPerIP: config.DEFAULT_MAX_CONN_IN_BOUND_FOR_SINGLE_IP,
		ReservedPeers:       p2p.AllAddrFilter(),
		dialer:              &noTlsDialer{},
	}
}

func (self ConnCtrlOption) MaxOutBound(n uint) ConnCtrlOption {
	self.MaxConnOutBound = n
	return self
}

func (self ConnCtrlOption) MaxInBound(n uint) ConnCtrlOption {
	self.MaxConnInBound = n
	return self
}

func (self ConnCtrlOption) MaxInBoundPerIp(n uint) ConnCtrlOption {
	self.MaxConnInBoundPerIP = n
	return self
}

func (self ConnCtrlOption) ReservedOnly(peers []string) ConnCtrlOption {
	self.ReservedPeers = NewStaticReserveFilter(peers)
	return self
}

func (self ConnCtrlOption) WithDialer(dialer Dialer) ConnCtrlOption {
	self.dialer = dialer
	return self
}

func ConnCtrlOptionFromConfig(config *config.P2PNodeConfig, reserveFilter p2p.AddressFilter) (option ConnCtrlOption, err error) {
	dialer, e := NewDialer(config)
	if e != nil {
		err = e
		return
	}
	return ConnCtrlOption{
		MaxConnOutBound:     config.MaxConnOutBound,
		MaxConnInBound:      config.MaxConnInBound,
		MaxConnInBoundPerIP: config.MaxConnInBoundForSingleIP,
		ReservedPeers:       reserveFilter,

		dialer: dialer,
	}, nil
}
