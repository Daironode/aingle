
package types

import (
	"fmt"

	" github.com/Daironode/aingle/vm/neovm/constants"
	" github.com/Daironode/aingle/vm/neovm/errors"
)

const (
	MAX_STRUCT_DEPTH = 10
	MAX_CLONE_LENGTH = 1024
)

// struct value is value type
type StructValue struct {
	Data []VmValue
}

func NewStructValue() *StructValue {
	return &StructValue{Data: make([]VmValue, 0, initArraySize)}
}

func (self *StructValue) Append(item VmValue) error {
	if len(self.Data) >= constants.MAX_ARRAY_SIZE {
		return errors.ERR_OVER_MAX_ARRAY_SIZE
	}
	self.Data = append(self.Data, item)
	return nil
}

func (self *StructValue) Len() int64 {
	return int64(len(self.Data))
}

func (self *StructValue) Clone() (*StructValue, error) {
	var i int
	return cloneStruct(self, &i)
}

func cloneStruct(s *StructValue, length *int) (*StructValue, error) {
	if *length > MAX_CLONE_LENGTH {
		return nil, fmt.Errorf("%s", "over max struct clone length")
	}
	var arr []VmValue
	for _, v := range s.Data {
		*length++
		if temp, err := v.AsStructValue(); err == nil {
			vc, err := cloneStruct(temp, length)
			if err != nil {
				return nil, err
			}
			arr = append(arr, VmValueFromStructVal(vc))
		} else {
			arr = append(arr, v)
		}
	}
	return &StructValue{Data: arr}, nil
}
