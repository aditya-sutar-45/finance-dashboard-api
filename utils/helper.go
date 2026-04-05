package utils

import "time"

func GetFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func GetString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
