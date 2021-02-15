
package types

import (
	"io"

	comm " github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/p2pserver/common"
)

type BlocksReq struct {
	HeaderHashCount uint8
	HashStart       comm.Uint256
	HashStop        comm.Uint256
}

//Serialize message payload
func (this *BlocksReq) Serialization(sink *comm.ZeroCopySink) {
	sink.WriteUint8(this.HeaderHashCount)
	sink.WriteHash(this.HashStart)
	sink.WriteHash(this.HashStop)
}

func (this *BlocksReq) CmdType() string {
	return common.GET_BLOCKS_TYPE
}

//Deserialize message payload
func (this *BlocksReq) Deserialization(source *comm.ZeroCopySource) error {
	var eof bool
	this.HeaderHashCount, eof = source.NextUint8()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.HashStart, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.HashStop, eof = source.NextHash()

	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}
