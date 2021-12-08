package dates

import (
	"errors"
	"github.com/snabb/isoweek"
	"strconv"
	"strings"
	"time"
)

func ParseIsoWeekTime(valueString string) (time.Time, error)  {
	parts := strings.Split(valueString, "-W")
	if len(parts) != 2 {
		return time.Time{}, errors.New("unexpected week field format")
	}
	var week int
	var year int
	for i, part := range parts {
		p, err := strconv.Atoi(part)
		if err != nil {
			return time.Time{}, errors.New("invalid week field format")
		}
		if i == 0 {
			year = p
		}
		if i == 1 {
			week = p
		}
	}
	if isoweek.Validate(year, week) {
		return isoweek.StartTime(year, week, time.UTC), nil
	} else {
		return time.Time{}, errors.New("invalid iso week format")
	}

}
