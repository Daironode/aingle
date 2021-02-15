
package vconfig

import (
	"encoding/json"
	"fmt"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
)

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub keypair.PublicKey) string {
	return common.PubKeyToHex(pub)
}

func Pubkey(nodeid string) (keypair.PublicKey, error) {
	return common.PubKeyFromHex(nodeid)
}

func VbftBlock(header *types.Header) (*VbftBlockInfo, error) {
	blkInfo := &VbftBlockInfo{}
	if err := json.Unmarshal(header.ConsensusPayload, blkInfo); err != nil {
		return nil, fmt.Errorf("unmarshal blockInfo: %s", err)
	}
	return blkInfo, nil
}
