package time

import "time"

type Time struct{}

type TimeInterface interface {
	Now() time.Time
}

func New() *Time {
	return &Time{}
}

func (Time) Now() time.Time {
	return time.Now()
}
