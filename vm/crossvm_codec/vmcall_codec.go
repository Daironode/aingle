
package crossvm_codec

import (
	"bytes"

	" github.com/Daironode/aingle/common"
)

//input byte array should be the following format
// version(1byte) + type(1byte) + data...
func DeserializeCallParam(input []byte) (interface{}, error) {
	if !bytes.HasPrefix(input, []byte{0}) {
		return nil, ERROR_PARAM_FORMAT
	}

	source := common.NewZeroCopySource(input[1:])
	return DecodeValue(source)
}
