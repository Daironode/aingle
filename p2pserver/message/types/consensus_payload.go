
package types

import (
	"fmt"
	"io"

	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/signature"
	" github.com/Daironode/aingle/errors"
	common2 " github.com/Daironode/aingle/p2pserver/common"
)

type ConsensusPayload struct {
	Version         uint32
	PrevHash        common.Uint256
	Height          uint32
	BookkeeperIndex uint16
	Timestamp       uint32
	Data            []byte
	Owner           keypair.PublicKey
	Signature       []byte
	PeerId          common2.PeerId
	hash            common.Uint256
}

//get the consensus payload hash
func (this *ConsensusPayload) Hash() common.Uint256 {
	return common.Uint256{}
}

//Check whether header is correct
func (this *ConsensusPayload) Verify() error {
	sink := common.NewZeroCopySink(nil)
	this.SerializationUnsigned(sink)

	err := signature.Verify(this.Owner, sink.Bytes(), this.Signature)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetVerifyFail, fmt.Sprintf("signature verify error. buf:%v", sink.Bytes()))
	}
	return nil
}

//serialize the consensus payload
func (this *ConsensusPayload) ToArray() []byte {
	return common.SerializeToBytes(this)
}

func (this *ConsensusPayload) GetMessage() []byte {
	//TODO: GetMessage
	//return sig.GetHashData(cp)
	return []byte{}
}

func (this *ConsensusPayload) Serialization(sink *common.ZeroCopySink) {
	this.SerializationUnsigned(sink)
	buf := keypair.SerializePublicKey(this.Owner)
	sink.WriteVarBytes(buf)
	sink.WriteVarBytes(this.Signature)
}

//Deserialize message payload
func (this *ConsensusPayload) Deserialization(source *common.ZeroCopySource) error {
	err := this.DeserializationUnsigned(source)
	if err != nil {
		return err
	}
	buf, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	this.Owner, err = keypair.DeserializePublicKey(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "deserialize publickey error")
	}

	this.Signature, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func (this *ConsensusPayload) SerializationUnsigned(sink *common.ZeroCopySink) {
	sink.WriteUint32(this.Version)
	sink.WriteHash(this.PrevHash)
	sink.WriteUint32(this.Height)
	sink.WriteUint16(this.BookkeeperIndex)
	sink.WriteUint32(this.Timestamp)
	sink.WriteVarBytes(this.Data)
}

func (this *ConsensusPayload) DeserializationUnsigned(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.Version, eof = source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.PrevHash, eof = source.NextHash()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Height, eof = source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.BookkeeperIndex, eof = source.NextUint16()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Timestamp, eof = source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Data, _, irregular, eof = source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	return nil
}
