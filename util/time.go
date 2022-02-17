package util

import "time"

func ParseTimeInLocation(timeZone, layout, timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation("2006-01-02T15:04:05", timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
