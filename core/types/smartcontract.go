
package types

import " github.com/Daironode/aingle/common"

type SmartCodeEvent struct {
	TxHash common.Uint256
	Action string
	Result interface{}
	Error  int64
}
