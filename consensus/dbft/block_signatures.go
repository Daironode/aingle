
package dbft

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type BlockSignatures struct {
	msgData    ConsensusMessageData
	Signatures []SignaturesData
}

type SignaturesData struct {
	Signature []byte
	Index     uint16
}

func (self *BlockSignatures) Serialization(sink *common.ZeroCopySink) {
	self.msgData.Serialization(sink)
	sink.WriteVarUint(uint64(len(self.Signatures)))

	for _, sign := range self.Signatures {
		sink.WriteVarBytes(sign.Signature)
		sink.WriteUint16(sign.Index)
	}
}

func (self *BlockSignatures) Deserialization(source *common.ZeroCopySource) error {
	err := self.msgData.Deserialization(source)
	if err != nil {
		return err
	}

	length, _, irregular, eof := source.NextVarUint()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	for i := uint64(0); i < length; i++ {
		sig := SignaturesData{}

		sig.Signature, _, irregular, eof = source.NextVarBytes()
		if eof {
			return io.ErrUnexpectedEOF
		}

		if irregular {
			return common.ErrIrregularData
		}

		sig.Index, eof = source.NextUint16()
		if eof {
			return io.ErrUnexpectedEOF
		}

		self.Signatures = append(self.Signatures, sig)
	}

	return nil
}

func (self *BlockSignatures) Type() ConsensusMessageType {
	return self.ConsensusMessageData().Type
}

func (self *BlockSignatures) ViewNumber() byte {
	return self.msgData.ViewNumber
}

func (self *BlockSignatures) ConsensusMessageData() *ConsensusMessageData {
	return &(self.msgData)
}
