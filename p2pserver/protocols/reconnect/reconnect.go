
package reconnect

import (
	"math/rand"
	"sync"
	"time"

	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/p2pserver/common"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/p2pserver/peer"
)

type ReconnectPeerInfo struct {
	count int // current retry count
	id    common.PeerId
}

const MaxRetryCountForReserveNode = 100

//ReconnectService contain addr need to reconnect
type ReconnectService struct {
	sync.RWMutex
	MaxRetryCount       int
	RetryAddrs          map[string]*ReconnectPeerInfo
	net                 p2p.P2P
	staticReserveFilter p2p.AddressFilter
	quit                chan bool
}

func NewReconectService(net p2p.P2P, staticReserveFilter p2p.AddressFilter) *ReconnectService {
	return &ReconnectService{
		net:                 net,
		staticReserveFilter: staticReserveFilter,
		MaxRetryCount:       common.MAX_RETRY_COUNT,
		quit:                make(chan bool),
		RetryAddrs:          make(map[string]*ReconnectPeerInfo),
	}
}

func (self *ReconnectService) Start() {
	go self.keepOnlineService()
}

func (self *ReconnectService) Stop() {
	close(self.quit)
}

func (this *ReconnectService) keepOnlineService() {
	tick := time.NewTicker(time.Second * common.CONN_MONITOR)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			this.retryInactivePeer()
		case <-this.quit:
			return
		}
	}
}

func (self *ReconnectService) OnAddPeer(p *peer.PeerInfo) {
	listenAddr := p.RemoteListenAddress()
	self.Lock()
	delete(self.RetryAddrs, listenAddr)
	self.Unlock()
}

func (self *ReconnectService) OnDelPeer(p *peer.PeerInfo) {
	nodeAddr := p.RemoteListenAddress()
	maxCount := self.MaxRetryCount
	if self.staticReserveFilter.Contains(nodeAddr) {
		maxCount = MaxRetryCountForReserveNode
	}
	self.Lock()
	self.RetryAddrs[nodeAddr] = &ReconnectPeerInfo{count: maxCount, id: p.Id}
	self.Unlock()
}

func (this *ReconnectService) retryInactivePeer() {
	net := this.net
	connCount := net.GetOutConnRecordLen()
	if connCount >= config.DefConfig.P2PNode.MaxConnOutBound {
		log.Warnf("[p2p]Connect: out connections(%d) reach max limit(%d)", connCount,
			config.DefConfig.P2PNode.MaxConnOutBound)
		return
	}

	//try connect
	var addrs []string
	this.Lock()
	if len(this.RetryAddrs) > 0 {
		list := make(map[string]*ReconnectPeerInfo)
		for addr, v := range this.RetryAddrs {
			v.count -= 1
			if v.count >= 0 && net.GetPeer(v.id) == nil {
				addrs = append(addrs, addr)
				list[addr] = v
			}
		}

		this.RetryAddrs = list
	}
	this.Unlock()
	for _, addr := range addrs {
		rand.Seed(time.Now().UnixNano())
		log.Debug("[p2p]Try to reconnect peer, peer addr is ", addr)
		<-time.After(time.Duration(rand.Intn(common.CONN_MAX_BACK)) * time.Millisecond)
		log.Debug("[p2p]Back off time`s up, start connect node")
		net.Connect(addr)
	}
}

func (self *ReconnectService) ReconnectCount() int {
	self.RLock()
	defer self.RUnlock()
	return len(self.RetryAddrs)
}
