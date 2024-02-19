package uuid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/infra/uuid"
)

func TestNew(t *testing.T) {
	newUUID := uuid.New()
	assert.True(t, uuid.IsUUID(newUUID))
}

func TestIsUUID(t *testing.T) {
	testCases := map[string]struct {
		givenUUID string
		expected  bool
	}{
		"valid uuid should return true": {
			givenUUID: "d1c5e07a-92eb-460e-b3f8-219bb1082278",
			expected:  true,
		},
		"invalid uuid should return false": {
			givenUUID: "xxx",
			expected:  false,
		},
		"empty string should return false": {
			givenUUID: "",
			expected:  false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, uuid.IsUUID(tc.givenUUID))
		})
	}
}
