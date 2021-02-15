
package dbft

import (
	"fmt"
	"io"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
	" github.com/Daironode/aingle/core/types"
)

type PrepareRequest struct {
	msgData        ConsensusMessageData
	Nonce          uint64
	NextBookkeeper common.Address
	Transactions   []*types.Transaction
	Signature      []byte
}

func (pr *PrepareRequest) Serialization(sink *common.ZeroCopySink) {
	pr.msgData.Serialization(sink)
	sink.WriteVarUint(pr.Nonce)
	sink.WriteAddress(pr.NextBookkeeper)
	sink.WriteVarUint(uint64(len(pr.Transactions)))
	for _, t := range pr.Transactions {
		t.Serialization(sink)
	}
	sink.WriteVarBytes(pr.Signature)
}

func (pr *PrepareRequest) Deserialization(source *common.ZeroCopySource) error {
	pr.msgData = ConsensusMessageData{}
	err := pr.msgData.Deserialization(source)
	if err != nil {
		return err
	}

	nonce, _, irregular, eof := source.NextVarUint()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	pr.Nonce = nonce
	pr.NextBookkeeper, eof = source.NextAddress()
	if eof {
		return io.ErrUnexpectedEOF
	}

	var length uint64
	length, _, irregular, eof = source.NextVarUint()
	if eof {
		return io.ErrUnexpectedEOF
	}

	if irregular {
		return common.ErrIrregularData
	}

	for i := 0; i < int(length); i++ {
		var t types.Transaction
		if err := t.Deserialization(source); err != nil {
			return fmt.Errorf("[PrepareRequest] transactions deserialization failed: %s", err)
		}
		pr.Transactions = append(pr.Transactions, &t)
	}

	pr.Signature, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}

	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func (pr *PrepareRequest) Type() ConsensusMessageType {
	log.Debug()
	return pr.ConsensusMessageData().Type
}

func (pr *PrepareRequest) ViewNumber() byte {
	log.Debug()
	return pr.msgData.ViewNumber
}

func (pr *PrepareRequest) ConsensusMessageData() *ConsensusMessageData {
	log.Debug()
	return &(pr.msgData)
}
