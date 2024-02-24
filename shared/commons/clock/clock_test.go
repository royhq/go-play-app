package clock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/shared/commons/clock"
)

func TestTimeFunc_Now(t *testing.T) {
	t.Parallel()

	t.Run("should return given datetime", func(t *testing.T) {
		// GIVEN
		aDate := time.Date(2024, 9, 10, 14, 30, 15, 0, time.UTC)
		timeFn := clock.TimeFunc(func() time.Time {
			return aDate
		})

		// WHEN
		funcValue := timeFn()
		nowValue := timeFn.Now()

		// THEN
		assert.Equal(t, funcValue, nowValue)
		assert.Equal(t, aDate, funcValue)
	})
}

func TestAt(t *testing.T) {
	t.Parallel()

	t.Run("should return a clock with the given date", func(t *testing.T) {
		// GIVEN
		aDate := time.Date(2024, 9, 10, 14, 30, 15, 0, time.UTC)

		// WHEN
		c := clock.At(aDate)

		// THEN
		assert.Equal(t, aDate, c.Now())
	})
}

func TestMustParseAt(t *testing.T) {
	t.Parallel()

	t.Run("should return a clock with the given date", func(t *testing.T) {
		// WHEN
		c := clock.MustParseAt(time.RFC3339, "2024-09-10T14:30:15.0Z")

		// THEN
		expected := time.Date(2024, 9, 10, 14, 30, 15, 0, time.UTC)
		assert.Equal(t, expected, c.Now())
	})

	t.Run("should panic when cannot parse date", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = clock.MustParseAt(time.RFC3339, "invalid date")
		})
	})
}

func TestDefault(t *testing.T) {
	c := clock.Default()
	assert.NotNil(t, c)
	assert.NotZero(t, c.Now())
}
