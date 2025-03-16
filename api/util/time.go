package util

import (
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(timeFormat, timeStr)
}

func FormatTime(time time.Time) string {
	return time.Format(timeFormat)
}

func Now() time.Time {
	nowTime, err := ParseTime(time.Now().Format(timeFormat))
	if err != nil {
		return time.Time{}
	}

	return nowTime
}
