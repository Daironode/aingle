
package req

import (
	"sync/atomic"
	"unsafe"

	"github.com/Daironode/aingle-event/actor"
)

var consensusPid unsafe.Pointer

func SetConsensusPid(conPid *actor.PID) {
	atomic.StorePointer(&consensusPid, unsafe.Pointer(conPid))
}

func GetConsensusPid() *actor.PID {
	return (*actor.PID)(atomic.LoadPointer(&consensusPid))
}
