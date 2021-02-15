
package types

import (
	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type AddrReq struct{}

//Serialize message payload
func (this AddrReq) Serialization(sink *common.ZeroCopySink) {
}

func (this *AddrReq) CmdType() string {
	return comm.GetADDR_TYPE
}

//Deserialize message payload
func (this *AddrReq) Deserialization(source *common.ZeroCopySource) error {
	return nil
}
