package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/royhq/go-play-app/use_cases/users/create"
)

func TestPgUsersRepository_Insert(t *testing.T) {
	t.Run("insert user successfully", func(t *testing.T) {
		// GIVEN
		mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
		require.NoError(t, err)
		t.Cleanup(mock.Close)

		repo := create.NewPgUsersRepository(mock, "test_table")
		date, _ := time.Parse(time.RFC3339, "2024-09-12T16:23:53Z")

		mock.ExpectExec(`INSERT INTO test_table (id, "name", age, created_at) 
									 VALUES($1, $2, $3, $4)`).
			WithArgs("9693cb98-f309-440b-93cb-98f309240bf3", "John Doe", 21, date).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// WHEN
		user := create.User{
			ID:      "9693cb98-f309-440b-93cb-98f309240bf3",
			Name:    "John Doe",
			Age:     21,
			Created: date,
		}

		err = repo.Insert(context.Background(), user)

		// THEN
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insert fails should return error", func(t *testing.T) {
		// GIVEN
		mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
		require.NoError(t, err)
		t.Cleanup(mock.Close)

		repo := create.NewPgUsersRepository(mock, "test_table")
		insertErr := errors.New("test error")

		mock.ExpectExec(`INSERT INTO test_table (id, "name", age, created_at) 
									 VALUES($1, $2, $3, $4)`).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnError(insertErr)

		// WHEN
		err = repo.Insert(context.Background(), create.User{})

		// THEN
		assert.EqualError(t, err, "insert user error: test error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
