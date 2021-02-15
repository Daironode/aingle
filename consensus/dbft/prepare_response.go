 
package dbft

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type PrepareResponse struct {
	msgData   ConsensusMessageData
	Signature []byte
}

func (pres *PrepareResponse) Serialization(sink *common.ZeroCopySink) {
	pres.msgData.Serialization(sink)
	sink.WriteVarBytes(pres.Signature)
}

//read data to reader
func (pres *PrepareResponse) Deserialization(source *common.ZeroCopySource) error {
	err := pres.msgData.Deserialization(source)
	if err != nil {
		return err
	}

	sign, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	pres.Signature = sign

	return nil
}

func (pres *PrepareResponse) Type() ConsensusMessageType {
	return pres.ConsensusMessageData().Type
}

func (pres *PrepareResponse) ViewNumber() byte {
	return pres.msgData.ViewNumber
}

func (pres *PrepareResponse) ConsensusMessageData() *ConsensusMessageData {
	return &(pres.msgData)
}
