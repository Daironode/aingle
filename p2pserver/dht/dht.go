
package dht

import (
	"time"

	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/p2pserver/common"
	kb " github.com/Daironode/aingle/p2pserver/dht/kbucket"
)

// Pool size is the number of nodes used for group find/set RPC calls
var PoolSize = 6

// K is the maximum number of requests to perform before returning failure.
var KValue = 20

// Alpha is the concurrency factor for asynchronous requests.
var AlphaValue = 3

type DHT struct {
	localId    common.PeerId
	bucketSize int
	routeTable *kb.RouteTable // Array of routing tables for differently distanced nodes

	AutoRefresh     bool
	RtRefreshPeriod time.Duration
}

// RouteTable return dht's routeTable
func (dht *DHT) RouteTable() *kb.RouteTable {
	return dht.routeTable
}

// NewDHT creates a new DHT with the specified host and options.
func NewDHT(id common.PeerId) *DHT {
	bucketSize := KValue
	rt := kb.NewRoutingTable(bucketSize, id)

	rt.PeerAdded = func(p common.PeerId) {
		log.Debugf("dht: peer: %s added to dht", p.ToHexString())
	}

	rt.PeerRemoved = func(p common.PeerId) {
		log.Debugf("dht: peer: %s removed from dht", p.ToHexString())
	}

	return &DHT{
		localId:         id,
		routeTable:      rt,
		bucketSize:      bucketSize,
		AutoRefresh:     true,
		RtRefreshPeriod: 10 * time.Second,
	}
}

// Update signals the routeTable to Update its last-seen status
// on the given peer.
func (dht *DHT) Update(peer common.PeerId, addr string) bool {
	err := dht.routeTable.Update(peer, addr)
	return err == nil
}

func (dht *DHT) Contains(peer common.PeerId) bool {
	_, ok := dht.routeTable.Find(peer)
	return ok
}

func (dht *DHT) Remove(peer common.PeerId) {
	dht.routeTable.Remove(peer)
}

func (dht *DHT) BetterPeers(id common.PeerId, count int) []common.PeerIDAddressPair {
	closer := dht.routeTable.NearestPeers(id, count)
	filtered := make([]common.PeerIDAddressPair, 0, len(closer))
	// don't include self and target id
	for _, curPair := range closer {
		if curPair.ID == dht.localId {
			continue
		}
		if curPair.ID == id {
			continue
		}
		filtered = append(filtered, curPair)
	}

	return filtered
}
