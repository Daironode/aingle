 
package neovm

import "testing"

// testcase enforce keys in ServiceMap and ServiceMapDeprecated
func TestNeoVmServiceMap(t *testing.T) {
	for k := range ServiceMap {
		if ServiceMapDeprecated[k] != nil || ServiceMapNew[k] != nil {
			panic("key in ServiceMap also in ServiceMapDeprecated or ServiceMapNew")
		}
	}
}
