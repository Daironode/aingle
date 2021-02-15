 
package init

import (
	"bytes"
	"math/big"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/service/native/auth"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/cross_chain_manager"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/header_sync"
	" github.com/Daironode/aingle/smartcontract/service/native/cross_chain/lock_proxy"
	params " github.com/Daironode/aingle/smartcontract/service/native/global_params"
	" github.com/Daironode/aingle/smartcontract/service/native/governance"
	" github.com/Daironode/aingle/smartcontract/service/native/ong"
	" github.com/Daironode/aingle/smartcontract/service/native/ont"
	" github.com/Daironode/aingle/smartcontract/service/native/ontfs"
	" github.com/Daironode/aingle/smartcontract/service/native/ontid"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
	" github.com/Daironode/aingle/smartcontract/service/neovm"
	vm " github.com/Daironode/aingle/vm/neovm"
)

var (
	COMMIT_DPOS_BYTES = InitBytes(utils.GovernanceContractAddress, governance.COMMIT_DPOS)
)

func init() {
	ong.InitOng()
	ont.InitOnt()
	params.InitGlobalParams()
	ontid.Init()
	auth.Init()
	governance.InitGovernance()
	cross_chain_manager.InitCrossChain()
	header_sync.InitHeaderSync()
	lock_proxy.InitLockProxy()
	ontfs.InitFs()
}

func InitBytes(addr common.Address, method string) []byte {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	return builder.ToArray()
}
