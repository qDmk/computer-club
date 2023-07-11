package utils

import "time"

const timeLayout = "15:04"

var Midnight, _ = ParseTime("00:00")

func ParseTime(s string) (time.Time, error) {
	return time.Parse(timeLayout, s)
}

func FormatTime(t time.Time) string {
	return t.Format(timeLayout)
}

func FormatDuration(d time.Duration) string {
	return FormatTime(Midnight.Add(d))
}
