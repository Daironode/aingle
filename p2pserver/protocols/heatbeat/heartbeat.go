 package heatbeat

import (
	"time"

	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/ledger"
	" github.com/Daironode/aingle/p2pserver/common"
	msgpack " github.com/Daironode/aingle/p2pserver/message/msg_pack"
	" github.com/Daironode/aingle/p2pserver/message/types"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
)

type HeartBeat struct {
	net    p2p.P2P
	id     common.PeerId
	quit   chan bool
	ledger *ledger.Ledger //ledger
}

func NewHeartBeat(net p2p.P2P, ld *ledger.Ledger) *HeartBeat {
	return &HeartBeat{
		id:     net.GetID(),
		net:    net,
		quit:   make(chan bool),
		ledger: ld,
	}
}

func (self *HeartBeat) Start() {
	go self.heartBeatService()
}

func (self *HeartBeat) Stop() {
	close(self.quit)
}
func (this *HeartBeat) heartBeatService() {
	var periodTime uint = config.DEFAULT_GEN_BLOCK_TIME / common.UPDATE_RATE_PER_BLOCK
	t := time.NewTicker(time.Second * (time.Duration(periodTime)))

	for {
		select {
		case <-t.C:
			this.ping()
			this.timeout()
		case <-this.quit:
			t.Stop()
			return
		}
	}
}

func (this *HeartBeat) ping() {
	height := this.ledger.GetCurrentBlockHeight()
	ping := msgpack.NewPingMsg(uint64(height))
	go this.net.Broadcast(ping)
}

//timeout trace whether some peer be long time no response
func (this *HeartBeat) timeout() {
	peers := this.net.GetNeighbors()
	var periodTime uint = config.DEFAULT_GEN_BLOCK_TIME / common.UPDATE_RATE_PER_BLOCK
	for _, p := range peers {
		t := p.GetContactTime()
		if t.Before(time.Now().Add(-1 * time.Second *
			time.Duration(periodTime) * common.KEEPALIVE_TIMEOUT)) {
			log.Warnf("[p2p]keep alive timeout!!!lost remote peer %d - %s from %s", p.GetID(), p.Link.GetAddr(), t.String())
			p.Close()
		}
	}
}

func (this *HeartBeat) PingHandle(ctx *p2p.Context, ping *types.Ping) {
	remotePeer := ctx.Sender()
	remotePeer.SetHeight(ping.Height)
	p2p := ctx.Network()

	height := ledger.DefLedger.GetCurrentBlockHeight()
	p2p.SetHeight(uint64(height))
	msg := msgpack.NewPongMsg(uint64(height))

	err := remotePeer.Send(msg)
	if err != nil {
		log.Warn(err)
	}
}

func (this *HeartBeat) PongHandle(ctx *p2p.Context, pong *types.Pong) {
	remotePeer := ctx.Network()
	remotePeer.SetHeight(pong.Height)
}
