

package types

import (
	"testing"

	" github.com/Daironode/aingle/p2pserver/common"
)

func TestFindNodeRequest(t *testing.T) {
	var req FindNodeReq
	req.TargetID = common.PeerId{}

	MessageTest(t, &req)
}

func TestFindNodeResponse(t *testing.T) {
	var resp FindNodeResp
	resp.TargetID = common.PeerId{}
	resp.Address = "127.0.0.1:1222"
	id := common.PseudoPeerIdFromUint64(uint64(0x456))
	resp.CloserPeers = []common.PeerIDAddressPair{
		common.PeerIDAddressPair{
			ID:      id,
			Address: "127.0.0.1:4222",
		},
	}
	resp.Success = true

	MessageTest(t, &resp)
}
