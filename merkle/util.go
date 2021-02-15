 
package merkle

import "math/bits"

// return the number of 1 bit
func countBit(num uint32) uint {
	return uint(bits.OnesCount32(num))
}

func isPower2(num uint32) bool {
	return countBit(num) == 1
}

// return the position of the heightest 1 bit
// 1-based index
func highBit(num uint32) uint {
	return uint(32 - bits.LeadingZeros32(num))
}

// return the position of the lowest 1 bit
// 1-based index
func lowBit(num uint32) uint {
	return highBit(num & -num)
}
