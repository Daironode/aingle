
package types

import (
	"fmt"
	"io"

	" github.com/Daironode/aingle/common"
	ct " github.com/Daironode/aingle/core/types"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type BlkHeader struct {
	BlkHdr []*ct.Header
}

type RawBlockHeader struct {
	BlkHdr []*ct.RawHeader
}

func (this *RawBlockHeader) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint32(uint32(len(this.BlkHdr)))

	for _, header := range this.BlkHdr {
		header.Serialization(sink)
	}
}
func (this *RawBlockHeader) Deserialization(source *common.ZeroCopySource) error {
	panic("[block_header] unsupport")
}

func (this *RawBlockHeader) CmdType() string {
	return comm.HEADERS_TYPE
}

//Serialize message payload
func (this BlkHeader) Serialization(sink *common.ZeroCopySink) {
	sink.WriteUint32(uint32(len(this.BlkHdr)))

	for _, header := range this.BlkHdr {
		header.Serialization(sink)
	}
}

func (this *BlkHeader) CmdType() string {
	return comm.HEADERS_TYPE
}

//Deserialize message payload
func (this *BlkHeader) Deserialization(source *common.ZeroCopySource) error {
	var count uint32
	count, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}

	for i := 0; i < int(count); i++ {
		var headers ct.Header
		err := headers.Deserialization(source)
		if err != nil {
			return fmt.Errorf("deserialze BlkHeader error: %v", err)
		}
		this.BlkHdr = append(this.BlkHdr, &headers)
	}
	return nil
}
