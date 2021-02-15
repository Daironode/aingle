
package p2p

type AddressFilter interface {
	// addr format : ip:port
	Contains(addr string) bool
}

func CombineAddrFilter(filter1, filter2 AddressFilter) AddressFilter {
	return &combineAddrFilter{filter1: filter1, filter2: filter2}
}

func NoneAddrFilter() AddressFilter {
	return &noneAddrFilter{}
}

type combineAddrFilter struct {
	filter1 AddressFilter
	filter2 AddressFilter
}

func (self *combineAddrFilter) Contains(addr string) bool {
	return self.filter1.Contains(addr) || self.filter2.Contains(addr)
}

type noneAddrFilter struct{}

func (self *noneAddrFilter) Contains(addr string) bool {
	return false
}

func AllAddrFilter() AddressFilter {
	return &allAddrFilter{}
}

type allAddrFilter struct{}

func (self *allAddrFilter) Contains(addr string) bool {
	return true
}
