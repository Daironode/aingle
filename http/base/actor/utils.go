
// Package actor privides communication with other actor
package actor

import (
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

func updateNativeSCAddr(hash common.Address) common.Address {
	if hash == utils.OntContractAddress {
		hash = common.AddressFromVmCode(utils.OntContractAddress[:])
	} else if hash == utils.OngContractAddress {
		hash = common.AddressFromVmCode(utils.OngContractAddress[:])
	} else if hash == utils.OntIDContractAddress {
		hash = common.AddressFromVmCode(utils.OntIDContractAddress[:])
	} else if hash == utils.ParamContractAddress {
		hash = common.AddressFromVmCode(utils.ParamContractAddress[:])
	} else if hash == utils.AuthContractAddress {
		hash = common.AddressFromVmCode(utils.AuthContractAddress[:])
	} else if hash == utils.GovernanceContractAddress {
		hash = common.AddressFromVmCode(utils.GovernanceContractAddress[:])
	}
	return hash
}
