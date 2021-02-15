
package types

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/types"
	comm " github.com/Daironode/aingle/p2pserver/common"
)

type Block struct {
	Blk        *types.Block
	MerkleRoot common.Uint256
	CCMsg      *types.CrossChainMsg
}

//Serialize message payload
func (this *Block) Serialization(sink *common.ZeroCopySink) {
	this.Blk.Serialization(sink)
	sink.WriteHash(this.MerkleRoot)
	sink.WriteBool(this.CCMsg != nil)
	if this.CCMsg != nil {
		this.CCMsg.Serialization(sink)
	}
}

func (this *Block) CmdType() string {
	return comm.BLOCK_TYPE
}

//Deserialize message payload
func (this *Block) Deserialization(source *common.ZeroCopySource) error {
	this.Blk = new(types.Block)
	err := this.Blk.Deserialization(source)
	if err != nil {
		return fmt.Errorf("read Blk error. err:%v", err)
	}
	var eof bool
	this.MerkleRoot, eof = source.NextHash()
	if eof {
		// to accept old node's block
		this.MerkleRoot = common.UINT256_EMPTY
	}
	hasCCM, irr, eof := source.NextBool()
	if irr || eof {
		// to accept old node's cross msg
		return nil
	}
	var ccMsg *types.CrossChainMsg
	if hasCCM {
		ccMsg = new(types.CrossChainMsg)
		if err := ccMsg.Deserialization(source); err != nil {
			return err
		}
	}
	this.CCMsg = ccMsg

	return nil
}
