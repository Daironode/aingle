 
package utils

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const HostsResolverCacheTime = time.Minute * 10

// host resolver with cache
type HostsResolver struct {
	hosts [][2]string

	lock  sync.Mutex     // avoid concurrent cache refresh
	cache unsafe.Pointer // atomic pointer to HostsCache, avoid read&write data race
}

type HostsCache struct {
	refreshTime time.Time
	addrs       []string
}

func NewHostsResolver(hosts []string) (*HostsResolver, []string) {
	resolver := &HostsResolver{}
	var invalids []string
	for _, n := range hosts {
		host, port, e := net.SplitHostPort(n)
		if e != nil {
			invalids = append(invalids, n)
			continue
		}
		resolver.hosts = append(resolver.hosts, [2]string{host, port})
	}

	return resolver, invalids
}

func (self *HostsResolver) GetHostAddrs() []string {
	// fast path test
	cached := (*HostsCache)(atomic.LoadPointer(&self.cache))
	if cached != nil && cached.refreshTime.Add(HostsResolverCacheTime).After(time.Now()) && len(cached.addrs) != 0 {
		return cached.addrs
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	cached = (*HostsCache)(self.cache)
	if cached != nil && cached.refreshTime.Add(HostsResolverCacheTime).After(time.Now()) && len(cached.addrs) != 0 {
		return cached.addrs
	}

	cache := make([]string, 0, len(self.hosts))
	for _, n := range self.hosts {
		host, port := n[0], n[1]
		ns, err := net.LookupHost(host)
		if err != nil || len(ns) == 0 {
			continue
		}

		for _, hs := range ns {
			cache = append(cache, net.JoinHostPort(hs, port))
		}
	}

	atomic.StorePointer(&self.cache, unsafe.Pointer(&HostsCache{refreshTime: time.Now(), addrs: cache}))

	return cache
}
