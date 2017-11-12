package util

import (
	"log"
	"sync"
	"time"
)

const chanSize = 36

const panicNegativeLength = "timeseries length must be positive"

// TimeSeries is a naive implementation of int timeseries.
// It holds a slice of ints between Start and End times with a maximum length.
// Each int is supposed to be separated by Interval time, but this is obviously dependant on user.
// maxLength is a positive (>=0) amount, if equals to 0, there is no maximum length for TimeSeries.
type TimeSeries struct {
	Start    time.Time
	End      time.Time
	Interval Duration
	Data     []int

	maxLength   int
	subscribers map[int]chan int
	currentId   int
	sync.RWMutex
}

// NewTimeSeries initiates a TimeSeries.
// A negative length panics, 0 means no limit.
func NewTimeSeries(length int, interval Duration) *TimeSeries {
	if length < 0 {
		panic(panicNegativeLength)
	}
	return &TimeSeries{
		Start:    time.Now(),
		End:      time.Now(),
		Interval: interval,
		Data:     make([]int, 0, length),

		maxLength:   length,
		subscribers: make(map[int]chan int),
	}
}

// Subscribe creates a new buffered chan, adds it to subscribers map for
// later broadcast on Add() calls. The id of the chan is returned for
// unsubscribing, and the chan itself.
func (d *TimeSeries) Subscribe() (id int, data chan int) {
	if d.subscribers == nil {
		d.subscribers = make(map[int]chan int)
	}

	data = make(chan int, chanSize)
	d.Lock()
	id = d.currentId
	d.subscribers[id] = data
	d.currentId++
	d.Unlock()
	return id, data
}

// Unsubscribe closes and deletes from subscribers channel with id.
func (d *TimeSeries) Unsubscribe(id int) {
	d.Lock()
	ch, ok := d.subscribers[id]
	if ok {
		close(ch)
		delete(d.subscribers, id)
	}
	d.Unlock()
}

// Add appends value to d.Data. If d.Data length is above d.maxLength, do a 1-shift
// operation on the slice to the left, and a naive shift of +d.Interval on d.Start (data loss).
// After appending value, it is broadcasted to subscribed chans.
func (d *TimeSeries) Add(v int) {
	if d.maxLength > 0 && len(d.Data) >= d.maxLength {
		d.Data = d.Data[1:]
		// trusting time shift
		// assumes d.Interval is somehow respected at each Add call
		d.Start = d.Start.Add(time.Duration(d.Interval))
	}
	d.Data = append(d.Data, v)
	d.End = time.Now()
	d.broadcast(v)
}

// Padded returns d.Data left-padded up to d.maxLength with d.Data[0] values.
func (d *TimeSeries) Padded() []int {
	if len(d.Data) >= d.maxLength {
		return d.Data[:d.maxLength]
	}

	var zero int
	if len(d.Data) > 0 {
		zero = d.Data[0]
	}

	data := make([]int, d.maxLength)
	for i := 0; i < d.maxLength-len(d.Data); i++ {
		data[i] = zero
	}
	copy(data[d.maxLength-len(d.Data):], d.Data)
	return data
}

// SetMaxLength resets d.maxLength.
// If necessary, oldest data are cropped to match new maximum length.
func (d *TimeSeries) SetMaxLength(length int) {
	if length < 0 {
		panic(panicNegativeLength)
	}
	d.maxLength = length
	if len(d.Data) > d.maxLength {
		stripN := len(d.Data) - d.maxLength
		d.Start = d.Start.Add(time.Duration(stripN * int(d.Interval)))
		d.Data = d.Data[stripN:]
	}
}

// ResetStartTime operates a timeshift on Start time from time.Now()
// of -len(d.Data) * d.Interval
func (d *TimeSeries) ResetStartTime() {
	d.Start = time.Now().Add(-time.Duration(int64(len(d.Data)) * int64(d.Interval)))
}

func (d *TimeSeries) broadcast(v int) {
	d.Lock()
	for i, ch := range d.subscribers {
		if len(ch) == cap(ch) {
			log.Printf("timeseries: unsubscribing full channel %d", i)
			delete(d.subscribers, i)
		} else {
			ch <- v
		}
	}
	d.Unlock()
}
