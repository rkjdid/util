package util

import "time"

func ParisTime(t time.Time) time.Time {
	return LocalizedTime(t, "Europe/Paris")
}

func LocalizedTime(t time.Time, s string) time.Time {
	loc, err := time.LoadLocation(s)
	if err != nil {
		return t
	}
	return t.In(loc)
}
