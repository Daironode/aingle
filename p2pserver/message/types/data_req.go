
package types

import (
	"io"

	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type DataReq struct {
	DataType common.InventoryType
	Hash     common.Uint256
}

//Serialize message payload
func (this DataReq) Serialization(sink *common.ZeroCopySink) {
	sink.WriteByte(byte(this.DataType))
	sink.WriteHash(this.Hash)
}

func (this *DataReq) CmdType() string {
	return comm.GET_DATA_TYPE
}

//Deserialize message payload
func (this *DataReq) Deserialization(source *common.ZeroCopySource) error {
	ty, eof := source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.DataType = common.InventoryType(ty)

	this.Hash, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}
