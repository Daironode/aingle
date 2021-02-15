
package types

import (
	"io"

	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type HeadersReq struct {
	Len       uint8
	HashStart common.Uint256
	HashEnd   common.Uint256
}

//Serialize message payload
func (this *HeadersReq) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint8(this.Len)
	sink.WriteHash(this.HashStart)
	sink.WriteHash(this.HashEnd)
}

func (this *HeadersReq) CmdType() string {
	return comm.GET_HEADERS_TYPE
}

//Deserialize message payload
func (this *HeadersReq) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	this.Len, eof = source.NextUint8()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.HashStart, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.HashEnd, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}
