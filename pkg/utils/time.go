package utils

import "time"

type TimeTeller interface {
	TellTime() time.Time
}

type timeTeller struct {
}

func NewUTCTimeTeller() TimeTeller {
	return timeTeller{}
}

func (t timeTeller) TellTime() time.Time {
	return time.Now().UTC()
}

type mockTimeTeller struct {
	time time.Time
}

func NewMockTimeTeller(time time.Time) TimeTeller {
	return mockTimeTeller{time: time}
}

func (t mockTimeTeller) TellTime() time.Time {
	return t.time
}

type delegateTimeTeller struct {
	timeFn func() time.Time
}

func NewDelegateTimeTeller(timeFn func() time.Time) TimeTeller {
	return delegateTimeTeller{timeFn: timeFn}
}

func (t delegateTimeTeller) TellTime() time.Time {
	return t.timeFn()
}
