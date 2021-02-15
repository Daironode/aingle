 
package types

import (
	"io"

	" github.com/Daironode/aingle/common"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type Addr struct {
	NodeAddrs []comm.PeerAddr
}

//Serialize message payload
func (this Addr) Serialization(sink *common.ZeroCopySink) {
	num := uint64(len(this.NodeAddrs))
	sink.WriteUint64(num)

	for _, addr := range this.NodeAddrs {
		sink.WriteInt64(addr.Time)
		sink.WriteUint64(addr.Services)
		sink.WriteBytes(addr.IpAddr[:])
		sink.WriteUint16(addr.Port)
		sink.WriteUint16(addr.ConsensusPort)
		sink.WriteUint64(addr.ID.ToUint64())
	}
}

func (this *Addr) CmdType() string {
	return comm.ADDR_TYPE
}

func (this *Addr) Deserialization(source *common.ZeroCopySource) error {
	count, eof := source.NextUint64()
	if eof {
		return io.ErrUnexpectedEOF
	}

	for i := 0; i < int(count); i++ {
		var addr comm.PeerAddr
		addr.Time, eof = source.NextInt64()
		if eof {
			return io.ErrUnexpectedEOF
		}
		addr.Services, eof = source.NextUint64()
		if eof {
			return io.ErrUnexpectedEOF
		}
		buf, _ := source.NextBytes(uint64(len(addr.IpAddr[:])))
		copy(addr.IpAddr[:], buf)
		addr.Port, eof = source.NextUint16()
		if eof {
			return io.ErrUnexpectedEOF
		}
		addr.ConsensusPort, eof = source.NextUint16()
		if eof {
			return io.ErrUnexpectedEOF
		}
		id, eof := source.NextUint64()
		if eof {
			return io.ErrUnexpectedEOF
		}
		addr.ID = comm.PseudoPeerIdFromUint64(id)

		this.NodeAddrs = append(this.NodeAddrs, addr)
	}

	if count > comm.MAX_ADDR_NODE_CNT {
		count = comm.MAX_ADDR_NODE_CNT
	}
	this.NodeAddrs = this.NodeAddrs[:count]

	return nil
}
