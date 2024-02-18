package clock

import "time"

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (f ClockFunc) Now() time.Time {
	return f()
}

func At(v time.Time) Clock {
	return ClockFunc(func() time.Time {
		return v
	})
}

func MustParseAt(layout, value string) Clock {
	parsed, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return At(parsed)
}

func Default() Clock {
	return ClockFunc(time.Now)
}
