package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/royhq/go-play-app/commons/clock"
	"github.com/royhq/go-play-app/features/users/create"
	"github.com/royhq/go-play-app/internal/mocks"
)

func TestCommandHandler_Handle(t *testing.T) {
	t.Run("command successfully executed", func(t *testing.T) {
		// GIVEN
		_clock := clock.MustParseAt(time.RFC3339, "2024-03-20T15:23:00Z")
		userInserter := mocks.NewUserInserterMock(t)
		uuidGenerator := fakeUUIDGeneratorReturns("5d1b769c-777d-434c-9b76-9c777d734ce1")

		cmdHandler := create.NewCommandHandler(noLogger(), _clock, userInserter, uuidGenerator)

		expectedUser := create.User{
			ID:      "5d1b769c-777d-434c-9b76-9c777d734ce1",
			Name:    "John Doe",
			Age:     21,
			Created: _clock.Now(),
		}

		userInserter.EXPECT().Insert(mock.Anything, expectedUser).Return(nil).Once()

		// WHEN
		cmd := create.Command{
			Name: "John Doe",
			Age:  21,
		}

		out, err := cmdHandler.Handle(context.Background(), cmd)

		// THEN
		assert.NoError(t, err)

		expectedOut := create.CommandOutput{CreatedUser: expectedUser}
		assert.Equal(t, expectedOut, out)
	})

	t.Run("command validation errors", testCommandValidationErrors)
}

func testCommandValidationErrors(t *testing.T) {
	testCases := map[string]struct {
		givenCmd    create.Command
		expectedErr error
	}{
		"command with no name should return error": {
			givenCmd: create.Command{Name: "", Age: 21},
			expectedErr: &create.ValidationError{
				Msg: "name could not be empty",
			},
		},
		"command with zero age should return error": {
			givenCmd: create.Command{Name: "John Doe"},
			expectedErr: &create.ValidationError{
				Msg: "age should be greater than zero",
			},
		},
		"command with negative age should return error": {
			givenCmd: create.Command{Name: "John Doe", Age: -1},
			expectedErr: &create.ValidationError{
				Msg: "age should be greater than zero",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			cmdHandler := create.NewCommandHandler(noLogger(), nil, nil, nil)

			// WHEN
			out, err := cmdHandler.Handle(context.Background(), tc.givenCmd)

			// THEN
			assert.Zero(t, out)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCommandError(t *testing.T) {
	// GIVEN
	cause := errors.New("original error")

	err := &create.CommandError{
		Msg:   "test message",
		Code:  "test_code",
		Cause: cause,
	}

	// WHEN & THEN
	assert.Equal(t, "test_code", err.Code)
	assert.Equal(t, "test message", err.Error())
	assert.Equal(t, cause, err.Unwrap())
	assert.ErrorIs(t, err, cause)
}

func fakeUUIDGeneratorReturns(uuid string) create.UUIDGenerator {
	return func() string { return uuid }
}
