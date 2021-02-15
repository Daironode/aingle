
package kbucket

import (
	"bytes"
	"container/list"
	"sort"

	" github.com/Daironode/aingle/p2pserver/common"
)

// A helper struct to sort peers by their distance to the local node
type peerDistance struct {
	p        common.PeerIDAddressPair
	distance [20]byte
}

// peerDistanceSorter implements sort.Interface to sort peers by xor distance
type peerDistanceSorter struct {
	peers  []peerDistance
	target common.PeerId
}

func (pds *peerDistanceSorter) Len() int { return len(pds.peers) }
func (pds *peerDistanceSorter) Swap(a, b int) {
	pds.peers[a], pds.peers[b] = pds.peers[b], pds.peers[a]
}
func (pds *peerDistanceSorter) Less(a, b int) bool {
	return bytes.Compare(pds.peers[a].distance[:], pds.peers[b].distance[:]) < 0
}

// Append the peer.ID to the sorter's slice. It may no longer be sorted.
func (pds *peerDistanceSorter) appendPeer(p common.PeerIDAddressPair) {
	pds.peers = append(pds.peers, peerDistance{
		p:        p,
		distance: pds.target.Distance(p.ID),
	})
}

// Append the peer.ID values in the list to the sorter's slice. It may no longer be sorted.
func (pds *peerDistanceSorter) appendPeersFromList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		pds.appendPeer(e.Value.(common.PeerIDAddressPair))
	}
}

func (pds *peerDistanceSorter) sort() {
	sort.Sort(pds)
}
