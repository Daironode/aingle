
package utils

import "time"

type Parker struct {
	c chan struct{}
}

func NewParker() *Parker {
	c := make(chan struct{}, 0)
	return &Parker{c: c}
}

func (self *Parker) ParkTimeout(d time.Duration) {
	select {
	case <-self.c:
	case <-time.After(d):
	}
}

func (self *Parker) Unpark() {
	select {
	case self.c <- struct{}{}:
	default:
	}
}
