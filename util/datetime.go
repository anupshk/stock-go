package util

import (
	"time"
)

const (
	DATE_FORMAT     = "2006-01-02"
	DATETIME_FORMAT = time.RFC3339
)

func GetCurrentTime() time.Time {
	return time.Now().In(TZLocation)
}

func GetFormattedCurrentTime(format string) string {
	return GetCurrentTime().Format(format)
}
