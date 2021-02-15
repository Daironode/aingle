 package connect_controller

import (
	"net"
	"sort"
)

type StaticReserveFilter struct {
	//format: host or ip
	ReservedPeers []string
}

func NewStaticReserveFilter(peers []string) *StaticReserveFilter {
	// put domain to the end
	sort.Slice(peers, func(i, j int) bool {
		return net.ParseIP(peers[i]) != nil
	})
	return &StaticReserveFilter{
		ReservedPeers: peers,
	}
}

// remoteAddr format 192.168.1.1:61234
// if reserved peers is empty, we should handle this case in subnet now
// since for gov node, reserve_result = in_subnet_set || in_static_set
// for normal node, reserve_result = in_static_set || static_set_is_empty
// because the information of whether self node is gov or not is in subnet module
func (self *StaticReserveFilter) Contains(remoteIPPort string) bool {
	// 192.168.1.1 in reserve list, 192.168.1.111:61234 and 192.168.1.11:61234 can connect in if we are using prefix matching
	// so get this IP to do fully match
	remoteAddr, _, err := net.SplitHostPort(remoteIPPort)
	if err != nil {
		return false
	}
	// we don't load domain in start because we consider domain's A/AAAA record may change sometimes
	for _, curIPOrName := range self.ReservedPeers {
		curIPs, err := net.LookupHost(curIPOrName)
		if err != nil {
			continue
		}
		for _, digIP := range curIPs {
			if digIP == remoteAddr {
				return true
			}
		}
	}

	return false
}
