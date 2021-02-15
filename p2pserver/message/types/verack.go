
package types

import (
	"io"

	comm " github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/p2pserver/common"
)

type VerACK struct {
	//TODO remove this legecy field when upgrade network layer protocal
	isConsensus bool
}

//Serialize message payload
func (this *VerACK) Serialization(sink *comm.ZeroCopySink) {
	sink.WriteBool(this.isConsensus)
}

func (this *VerACK) CmdType() string {
	return common.VERACK_TYPE
}

//Deserialize message payload
func (this *VerACK) Deserialization(source *comm.ZeroCopySource) error {
	var irregular, eof bool
	this.isConsensus, irregular, eof = source.NextBool()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return comm.ErrIrregularData
	}

	return nil
}
