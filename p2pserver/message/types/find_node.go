
package types

import (
	"io"

	" github.com/Daironode/aingle/common"
	ncomm " github.com/Daironode/aingle/p2pserver/common"
)

type FindNodeReq struct {
	TargetID ncomm.PeerId
}

// Serialization message payload
func (req FindNodeReq) Serialization(sink *common.ZeroCopySink) {
	req.TargetID.Serialization(sink)
}

// CmdType return this message type
func (req *FindNodeReq) CmdType() string {
	return ncomm.FINDNODE_TYPE
}

// Deserialization message payload
func (req *FindNodeReq) Deserialization(source *common.ZeroCopySource) error {
	return req.TargetID.Deserialization(source)
}

type FindNodeResp struct {
	TargetID    ncomm.PeerId
	Success     bool
	Address     string
	CloserPeers []ncomm.PeerIDAddressPair
}

// Serialization message payload
func (resp FindNodeResp) Serialization(sink *common.ZeroCopySink) {
	resp.TargetID.Serialization(sink)
	sink.WriteBool(resp.Success)
	sink.WriteString(resp.Address)
	sink.WriteUint32(uint32(len(resp.CloserPeers)))
	for _, curPeer := range resp.CloserPeers {
		curPeer.ID.Serialization(sink)
		sink.WriteString(curPeer.Address)
	}
}

// CmdType return this message type
func (resp *FindNodeResp) CmdType() string {
	return ncomm.FINDNODE_RESP_TYPE
}

// Deserialization message payload
func (resp *FindNodeResp) Deserialization(source *common.ZeroCopySource) error {
	err := resp.TargetID.Deserialization(source)
	if err != nil {
		return err
	}

	succ, _, eof := source.NextBool()
	if eof {
		return io.ErrUnexpectedEOF
	}
	resp.Success = succ

	addr, _, _, eof := source.NextString()
	if eof {
		return io.ErrUnexpectedEOF
	}
	resp.Address = addr

	numCloser, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}

	for i := 0; i < int(numCloser); i++ {
		var curpa ncomm.PeerIDAddressPair
		id := ncomm.PeerId{}
		err = id.Deserialization(source)
		if err != nil {
			return err
		}
		curpa.ID = id
		addr, _, _, eof := source.NextString()
		if eof {
			return io.ErrUnexpectedEOF
		}
		curpa.Address = addr

		resp.CloserPeers = append(resp.CloserPeers, curpa)
	}

	return nil
}
