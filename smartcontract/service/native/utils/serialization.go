
package utils

import (
	"fmt"
	"io"
	"math/big"

	" github.com/Daironode/aingle/common"
)

func EncodeAddress(sink *common.ZeroCopySink, addr common.Address) (size uint64) {
	return sink.WriteVarBytes(addr[:])
}

func EncodeVarUint(sink *common.ZeroCopySink, value uint64) (size uint64) {
	return sink.WriteVarBytes(common.BigIntToNeoBytes(big.NewInt(int64(value))))
}

func EncodeVarBytes(sink *common.ZeroCopySink, v []byte) (size uint64) {
	return sink.WriteVarBytes(v)
}

func EncodeString(sink *common.ZeroCopySink, str string) (size uint64) {
	return sink.WriteVarBytes([]byte(str))
}

func EncodeBool(sink *common.ZeroCopySink, value bool) {
	sink.WriteBool(value)
}

func DecodeVarUint(source *common.ZeroCopySource) (uint64, error) {
	value, _, irregular, eof := source.NextVarBytes()
	if eof {
		return 0, io.ErrUnexpectedEOF
	}
	if irregular {
		return 0, common.ErrIrregularData
	}
	v := common.BigIntFromNeoBytes(value)
	if v.Cmp(big.NewInt(0)) < 0 {
		return 0, fmt.Errorf("%s", "value should not be a negative number.")
	}
	return v.Uint64(), nil
}

func DecodeAddress(source *common.ZeroCopySource) (common.Address, error) {
	from, _, irregular, eof := source.NextVarBytes()
	if eof {
		return common.Address{}, io.ErrUnexpectedEOF
	}
	if irregular {
		return common.Address{}, common.ErrIrregularData
	}

	return common.AddressParseFromBytes(from)
}
func DecodeVarBytes(source *common.ZeroCopySource) ([]byte, error) {
	data, _, irregular, eof := source.NextVarBytes()
	if eof {
		return nil, io.ErrUnexpectedEOF
	}
	if irregular {
		return nil, common.ErrIrregularData
	}

	return data, nil
}
func DecodeUint64(source *common.ZeroCopySource) (uint64, error) {
	data, eof := source.NextUint64()
	if eof {
		return 0, io.ErrUnexpectedEOF
	}
	return data, nil
}
func DecodeUint32(source *common.ZeroCopySource) (uint32, error) {
	data, eof := source.NextUint32()
	if eof {
		return 0, io.ErrUnexpectedEOF
	}
	return data, nil
}
func DecodeBool(source *common.ZeroCopySource) (bool, error) {
	data, irregular, eof := source.NextBool()
	if eof {
		return false, io.ErrUnexpectedEOF
	}
	if irregular {
		return false, common.ErrIrregularData
	}
	return data, nil
}
func DecodeString(source *common.ZeroCopySource) (string, error) {
	data, _, irregular, eof := source.NextString()
	if eof {
		return "", io.ErrUnexpectedEOF
	}
	if irregular {
		return "", common.ErrIrregularData
	}

	return data, nil
}
