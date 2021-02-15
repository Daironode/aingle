
package nodeinfo

import (
	"time"

	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/core/ledger"
	" github.com/Daironode/aingle/p2pserver/net/netserver"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
	" github.com/Daironode/aingle/p2pserver/protocols"
	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	nodePortMetric = prom.NewGauge(prom.GaugeOpts{
		Name: "aingle_nodeport",
		Help: "aingle node port",
	})

	blockHeightMetric = prom.NewGauge(prom.GaugeOpts{
		Name: "aingle_block_height",
		Help: "aingle blockchain block height",
	})

	inboundsCountMetric = prom.NewGauge(prom.GaugeOpts{
		Name: "aingle_p2p_inbounds_count",
		Help: "aingle p2p inbloud count",
	})

	outboundsCountMetric = prom.NewGauge(prom.GaugeOpts{
		Name: "aingle_p2p_outbounds_count",
		Help: "aingle p2p outbloud count",
	})

	peerStatusMetric = prom.NewGaugeVec(prom.GaugeOpts{
		Name: "aingle_p2p_peer_status",
		Help: "aingle peer info",
	}, []string{"ip", "id"})

	reconnectCountMetric = prom.NewGauge(prom.GaugeOpts{
		Name: "aingle_p2p_reconnect_count",
		Help: "aingle p2p reconnect count",
	})
)

var (
	metrics = []prom.Collector{nodePortMetric, blockHeightMetric, inboundsCountMetric,
		outboundsCountMetric, peerStatusMetric, reconnectCountMetric}
)

func initMetric() error {
	for _, curMetric := range metrics {
		if err := prom.Register(curMetric); err != nil {
			return err
		}
	}

	return nil
}

func metricUpdate(n p2p.P2P) {
	nodePortMetric.Set(float64(config.DefConfig.P2PNode.NodePort))

	blockHeightMetric.Set(float64(ledger.DefLedger.GetCurrentBlockHeight()))

	ns, ok := n.(*netserver.NetServer)
	if !ok {
		return
	}

	inboundsCountMetric.Set(float64(ns.ConnectController().InboundsCount()))
	outboundsCountMetric.Set(float64(ns.ConnectController().OutboundsCount()))

	peers := ns.GetNeighbors()
	for _, curPeer := range peers {
		id := curPeer.GetID()

		// label: IP PeedID
		peerStatusMetric.WithLabelValues(curPeer.GetAddr(), id.ToHexString()).Set(float64(curPeer.GetHeight()))
	}

	pt := ns.Protocol()
	mh, ok := pt.(*protocols.MsgHandler)
	if !ok {
		return
	}

	reconnectCountMetric.Set(float64(mh.ReconnectService().ReconnectCount()))
}

func updateMetric(n p2p.P2P) {
	tk := time.NewTicker(time.Minute)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			metricUpdate(n)
		}
	}
}
