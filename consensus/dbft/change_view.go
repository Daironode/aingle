
package dbft

import (
	"io"

	" github.com/Daironode/aingle/common"
)

type ChangeView struct {
	msgData       ConsensusMessageData
	NewViewNumber byte
}

func (cv *ChangeView) Serialization(sink *common.ZeroCopySink) {
	cv.msgData.Serialization(sink)
	sink.WriteByte(cv.NewViewNumber)
}

//read data to reader
func (cv *ChangeView) Deserialization(source *common.ZeroCopySource) error {
	err := cv.msgData.Deserialization(source)
	if err != nil {
		return err
	}

	viewNum, eof := source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}
	cv.NewViewNumber = viewNum

	return nil
}

func (cv *ChangeView) Type() ConsensusMessageType {
	return cv.ConsensusMessageData().Type
}

func (cv *ChangeView) ViewNumber() byte {
	return cv.msgData.ViewNumber
}

func (cv *ChangeView) ConsensusMessageData() *ConsensusMessageData {
	return &(cv.msgData)
}
