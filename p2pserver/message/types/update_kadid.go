
package types

import (
	common2 " github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/p2pserver/common"
)

type UpdatePeerKeyId struct {
	//TODO remove this legecy field when upgrade network layer protocal
	KadKeyId *common.PeerKeyId
}

//Serialize message payload
func (this *UpdatePeerKeyId) Serialization(sink *common2.ZeroCopySink) {
	this.KadKeyId.Serialization(sink)
}

func (this *UpdatePeerKeyId) Deserialization(source *common2.ZeroCopySource) error {
	this.KadKeyId = &common.PeerKeyId{}
	return this.KadKeyId.Deserialization(source)
}

func (this *UpdatePeerKeyId) CmdType() string {
	return common.UPDATE_KADID_TYPE
}
