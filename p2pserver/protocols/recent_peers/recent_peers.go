
package recent_peers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	common2 " github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/p2pserver/common"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
)

type PersistRecentPeerService struct {
	net         p2p.P2P
	quit        chan bool
	recentPeers map[uint32][]*RecentPeer
	lock        sync.RWMutex
}

func (this *PersistRecentPeerService) contains(addr string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	for i := 0; i < len(this.recentPeers[netID]); i++ {
		if this.recentPeers[netID][i].Addr == addr {
			return true
		}
	}
	return false
}

func (this *PersistRecentPeerService) AddNodeAddr(addr string) {
	if this.contains(addr) {
		return
	}
	this.lock.Lock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	this.recentPeers[netID] = append(this.recentPeers[netID],
		&RecentPeer{
			Addr:  addr,
			Birth: time.Now().Unix(),
		})
	this.lock.Unlock()
}

func (this *PersistRecentPeerService) DelNodeAddr(addr string) {
	this.lock.Lock()
	netID := config.DefConfig.P2PNode.NetworkMagic
	for i := 0; i < len(this.recentPeers[netID]); i++ {
		if this.recentPeers[netID][i].Addr == addr {
			this.recentPeers[netID] = append(this.recentPeers[netID][:i], this.recentPeers[netID][i+1:]...)
		}
	}
	this.lock.Unlock()
}

type RecentPeer struct {
	Addr  string
	Birth int64
}

func (this *PersistRecentPeerService) saveToFile() {
	temp := make(map[uint32][]string)
	this.lock.RLock()
	for networkId, rps := range this.recentPeers {
		temp[networkId] = make([]string, 0)
		for _, rp := range rps {
			elapse := time.Now().Unix() - rp.Birth
			if elapse > common.RecentPeerElapseLimit {
				temp[networkId] = append(temp[networkId], rp.Addr)
			}
		}
	}
	this.lock.RUnlock()
	buf, err := json.Marshal(temp)
	if err != nil {
		log.Warn("[p2p]package recent peer fail: ", err)
		return
	}
	err = ioutil.WriteFile(common.RECENT_FILE_NAME, buf, os.ModePerm)
	if err != nil {
		log.Warn("[p2p]write recent peer fail: ", err)
	}
}

func NewPersistRecentPeerService(net p2p.P2P) *PersistRecentPeerService {
	return &PersistRecentPeerService{
		net:         net,
		recentPeers: make(map[uint32][]*RecentPeer),
		quit:        make(chan bool),
	}
}

func (self *PersistRecentPeerService) Stop() {
	close(self.quit)
}

func (this *PersistRecentPeerService) loadRecentPeers() {
	if common2.FileExisted(common.RECENT_FILE_NAME) {
		buf, err := ioutil.ReadFile(common.RECENT_FILE_NAME)
		if err != nil {
			log.Warn("[p2p]read %s fail:%s, connect recent peers cancel", common.RECENT_FILE_NAME, err.Error())
			return
		}

		temp := make(map[uint32][]string)
		err = json.Unmarshal(buf, &temp)
		if err != nil {
			log.Warn("[p2p]parse recent peer file fail: ", err)
			return
		}
		this.lock.Lock()
		defer this.lock.Unlock()
		for networkId, addrs := range temp {
			for _, addr := range addrs {
				this.recentPeers[networkId] = append(this.recentPeers[networkId], &RecentPeer{
					Addr:  addr,
					Birth: time.Now().Unix(),
				})
			}
		}
	}
}

func (this *PersistRecentPeerService) Start() {
	this.loadRecentPeers()
	this.tryRecentPeers()
	go this.syncUpRecentPeers()
}

//tryRecentPeers try connect recent contact peer when service start
func (this *PersistRecentPeerService) tryRecentPeers() {
	netID := config.DefConfig.P2PNode.NetworkMagic
	this.lock.RLock()
	defer this.lock.RUnlock()
	if len(this.recentPeers[netID]) > 0 {
		log.Info("[p2p] try to connect recent peer")
	}
	for _, v := range this.recentPeers[netID] {
		go this.net.Connect(v.Addr)
	}
}

//syncUpRecentPeers sync up recent peers periodically
func (this *PersistRecentPeerService) syncUpRecentPeers() {
	periodTime := common.RECENT_TIMEOUT
	t := time.NewTicker(time.Second * (time.Duration(periodTime)))
	for {
		select {
		case <-t.C:
			this.saveToFile()
		case <-this.quit:
			t.Stop()
			return
		}
	}
}
