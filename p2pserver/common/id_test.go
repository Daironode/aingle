
package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertPeerID(t *testing.T) {
	start := time.Now().Unix()
	fmt.Println("start:", start)
	RandPeerKeyId()

	end := time.Now().Unix()
	fmt.Println("end:", end)
	fmt.Println(end - start)
}

func TestKIdToUint64(t *testing.T) {
	for i := 0; i < 100; i++ {
		data := rand.Uint64()
		id := PseudoPeerIdFromUint64(data)
		data2 := id.ToUint64()
		assert.Equal(t, data, data2)
	}
}

func TestKadId_IsEmpty(t *testing.T) {
	id := PeerId{}
	assert.True(t, id.IsEmpty())
	kid := RandPeerKeyId()
	assert.False(t, kid.Id.IsEmpty())
}
