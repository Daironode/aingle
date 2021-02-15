
package subnet

import (
	"github.com/blang/semver"
)

const MIN_VERSION_FOR_SUBNET = "2.0.0-0"

func supportSubnet(version string) bool {
	if version == "" {
		return false
	}
	v1, err := semver.ParseTolerant(version)
	if err != nil {
		return false
	}
	min, err := semver.ParseTolerant(MIN_VERSION_FOR_SUBNET)
	if err != nil {
		panic(err) // enforced by testcase
	}

	return v1.GTE(min)
}
