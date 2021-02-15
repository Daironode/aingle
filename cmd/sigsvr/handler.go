
package sigsvr

import " github.com/Daironode/aingle/cmd/sigsvr/handlers"

func init() {
	DefCliRpcSvr.RegHandler("createaccount", handlers.CreateAccount)
	DefCliRpcSvr.RegHandler("exportaccount", handlers.ExportAccount)
	DefCliRpcSvr.RegHandler("sigdata", handlers.SigData)
	DefCliRpcSvr.RegHandler("sigrawtx", handlers.SigRawTransaction)
	DefCliRpcSvr.RegHandler("sigmutilrawtx", handlers.SigMutilRawTransaction)
	DefCliRpcSvr.RegHandler("sigtransfertx", handlers.SigTransferTransaction)
	DefCliRpcSvr.RegHandler("signeovminvoketx", handlers.SigNeoVMInvokeTx)
	DefCliRpcSvr.RegHandler("signeovminvokeabitx", handlers.SigNeoVMInvokeAbiTx)
	DefCliRpcSvr.RegHandler("signativeinvoketx", handlers.SigNativeInvokeTx)
}
