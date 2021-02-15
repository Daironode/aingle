
package subnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompatibility(t *testing.T) {
	unsupported := []string{"v1.10.0", "1.10.0-alpha", "v1.9"}
	for _, version := range unsupported {
		assert.False(t, supportSubnet(version))
	}

	supported := []string{"v2.0.0-0-gfcbf82c", "v2.0.0", "2.0.0-alpha", "2.0.0-beta", "v2.0.0-alpha.9", "v2.0.0-laizy", "v2.0.0-laizy1"}
	for _, version := range supported {
		assert.True(t, supportSubnet(version))
	}
}
