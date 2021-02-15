
package neovm

import (
	"math/big"
)

func ToBigInt(data interface{}) *big.Int {
	var bi big.Int
	switch t := data.(type) {
	case int64:
		bi.SetInt64(int64(t))
	case int32:
		bi.SetInt64(int64(t))
	case int16:
		bi.SetInt64(int64(t))
	case int8:
		bi.SetInt64(int64(t))
	case int:
		bi.SetInt64(int64(t))
	case uint64:
		bi.SetUint64(uint64(t))
	case uint32:
		bi.SetUint64(uint64(t))
	case uint16:
		bi.SetUint64(uint64(t))
	case uint8:
		bi.SetUint64(uint64(t))
	case uint:
		bi.SetUint64(uint64(t))
	case big.Int:
		bi = t
	case *big.Int:
		bi = *t
	}
	return &bi
}

func BigIntZip(ints1 *big.Int, ints2 *big.Int, op OpCode) *big.Int {
	nb := new(big.Int)
	switch op {
	case AND:
		nb.And(ints1, ints2)
	case OR:
		nb.Or(ints1, ints2)
	case XOR:
		nb.Xor(ints1, ints2)
	case ADD:
		nb.Add(ints1, ints2)
	case SUB:
		nb.Sub(ints1, ints2)
	case MUL:
		nb.Mul(ints1, ints2)
	case DIV:
		nb.Quo(ints1, ints2)
	case MOD:
		nb.Rem(ints1, ints2)
	case SHL:
		nb.Lsh(ints1, uint(ints2.Int64()))
	case SHR:
		nb.Rsh(ints1, uint(ints2.Int64()))
	case MIN:
		c := ints1.Cmp(ints2)
		if c <= 0 {
			nb.Set(ints1)
		} else {
			nb.Set(ints2)
		}
	case MAX:
		c := ints1.Cmp(ints2)
		if c <= 0 {
			nb.Set(ints2)
		} else {
			nb.Set(ints1)
		}
	}
	return nb
}
