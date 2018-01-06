package util

import "time"

func ParisTime(t time.Time) time.Time {
	return LocalizedTime(t, "Europe/Paris")
}

func LocalizedTime(t time.Time, loc string) time.Time {
	loc, err := time.LoadLocation(loc)
	if err != nil {
		return t
	}
	return t.In(loc)
}
