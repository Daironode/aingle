
package vbft

import "testing"

func constructEventTimer() *EventTimer {
	server := constructServer()
	return NewEventTimer(server)
}

func TestStartTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	eventtimer.StartTimer(1, 10)
}

func TestCancelTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	eventtimer.StartTimer(1, 10)
	eventtimer.CancelTimer(1)
}
func TestStartEventTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.startEventTimer(EventProposeBlockTimeout, 1)
	t.Logf("TestStartEventTimer: %v", err)
}

func TestCancelEventTimer(t *testing.T) {
	eventtimer := constructEventTimer()
	err := eventtimer.startEventTimer(EventProposeBlockTimeout, 1)
	t.Logf("startEventTimer: %v", err)
	eventtimer.cancelEventTimer(EventProposeBlockTimeout, 1)
}
