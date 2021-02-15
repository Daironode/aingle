 package testsuite

import (
	" github.com/Daironode/aingle/core/store/leveldbstore"
	" github.com/Daironode/aingle/core/store/overlaydb"
)

func NewOverlayDB() *overlaydb.OverlayDB {
	store, _ := leveldbstore.NewMemLevelDBStore()
	return overlaydb.NewOverlayDB(store)
}
