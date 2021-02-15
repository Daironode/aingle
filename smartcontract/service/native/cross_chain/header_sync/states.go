
package header_sync

import (
	"fmt"
	"math"
	"sort"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type Peer struct {
	Index      uint32
	PeerPubkey string
}

func (this *Peer) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, uint64(this.Index))
	utils.EncodeString(sink, this.PeerPubkey)
}

func (this *Peer) Deserialization(source *common.ZeroCopySource) error {
	index, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize index error: %v", err)
	}
	if index > math.MaxUint32 {
		return fmt.Errorf("deserialize index error: index more than max uint32")
	}
	peerPubkey, err := utils.DecodeString(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeString, deserialize peerPubkey error: %v", err)
	}
	this.Index = uint32(index)
	this.PeerPubkey = peerPubkey
	return nil
}

type KeyHeights struct {
	HeightList []uint32
}

func (this *KeyHeights) Serialization(sink *common.ZeroCopySink) {
	//first sort the list  (big -> small)
	sort.SliceStable(this.HeightList, func(i, j int) bool {
		return this.HeightList[i] > this.HeightList[j]
	})
	utils.EncodeVarUint(sink, uint64(len(this.HeightList)))
	for _, v := range this.HeightList {
		utils.EncodeVarUint(sink, uint64(v))
	}
}

func (this *KeyHeights) Deserialization(source *common.ZeroCopySource) error {
	n, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize HeightList length error: %v", err)
	}
	heightList := make([]uint32, 0)
	for i := 0; uint64(i) < n; i++ {
		height, err := utils.DecodeVarUint(source)
		if err != nil {
			return fmt.Errorf("utils.DecodeVarUint, deserialize height error: %v", err)
		}
		if height > math.MaxUint32 {
			return fmt.Errorf("deserialize height error: height more than max uint32")
		}
		heightList = append(heightList, uint32(height))
	}
	this.HeightList = heightList
	return nil
}

type ConsensusPeers struct {
	ChainID uint64
	Height  uint32
	PeerMap map[string]*Peer
}

func (this *ConsensusPeers) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.ChainID)
	utils.EncodeVarUint(sink, uint64(this.Height))
	utils.EncodeVarUint(sink, uint64(len(this.PeerMap)))
	var peerList []*Peer
	for _, v := range this.PeerMap {
		peerList = append(peerList, v)
	}
	sort.SliceStable(peerList, func(i, j int) bool {
		return peerList[i].PeerPubkey > peerList[j].PeerPubkey
	})
	for _, v := range peerList {
		v.Serialization(sink)
	}
}

func (this *ConsensusPeers) Deserialization(source *common.ZeroCopySource) error {
	chainID, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize chainID error: %v", err)
	}
	height, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize height error: %v", err)
	}
	if height > math.MaxUint32 {
		return fmt.Errorf("deserialize height error: height more than max uint32")
	}
	n, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize HeightList length error: %v", err)
	}
	peerMap := make(map[string]*Peer)
	for i := 0; uint64(i) < n; i++ {
		peer := new(Peer)
		if err := peer.Deserialization(source); err != nil {
			return fmt.Errorf("deserialize peer error: %v", err)
		}
		peerMap[peer.PeerPubkey] = peer
	}
	this.ChainID = chainID
	this.Height = uint32(height)
	this.PeerMap = peerMap
	return nil
}
