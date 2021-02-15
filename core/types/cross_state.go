 package types

import (
	"crypto/sha256"
	"fmt"

	" github.com/Daironode/aingle/common"
)

const (
	CURR_CROSS_STATES_VERSION = 0
)

type CrossChainMsg struct {
	Version    byte
	Height     uint32
	StatesRoot common.Uint256

	SigData [][]byte

	hash *common.Uint256
}

func (this *CrossChainMsg) serializationUnsigned(sink *common.ZeroCopySink) {
	sink.WriteByte(this.Version)
	sink.WriteUint32(this.Height)
	sink.WriteBytes(this.StatesRoot[:])
}

func (this *CrossChainMsg) Serialization(sink *common.ZeroCopySink) {
	this.serializationUnsigned(sink)
	sink.WriteVarUint(uint64(len(this.SigData)))
	for _, sig := range this.SigData {
		sink.WriteVarBytes(sig)
	}
}

func (this *CrossChainMsg) Deserialization(source *common.ZeroCopySource) error {
	var eof bool
	this.Version, eof = source.NextByte()
	if eof {
		return fmt.Errorf("CrossChainMsg, deserialization read version error")
	}
	this.Height, eof = source.NextUint32()
	if eof {
		return fmt.Errorf("CrossChainMsg, deserialization read height error")
	}
	this.StatesRoot, eof = source.NextHash()
	if eof {
		return fmt.Errorf("CrossChainMsg, deserialization read statesRoot error")
	}
	sigLen, _, irr, eof := source.NextVarUint()
	if irr || eof {
		return fmt.Errorf("CrossChainMsg, deserialization read sigData lenght error")
	}
	sigData := make([][]byte, 0, sigLen)
	for i := 0; i < int(sigLen); i++ {
		v, _, irr, eof := source.NextVarBytes()
		if irr || eof {
			return fmt.Errorf("CrossChainMsg, deserialization read sigData value error")
		}
		sigData = append(sigData, v)
	}
	this.SigData = sigData
	return nil
}

func (this *CrossChainMsg) Hash() common.Uint256 {
	if this.hash != nil {
		return *this.hash
	}
	sink := common.NewZeroCopySink(nil)
	this.serializationUnsigned(sink)
	temp := sha256.Sum256(sink.Bytes())
	hash := common.Uint256(sha256.Sum256(temp[:]))
	this.hash = &hash
	return hash
}

func (this *CrossChainMsg) SetHash(hash common.Uint256) {
	this.hash = &hash
}
