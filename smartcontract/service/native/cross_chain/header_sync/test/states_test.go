
package test

import (
	"testing"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/header_sync"
	"github.com/stretchr/testify/assert"
)

func TestPeer(t *testing.T) {
	peer := header_sync.Peer{
		Index:      1,
		PeerPubkey: "testPubkey",
	}
	sink := common.NewZeroCopySink(nil)
	peer.Serialization(sink)

	var p header_sync.Peer
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, peer, p)
}

func TestKeyHeights(t *testing.T) {
	key := header_sync.KeyHeights{
		HeightList: []uint32{1, 2, 3, 4},
	}
	sink := common.NewZeroCopySink(nil)
	key.Serialization(sink)

	var k header_sync.KeyHeights
	err := k.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, key, k)
}

func TestConsensusPeers(t *testing.T) {
	peers := header_sync.ConsensusPeers{
		ChainID: 1,
		Height:  2,
		PeerMap: map[string]*header_sync.Peer{
			"testPubkey1": {
				Index:      1,
				PeerPubkey: "testPubkey1",
			},
			"testPubkey2": {
				Index:      2,
				PeerPubkey: "testPubkey2",
			},
		},
	}
	sink := common.NewZeroCopySink(nil)
	peers.Serialization(sink)

	var p header_sync.ConsensusPeers
	err := p.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.NoError(t, err)
	assert.Equal(t, p, peers)
}
