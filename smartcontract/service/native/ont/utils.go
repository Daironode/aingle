
package ont

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/constants"
	cstates " github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/event"
	" github.com/Daironode/aingle/smartcontract/service/native"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

const (
	UNBOUND_TIME_OFFSET       = "unboundTimeOffset"
	TOTAL_SUPPLY_NAME         = "totalSupply"
	INIT_NAME                 = "init"
	TRANSFER_NAME             = "transfer"
	APPROVE_NAME              = "approve"
	TRANSFERFROM_NAME         = "transferFrom"
	NAME_NAME                 = "name"
	SYMBOL_NAME               = "symbol"
	DECIMALS_NAME             = "decimals"
	TOTALSUPPLY_NAME          = "totalSupply"
	BALANCEOF_NAME            = "balanceOf"
	ALLOWANCE_NAME            = "allowance"
	TOTAL_ALLOWANCE_NAME      = "totalAllowance"
	UNBOUND_ONG_TO_GOVERNANCE = "unboundOngToGovernance"
)

func AddNotifications(native *native.NativeService, contract common.Address, state *State) {
	if !config.DefConfig.Common.EnableEventLog {
		return
	}
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{TRANSFER_NAME, state.From.ToBase58(), state.To.ToBase58(), state.Value},
		})
}
func GetToUInt64StorageItem(toBalance, value uint64) *cstates.StorageItem {
	sink := common.NewZeroCopySink(nil)
	sink.WriteUint64(toBalance + value)
	return &cstates.StorageItem{Value: sink.Bytes()}
}

func GenTotalSupplyKey(contract common.Address) []byte {
	return append(contract[:], TOTAL_SUPPLY_NAME...)
}

func GenBalanceKey(contract, addr common.Address) []byte {
	return append(contract[:], addr[:]...)
}

func Transfer(native *native.NativeService, contract common.Address, state *State) (uint64, uint64, error) {
	if !native.ContextRef.CheckWitness(state.From) {
		return 0, 0, errors.NewErr("authentication failed!")
	}

	fromBalance, err := fromTransfer(native, GenBalanceKey(contract, state.From), state.Value)
	if err != nil {
		return 0, 0, err
	}

	toBalance, err := toTransfer(native, GenBalanceKey(contract, state.To), state.Value)
	if err != nil {
		return 0, 0, err
	}
	return fromBalance, toBalance, nil
}

func GenApproveKey(contract, from, to common.Address) []byte {
	temp := append(contract[:], from[:]...)
	return append(temp, to[:]...)
}

func TransferedFrom(native *native.NativeService, currentContract common.Address, state *TransferFrom) (uint64, uint64, error) {
	if native.Time <= config.GetOntHolderUnboundDeadline()+constants.GENESIS_BLOCK_TIMESTAMP {
		if !native.ContextRef.CheckWitness(state.Sender) {
			return 0, 0, errors.NewErr("authentication failed!")
		}
	} else {
		if state.Sender != state.To && !native.ContextRef.CheckWitness(state.Sender) {
			return 0, 0, errors.NewErr("authentication failed!")
		}
	}

	if err := fromApprove(native, genTransferFromKey(currentContract, state), state.Value); err != nil {
		return 0, 0, err
	}

	fromBalance, err := fromTransfer(native, GenBalanceKey(currentContract, state.From), state.Value)
	if err != nil {
		return 0, 0, err
	}

	toBalance, err := toTransfer(native, GenBalanceKey(currentContract, state.To), state.Value)
	if err != nil {
		return 0, 0, err
	}
	return fromBalance, toBalance, nil
}

func getUnboundOffset(native *native.NativeService, contract, address common.Address) (uint32, error) {
	offset, err := utils.GetStorageUInt32(native, genAddressUnboundOffsetKey(contract, address))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func getGovernanceUnboundOffset(native *native.NativeService, contract common.Address) (uint32, error) {
	offset, err := utils.GetStorageUInt32(native, genGovernanceUnboundOffsetKey(contract))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func genTransferFromKey(contract common.Address, state *TransferFrom) []byte {
	temp := append(contract[:], state.From[:]...)
	return append(temp, state.Sender[:]...)
}

func fromApprove(native *native.NativeService, fromApproveKey []byte, value uint64) error {
	approveValue, err := utils.GetStorageUInt64(native, fromApproveKey)
	if err != nil {
		return err
	}
	if approveValue < value {
		return fmt.Errorf("[TransferFrom] approve balance insufficient! have %d, got %d", approveValue, value)
	} else if approveValue == value {
		native.CacheDB.Delete(fromApproveKey)
	} else {
		native.CacheDB.Put(fromApproveKey, utils.GenUInt64StorageItem(approveValue-value).ToArray())
	}
	return nil
}

func fromTransfer(native *native.NativeService, fromKey []byte, value uint64) (uint64, error) {
	fromBalance, err := utils.GetStorageUInt64(native, fromKey)
	if err != nil {
		return 0, err
	}
	if fromBalance < value {
		addr, _ := common.AddressParseFromBytes(fromKey[20:])
		return 0, fmt.Errorf("[Transfer] balance insufficient. contract:%s, account:%s,balance:%d, transfer amount:%d",
			native.ContextRef.CurrentContext().ContractAddress.ToHexString(), addr.ToBase58(), fromBalance, value)
	} else if fromBalance == value {
		native.CacheDB.Delete(fromKey)
	} else {
		native.CacheDB.Put(fromKey, utils.GenUInt64StorageItem(fromBalance-value).ToArray())
	}
	return fromBalance, nil
}

func toTransfer(native *native.NativeService, toKey []byte, value uint64) (uint64, error) {
	toBalance, err := utils.GetStorageUInt64(native, toKey)
	if err != nil {
		return 0, err
	}
	native.CacheDB.Put(toKey, GetToUInt64StorageItem(toBalance, value).ToArray())
	return toBalance, nil
}

func genAddressUnboundOffsetKey(contract, address common.Address) []byte {
	temp := append(contract[:], UNBOUND_TIME_OFFSET...)
	return append(temp, address[:]...)
}

func genGovernanceUnboundOffsetKey(contract common.Address) []byte {
	temp := append(contract[:], UNBOUND_TIME_OFFSET...)
	return temp
}