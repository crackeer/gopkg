package util

import (
	"time"
)

// GetDateByUnix ...
func GetDateByUnix(unix int64) string {
	if unix == 0 {
		return ""
	}
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}

// GetStandardNowTime ...
func GetStandardNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FormatRFC3339Time(rfcTime string) string {
	if t, err := time.Parse(time.RFC3339, rfcTime); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return rfcTime
}

// GetTodayDate
//  @return string
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// GetTomorrowDate
//  @return string
func GetTomorrowDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}
