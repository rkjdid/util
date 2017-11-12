package util

import (
	"fmt"
	"strconv"
)

type Float float64

func (f Float) String() string {
	return fmt.Sprint(float64(f))
}

func (f Float) MarshalJSON() ([]byte, error) {
	b, err := f.MarshalText()
	if err == nil {
		b = []byte(fmt.Sprintf("%s", string(b)))
	}
	return b, err
}

func (f *Float) UnmarshalJSON(data []byte) error {
	return f.UnmarshalText(data)
}

func (f Float) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Float) UnmarshalText(b []byte) (err error) {
	s := string(b)
	fl, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*f = Float(fl)
	} else if err != nil {
		i, err2 := strconv.Atoi(s)
		if err2 != nil {
			return fmt.Errorf("Float.UnmarshalText: %s", err)
		}
		*f = Float(i)
	}
	return nil
}
