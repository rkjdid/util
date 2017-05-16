package util

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Duration time.Duration

const durationBinaryVersion byte = 1

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Duration) MarshalBinary() ([]byte, error) {
	enc := []byte{
		durationBinaryVersion, // byte 0 : version
		byte(d >> 56),         // bytes 1-8: nanoseconds
		byte(d >> 48),
		byte(d >> 40),
		byte(d >> 32),
		byte(d >> 24),
		byte(d >> 16),
		byte(d >> 8),
		byte(d),
	}

	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (d *Duration) UnmarshalBinary(data []byte) error {
	buf := data
	if len(buf) == 0 {
		return errors.New("Duration.UnmarshalBinary: no data")
	}

	if buf[0] != durationBinaryVersion {
		return errors.New("Duration.UnmarshalBinary: unsupported version")
	}

	if len(buf) != /*version*/ 1+ /*nanoseconds*/ 8 {
		return errors.New("Duration.UnmarshalBinary: invalid length")
	}

	buf = buf[1:]
	*d = Duration(int64(buf[7]) |
		int64(buf[6])<<8 |
		int64(buf[5])<<16 |
		int64(buf[4])<<24 |
		int64(buf[3])<<32 |
		int64(buf[2])<<40 |
		int64(buf[1])<<48 |
		int64(buf[0])<<56)
	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (d Duration) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface.
func (d *Duration) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

func (d Duration) String() string {
	return fmt.Sprintf("\"%s\"", time.Duration(d))
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return d.MarshalText()
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	dataLength := len(data)
	if data[0] != '"' || data[dataLength-1] != '"' {
		return errors.New("Duration.UnmarshalJSON: Invalid JSON provided")
	}
	return d.UnmarshalText(data[1 : dataLength-1])
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalText(b []byte) (err error) {
	s := string(b)
	timeD, err := time.ParseDuration(s)
	if err == nil {
		*d = Duration(timeD)
		return nil
	}
	i, err2 := strconv.Atoi(s)
	if err2 == nil {
		*d = Duration(i)
		return nil
	}
	return fmt.Errorf("ParseDuration: %s; Atoi: %s", err, err2)
}
