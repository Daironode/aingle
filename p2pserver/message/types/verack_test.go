
package types

import (
	"testing"
)

func TestVerackSerializationDeserialization(t *testing.T) {
	var msg VerACK
	msg.isConsensus = false

	MessageTest(t, &msg)
}
