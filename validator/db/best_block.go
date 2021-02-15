
package db

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
)

type BestBlock struct {
	Height uint32
	Hash   common.Uint256
}

type BestStateProvider interface {
	GetBestBlock() (BestBlock, error)
	GetBestHeader() (*types.Header, error)
}
