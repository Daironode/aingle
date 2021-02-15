
package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
)

const I128_SIZE = 16

type U128 [I128_SIZE]byte // little endian u128
type I128 [I128_SIZE]byte // little endian i128

var U128_EMPTY = U128{}
var oneBits128 = func() I128 {
	var i128 I128
	for i := 0; i < I128_SIZE; i++ {
		i128[i] = 255
	}
	return i128
}()

var pow128 = bigPow(2, 128)
var maxBigU128 = new(big.Int).Sub(bigPow(2, 128), big.NewInt(1))
var maxI128 = new(big.Int).Sub(bigPow(2, 127), big.NewInt(1))
var minI128 = new(big.Int).Neg(bigPow(2, 127))

// returns a ** b
func bigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

func I128FromUint64(val uint64) I128 {
	var i128 I128
	binary.LittleEndian.PutUint64(i128[:], val)
	return i128
}

func I128FromInt64(val int64) I128 {
	var i128 I128
	if val < 0 {
		i128 = oneBits128
	}
	binary.LittleEndian.PutUint64(i128[:], uint64(val))

	return i128
}

// val should in i128 range
func I128FromBigInt(val *big.Int) (I128, error) {
	var u128 I128
	if val.Cmp(maxI128) > 0 || val.Cmp(minI128) < 0 {
		return u128, errors.New("big int out of i128 range")
	}

	if val.Sign() < 0 {
		val = new(big.Int).Add(val, pow128)
	}
	buf := val.Bytes()
	buf = ToArrayReverse(buf)
	copy(u128[:], buf)

	return u128, nil
}

func (self U128) ToBigInt() *big.Int {
	buf := append(self[:], 0)
	buf = ToArrayReverse(buf)
	return new(big.Int).SetBytes(buf)
}

func (self U128) ToI128() I128 {
	return I128(self)
}

func (self I128) ToBigInt() *big.Int {
	val := U128(self).ToBigInt()

	if val.Cmp(maxI128) > 0 {
		val.Sub(val, pow128)
	}

	return val
}

// to big endian hex string
func (self *I128) ToBEHex() string {
	return fmt.Sprintf("%x", ToArrayReverse(self[:]))
}

func (self *I128) ToNumString() string {
	val := self.ToBigInt()
	return fmt.Sprintf("%d", val)
}

// to little endian hex string
func (self *I128) ToLEHex() string {
	return fmt.Sprintf("%x", self[:])
}
