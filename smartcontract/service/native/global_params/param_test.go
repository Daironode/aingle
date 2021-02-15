
package global_params

import (
	"strconv"
	"testing"

	" github.com/Daironode/aingle/common"
	"github.com/stretchr/testify/assert"
)

func TestParams_Serialize_Deserialize(t *testing.T) {
	params := Params{}
	for i := 0; i < 10; i++ {
		k := "key" + strconv.Itoa(i)
		v := "value" + strconv.Itoa(i)
		params.SetParam(Param{k, v})
	}
	sink := common.NewZeroCopySink(nil)
	params.Serialization(sink)
	deserializeParams := Params{}
	source := common.NewZeroCopySource(sink.Bytes())
	if err := deserializeParams.Deserialization(source); err != nil {
		t.Fatalf("params deserialize error: %v", err)
	}
	for i := 0; i < 10; i++ {
		originParam := params[i]
		deseParam := deserializeParams[i]
		if originParam.Key != deseParam.Key || originParam.Value != deseParam.Value {
			t.Fatal("params deserialize error")
		}
	}
}

func TestParamNameList_Serialize_Deserialize(t *testing.T) {
	nameList := ParamNameList{}
	for i := 0; i < 3; i++ {
		nameList = append(nameList, strconv.Itoa(i))
	}
	sink := common.NewZeroCopySink(nil)
	nameList.Serialization(sink)
	deserializeNameList := ParamNameList{}
	source := common.NewZeroCopySource(sink.Bytes())
	err := deserializeNameList.Deserialization(source)
	assert.Nil(t, err)
	assert.Equal(t, nameList, deserializeNameList)
}
