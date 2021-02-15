
package actor

import " github.com/Daironode/aingle/core/types"

type StartConsensus struct{}
type StopConsensus struct{}

//internal Message
type TimeOut struct{}
type BlockCompleted struct {
	Block *types.Block
}
