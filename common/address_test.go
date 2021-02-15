
package common

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressFromBase58(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	base58 := addr.ToBase58()
	b1 := string(append([]byte{'X'}, []byte(base58)...))
	_, err := AddressFromBase58(b1)

	assert.NotNil(t, err)

	b2 := string([]byte(base58)[1:10])
	_, err = AddressFromBase58(b2)

	assert.NotNil(t, err)
}

func TestAddressParseFromBytes(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	addr2, _ := AddressParseFromBytes(addr[:])

	assert.Equal(t, addr, addr2)
}

func TestAddress_Serialize(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	sink := NewZeroCopySink(nil)
	addr.Serialization(sink)

	var addr2 Address
	source := NewZeroCopySource(sink.Bytes())
	addr2.Deserialization(source)
	assert.Equal(t, addr, addr2)
}
