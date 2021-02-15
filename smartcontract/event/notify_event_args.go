
package event

import (
	" github.com/Daironode/aingle/common"
)

const (
	CONTRACT_STATE_FAIL    byte = 0
	CONTRACT_STATE_SUCCESS byte = 1
)

// NotifyEventInfo describe smart contract event notify info struct
type NotifyEventInfo struct {
	ContractAddress common.Address
	States          interface{}
}

type ExecuteNotify struct {
	TxHash      common.Uint256
	State       byte
	GasConsumed uint64
	Notify      []*NotifyEventInfo
}
