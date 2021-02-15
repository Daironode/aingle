
package common

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"

	"github.com/Daironode/aingle-crypto/keypair"
)

// GetNonce returns random nonce
func GetNonce() uint64 {
	// Fixme replace with the real random number generator
	nonce := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return nonce
}

// ToHexString convert []byte to hex string
func ToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

// HexToBytes convert hex string to []byte
func HexToBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0; i-- {
		x = append(x, arr[i])
	}
	return x
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func PubKeyToHex(pub keypair.PublicKey) string {
	nodeid := hex.EncodeToString(keypair.SerializePublicKey(pub))
	return nodeid
}

func PubKeyFromHex(nodeid string) (keypair.PublicKey, error) {
	pubKey, err := hex.DecodeString(nodeid)
	if err != nil {
		return nil, err
	}
	pk, err := keypair.DeserializePublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("deserialize failed: %s", err)
	}
	return pk, err
}
