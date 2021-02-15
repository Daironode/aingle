
package subnet

import "net"

type SubNetReservedAddrFilter struct {
	staticFilterEnabled bool
	subnet              *SubNet
}

func (self *SubNetReservedAddrFilter) Contains(addr string) bool {
	// seed node should allow all node connection
	if self.subnet.IsSeedNode() {
		return true
	}

	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	// gov node
	if self.subnet.acct != nil && self.subnet.gov.IsGovNodePubKey(self.subnet.acct.PublicKey) {
		return self.subnet.isSeedIp(ip) || self.subnet.IpInMembers(ip)
	}

	// normal node, if static filter is disabled, then allow all node connection
	return !self.staticFilterEnabled
}

type SubNetMaskAddrFilter struct {
	subnet *SubNet
}

func (self *SubNetMaskAddrFilter) Contains(addr string) bool {
	self.subnet.lock.Lock()
	defer self.subnet.lock.Unlock()
	_, ok := self.subnet.members[addr]

	return ok
}
