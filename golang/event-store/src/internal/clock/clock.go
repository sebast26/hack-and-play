package clock

import "time"

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type FixedClock struct {
	time time.Time
}

func NewFixedClock(time time.Time) FixedClock {
	return FixedClock{
		time: time,
	}
}

func (f FixedClock) Now() time.Time                         { return f.time }
func (f FixedClock) After(d time.Duration) <-chan time.Time { return f.After(d) }
