package clock

import "time"

type Clock interface {
	Now() time.Time
}

type TimeFunc func() time.Time

func (f TimeFunc) Now() time.Time {
	return f()
}

func At(v time.Time) TimeFunc {
	return func() time.Time { return v }
}

func MustParseAt(layout, value string) TimeFunc {
	parsed, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return At(parsed)
}

func Default() TimeFunc { return time.Now }
