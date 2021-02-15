 package utils

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	cstates " github.com/Daironode/aingle/core/states"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native"
)

type LinkedlistNode struct {
	next    []byte
	prev    []byte
	payload []byte
}

func (this *LinkedlistNode) GetPrevious() []byte {
	return this.prev
}

func (this *LinkedlistNode) GetNext() []byte {
	return this.next
}

func (this *LinkedlistNode) GetPayload() []byte {
	return this.payload
}

func makeLinkedlistNode(next []byte, prev []byte, payload []byte) ([]byte, error) {
	node := &LinkedlistNode{next: next, prev: prev, payload: payload}
	node_bytes, err := node.Serialization()
	if err != nil {
		return nil, err
	}
	return node_bytes, nil
}
func (this *LinkedlistNode) Serialization() ([]byte, error) {
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarBytes(this.next)
	sink.WriteVarBytes(this.prev)
	sink.WriteVarBytes(this.payload)
	return sink.Bytes(), nil
}

func (this *LinkedlistNode) Deserialization(r []byte) error {
	source := common.NewZeroCopySource(r)
	next, err := DecodeVarBytes(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[linked list] deserialize next error!")
	}
	prev, err := DecodeVarBytes(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[linked list] deserialize prev error!")
	}
	payload, err := DecodeVarBytes(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[linked list] deserialize payload error!")
	}
	this.next = next
	this.prev = prev
	this.payload = payload
	return nil
}

func getListHead(native *native.NativeService, index []byte) ([]byte, error) {
	head, err := native.CacheDB.Get(index)
	if err != nil {
		return nil, err
	}
	if head == nil {
		return nil, nil
	}
	value, err := cstates.GetValueFromRawStorageItem(head)
	if err != nil {
		return nil, fmt.Errorf("[linked list] get header error:%v", err)
	}
	return value, nil
}

func getListNode(native *native.NativeService, index []byte, item []byte) (*LinkedlistNode, error) {
	node := new(LinkedlistNode)
	data, err := native.CacheDB.Get(append(index, item...))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	rawNode, err := cstates.GetValueFromRawStorageItem(data)
	if err != nil {
		return nil, fmt.Errorf("[linked list] get list node error:%v", err)
	}
	if len(rawNode) == 0 {
		return nil, nil
	}
	err = node.Deserialization(rawNode)
	if err != nil {
		//log.Tracef("[index: %s, item: %s] error %s", hex.EncodeToString(index), hex.EncodeToString(item), err)
		return nil, err
	}
	return node, nil
}

func LinkedlistInsert(native *native.NativeService, index []byte, item []byte, payload []byte) error {
	null := []byte{}
	if item == nil {
		return errors.NewErr("[linked list] invalid item")
	}
	head, err := getListHead(native, index) //list head
	if err != nil {
		//log.Trace(err)
		return err
	}

	q, err := getListNode(native, index, item) //list node
	if err != nil {
		//log.Trace(err)
		return err
	}

	if q != nil { //already exists
		//log.Trace(err)
		node, err := makeLinkedlistNode(q.next, q.prev, payload)
		if err != nil {
			return err
		}
		PutBytes(native, append(index, item...), node) //update it
		return nil
	}
	if head == nil { //doubly-linked list contains zero element
		node, err := makeLinkedlistNode(null, null, payload)
		if err != nil {
			//log.Trace(err)
			return err
		}
		PutBytes(native, append(index, item...), node) //item is the only element
		PutBytes(native, index, item)                  //item becomes head
	} else {
		null := []byte{}
		node, err := makeLinkedlistNode(head, null, payload)
		if err != nil {
			//log.Trace(err)
			return err
		}
		PutBytes(native, append(index, item...), node) //item.next = head, item.prev = null,
		// item.payload = payload
		qhead, err := getListNode(native, index, head)
		if err != nil {
			//log.Trace(err)
			return err
		}

		node, err = makeLinkedlistNode(qhead.next, item, qhead.payload)
		if err != nil {
			//log.Trace(err)
			return err
		}
		PutBytes(native, append(index, head...), node) //head.next = head.next, head.prev = item,
		// head.payload = head.payload
		PutBytes(native, index, item) // item becomes head
	}
	return nil
}

func LinkedlistDelete(native *native.NativeService, index []byte, item []byte) (bool, error) {
	null := []byte{}
	if item == nil {
		return false, errors.NewErr("[linked list] invalid item")
	}
	q, err := getListNode(native, index, item)
	if err != nil {
		return false, err
	}
	if q == nil {
		return false, nil
	}

	prev, next := q.prev, q.next
	if prev == nil {
		if next == nil {
			//clear linked list
			native.CacheDB.Delete(index)
		} else {
			qnext, err := getListNode(native, index, next)
			if err != nil {
				return false, err
			}
			node, err := makeLinkedlistNode(qnext.next, null, qnext.payload) //qnext.next = qnext.next
			if err != nil {                                                  // qnext.prev = nil
				return false, err
			}
			PutBytes(native, append(index, next...), node)
			PutBytes(native, index, next) //next becomes head
		}
	} else {
		if next == nil {
			qprev, err := getListNode(native, index, prev)
			if err != nil {
				return false, err
			}
			node, err := makeLinkedlistNode(null, qprev.prev, qprev.payload) //qprev becomes end
			if err != nil {
				return false, err
			}
			PutBytes(native, append(index, prev...), node)
		} else {
			qprev, err := getListNode(native, index, prev)
			if err != nil {
				return false, err
			}
			qnext, err := getListNode(native, index, next)
			if err != nil {
				return false, err
			}
			node_prev, err := makeLinkedlistNode(next, qprev.prev, qprev.payload) //
			if err != nil {
				return false, err
			}
			node_next, err := makeLinkedlistNode(qnext.next, prev, qnext.payload)
			if err != nil {
				return false, err
			}
			PutBytes(native, append(index, prev...), node_prev)
			PutBytes(native, append(index, next...), node_next)
		}
	}
	native.CacheDB.Delete(append(index, item...))
	return true, nil
}

func LinkedlistGetItem(native *native.NativeService, index []byte, item []byte) (*LinkedlistNode, error) {
	if item == nil {
		return nil, errors.NewErr("[linkedlist getNext] item is nil")
	}
	q, err := getListNode(native, index, item)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func LinkedlistGetHead(native *native.NativeService, index []byte) ([]byte, error) {
	head, err := getListHead(native, index)
	if err != nil {
		return nil, err
	}
	return head, nil
}

func LinkedlistGetNumOfItems(native *native.NativeService, index []byte) (int, error) {
	n := 0
	head, err := getListHead(native, index)
	if err != nil {
		return 0, err
	}
	q := head
	for q != nil {
		n += 1
		qnode, err := getListNode(native, index, q)
		if err != nil {
			return 0, err
		}
		q = qnode.next
	}
	return n, nil
}

func LinkedlistDeleteAll(native *native.NativeService, index []byte) error {
	head, err := getListHead(native, index)
	if err != nil {
		return err
	}
	q := head
	for q != nil {
		qnode, err := getListNode(native, index, q)
		if err != nil {
			return err
		}
		native.CacheDB.Delete(append(index, q...))
		q = qnode.next
	}
	native.CacheDB.Delete(index)
	return nil
}
