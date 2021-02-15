 package bootstrap

import (
	"math/rand"
	"sync/atomic"
	"time"

	" github.com/Daironode/aingle/p2pserver/common"
	msgpack " github.com/Daironode/aingle/p2pserver/message/msg_pack"
	" github.com/Daironode/aingle/p2pserver/message/types"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/p2pserver/peer"
	" github.com/Daironode/aingle/p2pserver/protocols/utils"
)

const activeConnect = 4 // when connection num less than this value, we connect seeds node actively.

type BootstrapService struct {
	seeds     *utils.HostsResolver
	connected uint32
	net       p2p.P2P
	quit      chan bool
}

func NewBootstrapService(net p2p.P2P, seeds *utils.HostsResolver) *BootstrapService {
	return &BootstrapService{
		seeds: seeds,
		net:   net,
		quit:  make(chan bool),
	}
}

func (self *BootstrapService) Start() {
	go self.connectSeedService()
}

func (self *BootstrapService) Stop() {
	close(self.quit)
}

func (self *BootstrapService) OnAddPeer(info *peer.PeerInfo) {
	atomic.AddUint32(&self.connected, 1)
}

func (self *BootstrapService) OnDelPeer(info *peer.PeerInfo) {
	atomic.AddUint32(&self.connected, ^uint32(0))
}

//connectSeedService make sure seed peer be connected
func (self *BootstrapService) connectSeedService() {
	t := time.NewTimer(0) // let it timeout to start connect immediately
	for {
		select {
		case <-t.C:
			self.connectSeeds()
			t.Stop()
			connected := atomic.LoadUint32(&self.connected)
			if connected >= activeConnect {
				t.Reset(time.Second * time.Duration(10*common.CONN_MONITOR))
			} else {
				t.Reset(time.Second * common.CONN_MONITOR)
			}
		case <-self.quit:
			t.Stop()
			return
		}
	}
}

//connectSeeds connect the seeds in seedlist and call for nbr list
func (self *BootstrapService) connectSeeds() {
	connPeers := make(map[string]*peer.Peer)
	nps := self.net.GetNeighbors()
	for _, tn := range nps {
		listenAddr := tn.Info.RemoteListenAddress()
		connPeers[listenAddr] = tn
	}

	seedConnList := make([]*peer.Peer, 0)
	seedDisconn := make([]string, 0)
	isSeed := false
	for _, nodeAddr := range self.seeds.GetHostAddrs() {
		if p, ok := connPeers[nodeAddr]; ok {
			seedConnList = append(seedConnList, p)
		} else {
			seedDisconn = append(seedDisconn, nodeAddr)
		}

		if self.net.IsOwnAddress(nodeAddr) {
			isSeed = true
		}
	}

	if len(seedConnList) > 0 {
		rand.Seed(time.Now().UnixNano())
		// close NewAddrReq
		index := rand.Intn(len(seedConnList))
		self.reqNbrList(seedConnList[index])
		if isSeed && len(seedDisconn) > 0 {
			index := rand.Intn(len(seedDisconn))
			go self.net.Connect(seedDisconn[index])
		}
	} else { //not found
		for _, nodeAddr := range self.seeds.GetHostAddrs() {
			go self.net.Connect(nodeAddr)
		}
	}
}

func (this *BootstrapService) reqNbrList(p *peer.Peer) {
	id := p.GetID()
	var msg types.Message
	if id.IsPseudoPeerId() {
		msg = msgpack.NewAddrReq()
	} else {
		msg = msgpack.NewFindNodeReq(this.net.GetID())
	}

	go this.net.SendTo(id, msg)
}
