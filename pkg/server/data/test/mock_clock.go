package test

import (
	"time"

	"github.com/nrc-no/core/pkg/server/data/api"
)

// MockClock is a mock clock that can be used for testing.
type MockClock struct {
	clock   api.Clock
	TheTime int64
}

func (c *MockClock) Now() int64 {
	if c.clock != nil {
		return c.clock.Now()
	}
	return c.TheTime
}

func (c *MockClock) Tick(d time.Duration) {
	c.TheTime += int64(d)
}

func (c *MockClock) UseClock(clock api.Clock) func() {
	oldClock := c.clock
	c.clock = clock
	return func() {
		c.clock = oldClock
	}
}
