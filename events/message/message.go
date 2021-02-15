
package message

import (
	" github.com/Daironode/aingle/core/types"
)

const (
	TOPIC_SAVE_BLOCK_COMPLETE = "svblkcmp"
	TOPIC_SMART_CODE_EVENT    = "scevt"
)

type SaveBlockCompleteMsg struct {
	Block *types.Block
}

type SmartCodeEventMsg struct {
	Event *types.SmartCodeEvent
}

type BlockConsensusComplete struct {
	Block *types.Block
}
