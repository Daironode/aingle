 package crossvm_codec

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/common/log"
)

func DeserializeNotify(input []byte) interface{} {
	val, err := parseNotify(input)
	if err != nil {
		return input
	}

	return stringify(val)
}

func stringify(notify interface{}) interface{} {
	switch val := notify.(type) {
	case []byte:
		return hex.EncodeToString(val)
	case common.Address:
		return val.ToBase58()
	case bool, string:
		return val
	case common.Uint256:
		return val.ToHexString()
	case *big.Int:
		return fmt.Sprintf("%d", val)
	case []interface{}:
		list := make([]interface{}, 0, len(val))
		for _, v := range val {
			list = append(list, stringify(v))
		}
		return list
	default:
		log.Warn("notify codec: unsupported type:", reflect.TypeOf(val).String())

		return val
	}
}

// input byte array should be the following format
// evt\0(4byte) + type(1byte) + usize( bytearray or list) (4 bytes) + data...
func parseNotify(input []byte) (interface{}, error) {
	if !bytes.HasPrefix(input, []byte("evt\x00")) {
		return nil, ERROR_PARAM_FORMAT
	}

	source := common.NewZeroCopySource(input[4:])

	return DecodeValue(source)
}
