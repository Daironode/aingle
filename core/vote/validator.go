
package vote

import (
	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/core/genesis"
	" github.com/Daironode/aingle/core/types"
)

func GetValidators(txs []*types.Transaction) ([]keypair.PublicKey, error) {
	// TODO implement vote
	return genesis.GenesisBookkeepers, nil
}
