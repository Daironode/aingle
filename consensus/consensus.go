
package consensus

import (
	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/account"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/consensus/dbft"
	" github.com/Daironode/aingle/consensus/solo"
	" github.com/Daironode/aingle/consensus/vbft"
	p2p " github.com/Daironode/aingle/p2pserver/net/protocol"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
	CONSENSUS_VBFT = "vbft"
)

func NewConsensusService(consensusType string, account *account.Account, txpool *actor.PID, ledger *actor.PID, p2p p2p.P2P) (ConsensusService, error) {
	if consensusType == "" {
		consensusType = CONSENSUS_DBFT
	}
	var consensus ConsensusService
	var err error
	switch consensusType {
	case CONSENSUS_DBFT:
		consensus, err = dbft.NewDbftService(account, txpool, p2p)
	case CONSENSUS_SOLO:
		consensus, err = solo.NewSoloService(account, txpool)
	case CONSENSUS_VBFT:
		consensus, err = vbft.NewVbftServer(account, txpool, p2p)
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
