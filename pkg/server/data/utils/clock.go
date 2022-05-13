package utils

import "time"

type Clock struct{}

func (c *Clock) Now() int64 {
	return time.Now().Unix()
}
