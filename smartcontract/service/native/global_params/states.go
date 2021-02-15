
package global_params

import (
	"fmt"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/errors"
	" github.com/Daironode/aingle/smartcontract/service/native/utils"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

type ParamNameList []string

func (params *Params) SetParam(value Param) {
	for index, param := range *params {
		if param.Key == value.Key {
			(*params)[index] = value
			return
		}
	}
	*params = append(*params, value)
}

func (params *Params) GetParam(key string) (int, Param) {
	for index, param := range *params {
		if param.Key == key {
			return index, param
		}
	}
	return -1, Param{}
}

func (params *Params) Serialization(sink *common.ZeroCopySink) {
	paramNum := len(*params)
	utils.EncodeVarUint(sink, uint64(paramNum))
	for _, param := range *params {
		sink.WriteString(param.Key)
		sink.WriteString(param.Value)
	}
}
func (params *Params) Deserialization(source *common.ZeroCopySource) error {
	paramNum, err := utils.DecodeVarUint(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "param config, deserialize params length error!")
	}

	for i := 0; uint64(i) < paramNum; i++ {
		param := Param{}
		var irregular, eof bool
		param.Key, _, irregular, eof = source.NextString()
		if irregular || eof {
			return errors.NewDetailErr(err, errors.ErrNoCode, fmt.Sprintf("param config, deserialize param key %v error!", param.Key))
		}
		param.Value, _, irregular, eof = source.NextString()
		if irregular || eof {
			return errors.NewDetailErr(err, errors.ErrNoCode, fmt.Sprintf("param config, deserialize param value %v error!", param.Value))
		}
		*params = append(*params, param)
	}
	return nil
}
func (nameList *ParamNameList) Serialization(sink *common.ZeroCopySink) {
	nameNum := len(*nameList)
	utils.EncodeVarUint(sink, uint64(nameNum))
	for _, value := range *nameList {
		sink.WriteString(value)
	}
}

func (nameList *ParamNameList) Deserialization(source *common.ZeroCopySource) error {
	nameNum, err := utils.DecodeVarUint(source)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "param config, deserialize param name list length error!")
	}
	for i := 0; uint64(i) < nameNum; i++ {
		name, _, irregular, eof := source.NextString()
		if irregular || eof {
			return errors.NewDetailErr(err, errors.ErrNoCode, fmt.Sprintf("param config, deserialize param name %v error!", name))
		}
		*nameList = append(*nameList, name)
	}
	return nil
}
