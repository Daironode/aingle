 
package common

import (
	"errors"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/event"
)

var ErrNotFound = errors.New("not found")

//Store iterator for iterate store
type StoreIterator interface {
	Next() bool //Next item. If item available return true, otherwise return false
	//Prev() bool           //previous item. If item available return true, otherwise return false
	First() bool //First item. If item available return true, otherwise return false
	//Last() bool           //Last item. If item available return true, otherwise return false
	//Seek(key []byte) bool //Seek key. If item available return true, otherwise return false
	Key() []byte   //Return the current item key
	Value() []byte //Return the current item value
	Release()      //Close iterator
	Error() error  // Error returns any accumulated error.
}

//PersistStore of ledger
type PersistStore interface {
	Put(key []byte, value []byte) error      //Put the key-value pair to store
	Get(key []byte) ([]byte, error)          //Get the value if key in store
	Has(key []byte) (bool, error)            //Whether the key is exist in store
	Delete(key []byte) error                 //Delete the key in store
	NewBatch()                               //Start commit batch
	BatchPut(key []byte, value []byte)       //Put a key-value pair to batch
	BatchDelete(key []byte)                  //Delete the key in batch
	BatchCommit() error                      //Commit batch to store
	Close() error                            //Close store
	NewIterator(prefix []byte) StoreIterator //Return the iterator of store
}

//EventStore save event notify
type EventStore interface {
	//SaveEventNotifyByTx save event notify gen by smart contract execution
	SaveEventNotifyByTx(txHash common.Uint256, notify *event.ExecuteNotify) error
	//Save transaction hashes which have event notify gen
	SaveEventNotifyByBlock(height uint32, txHashs []common.Uint256)
	//GetEventNotifyByTx return event notify by transaction hash
	GetEventNotifyByTx(txHash common.Uint256) (*event.ExecuteNotify, error)
	//Commit event notify to store
	CommitTo() error
}
