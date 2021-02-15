
package types

import (
	"io"

	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type NotFound struct {
	Hash common.Uint256
}

//Serialize message payload
func (this NotFound) Serialization(sink *common.ZeroCopySink) {
	sink.WriteHash(this.Hash)
}

func (this NotFound) CmdType() string {
	return comm.NOT_FOUND_TYPE
}

//Deserialize message payload
func (this *NotFound) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	this.Hash, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}
