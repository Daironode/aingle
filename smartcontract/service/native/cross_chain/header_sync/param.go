
package header_sync

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type SyncBlockHeaderParam struct {
	Address common.Address
	Headers [][]byte
}

func (this *SyncBlockHeaderParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeAddress(sink, this.Address)
	utils.EncodeVarUint(sink, uint64(len(this.Headers)))
	for _, v := range this.Headers {
		utils.EncodeVarBytes(sink, v)
	}
}

func (this *SyncBlockHeaderParam) Deserialization(source *common.ZeroCopySource) error {
	address, err := utils.DecodeAddress(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeAddress, deserialize address error:%s", err)
	}
	n, err := utils.DecodeVarUint(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarUint, deserialize header count error:%s", err)
	}
	var headers [][]byte
	for i := 0; uint64(i) < n; i++ {
		header, err := utils.DecodeVarBytes(source)
		if err != nil {
			return fmt.Errorf("utils.DecodeVarBytes, deserialize header error: %v", err)
		}
		headers = append(headers, header)
	}
	this.Address = address
	this.Headers = headers
	return nil
}

type SyncGenesisHeaderParam struct {
	GenesisHeader []byte
}

func (this *SyncGenesisHeaderParam) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarBytes(sink, this.GenesisHeader)
}

func (this *SyncGenesisHeaderParam) Deserialization(source *common.ZeroCopySource) error {
	genesisHeader, err := utils.DecodeVarBytes(source)
	if err != nil {
		return fmt.Errorf("utils.DecodeVarBytes, deserialize genesisHeader count error:%s", err)
	}
	this.GenesisHeader = genesisHeader
	return nil
}
