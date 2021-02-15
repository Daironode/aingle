 
package neovm

import (
	"fmt"

	" github.com/Daironode/aingle/vm/neovm/errors"
	" github.com/Daironode/aingle/vm/neovm/types"
)

const initialStackCap = 64 // to avoid reallocation

type ValueStack struct {
	data  []types.VmValue
	limit int64
}

func NewValueStack(limit int64) *ValueStack {
	return &ValueStack{
		data:  make([]types.VmValue, 0, initialStackCap),
		limit: limit,
	}
}

func (self *ValueStack) Count() int {
	return len(self.data)
}

func (self *ValueStack) Insert(index int64, t types.VmValue) error {
	l := int64(len(self.data))
	if l >= self.limit {
		return errors.ERR_OVER_LIMIT_STACK
	}
	if index > l || index < 0 {
		return errors.ERR_INDEX_OUT_OF_BOUND
	}
	index = l - index
	self.data = append(self.data, t)
	copy(self.data[index+1:], self.data[index:])
	self.data[index] = t
	return nil
}

func (self *ValueStack) Peek(index int64) (value types.VmValue, err error) {
	l := int64(len(self.data))
	if index >= l || index < 0 {
		err = errors.ERR_INDEX_OUT_OF_BOUND
		return
	}
	index = l - index
	value = self.data[index-1]
	return
}

func (self *ValueStack) Remove(index int64) (value types.VmValue, err error) {
	l := int64(len(self.data))
	if index >= l || index < 0 {
		err = errors.ERR_INDEX_OUT_OF_BOUND
		return
	}
	index = l - index
	value = self.data[index-1]
	self.data = append(self.data[:index-1], self.data[index:]...)
	return
}

func (self *ValueStack) Set(index int, t types.VmValue) error {
	l := len(self.data)
	if index >= l || index < 0 {
		return errors.ERR_INDEX_OUT_OF_BOUND
	}
	self.data[index] = t
	return nil
}

func (self *ValueStack) Push(t types.VmValue) error {
	if int64(len(self.data)) >= self.limit {
		return errors.ERR_OVER_STACK_LEN
	}
	self.data = append(self.data, t)
	return nil
}

func (self *ValueStack) PushMany(vals ...types.VmValue) error {
	if int64(len(self.data)+len(vals)) > self.limit {
		return errors.ERR_OVER_STACK_LEN
	}
	self.data = append(self.data, vals...)
	return nil
}

func (self *ValueStack) PushAsArray(vals []types.VmValue) error {

	if int64(len(self.data)+1) > self.limit {
		return errors.ERR_OVER_STACK_LEN
	}
	arrayValue := types.NewArrayValue()
	for _, val := range vals {
		err := arrayValue.Append(val)
		if err != nil {
			return err
		}
	}
	v := types.VmValueFromArrayVal(arrayValue)
	self.data = append(self.data, v)
	return nil
}

func (self *ValueStack) Pop() (value types.VmValue, err error) {
	length := len(self.data)
	if length == 0 {
		err = errors.ERR_INDEX_OUT_OF_BOUND
		return
	}
	value = self.data[length-1]
	self.data = self.data[:length-1]
	return
}

func (self *ValueStack) PopPair() (left, right types.VmValue, err error) {
	right, err = self.Pop()
	if err != nil {
		return
	}
	left, err = self.Pop()
	return
}

func (self *ValueStack) PopTriple() (left, middle, right types.VmValue, err error) {
	middle, right, err = self.PopPair()
	if err != nil {
		return
	}
	left, err = self.Pop()
	return
}

func (self *ValueStack) Swap(i, j int64) error {
	l := int64(len(self.data))
	if i >= l || i < 0 {
		return errors.ERR_INDEX_OUT_OF_BOUND
	}
	if j >= l || j < 0 {
		return errors.ERR_INDEX_OUT_OF_BOUND
	}
	if i == j {
		return nil
	}
	self.data[l-i-1], self.data[l-j-1] = self.data[l-j-1], self.data[l-i-1]

	return nil
}

func (self *ValueStack) CopyTo(stack *ValueStack) error {
	if int64(len(self.data)+len(stack.data)) > stack.limit {
		return errors.ERR_OVER_STACK_LEN
	}
	stack.data = append(stack.data, self.data...)
	return nil
}

func (self *ValueStack) Dump() string {
	data := fmt.Sprintf("stack[%d]:\n", len(self.data))
	for i, item := range self.data {
		i = len(self.data) - i
		res := item.Dump()
		data += fmt.Sprintf("%d:\t%s\n", i, res)
	}
	return data
}
