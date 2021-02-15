
package account

import (
	" github.com/Daironode/aingle/common"
)

type ClientStore interface {
	BuildDatabase(path string)

	SaveStoredData(name string, value []byte)

	LoadStoredData(name string) []byte

	LoadAccount() map[common.Address]*Account
}
