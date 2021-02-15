
package utils

import (
	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/common/constants"
)

var (
	TIME_INTERVAL         = constants.UNBOUND_TIME_INTERVAL
	GENERATION_AMOUNT     = constants.UNBOUND_GENERATION_AMOUNT
	NEW_GENERATION_AMOUNT = constants.NEW_UNBOUND_GENERATION_AMOUNT
)

// startOffset : start timestamp offset from genesis block
// endOffset :  end timestamp offset from genesis block
func CalcUnbindOng(balance uint64, startOffset, endOffset uint32) uint64 {
	var amount uint64 = 0
	if startOffset >= endOffset {
		return 0
	}
	if startOffset < config.GetOntHolderUnboundDeadline() {
		ustart := startOffset / TIME_INTERVAL
		istart := startOffset % TIME_INTERVAL
		if endOffset >= config.GetOntHolderUnboundDeadline() {
			endOffset = config.GetOntHolderUnboundDeadline()
		}
		uend := endOffset / TIME_INTERVAL
		iend := endOffset % TIME_INTERVAL
		for ustart < uend {
			amount += uint64(TIME_INTERVAL-istart) * GENERATION_AMOUNT[ustart]
			ustart++
			istart = 0
		}
		amount += uint64(iend-istart) * GENERATION_AMOUNT[ustart]
	}

	return uint64(amount) * balance
}

// startOffset : start timestamp offset from genesis block
// endOffset :  end timestamp offset from genesis block
func CalcGovernanceUnbindOng(startOffset, endOffset uint32) uint64 {
	if endOffset < config.GetOntHolderUnboundDeadline() {
		return 0
	}
	if startOffset < config.GetOntHolderUnboundDeadline() {
		startOffset = config.GetOntHolderUnboundDeadline()
	}

	var amount uint64 = 0
	if startOffset >= endOffset {
		return 0
	}
	deadline, _ := config.GetGovUnboundDeadline()
	var gap uint64 = 0
	if startOffset < deadline {
		ustart := startOffset / TIME_INTERVAL
		istart := startOffset % TIME_INTERVAL
		if endOffset > deadline {
			endOffset = deadline
			_, gap = config.GetGovUnboundDeadline()
		}
		uend := endOffset / TIME_INTERVAL
		iend := endOffset % TIME_INTERVAL
		for ustart < uend {
			amount += uint64(TIME_INTERVAL-istart) * NEW_GENERATION_AMOUNT[ustart]
			ustart++
			istart = 0
		}
		amount += uint64(iend-istart) * NEW_GENERATION_AMOUNT[ustart]
		amount += gap
	}

	return uint64(amount) * constants.ONT_TOTAL_SUPPLY
}
