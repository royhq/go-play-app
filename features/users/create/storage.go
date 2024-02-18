package create

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/jackc/pgx/v5"
)

type PgUsersRepository struct {
	tableName string
	db        *pgxpool.Pool
}

func (r *PgUsersRepository) Insert(ctx context.Context, user User) error {
	qFmt := `INSERT INTO %s (id, "name", age, created_at) 
			 VALUES($1, $2, $3, $4)`

	q := fmt.Sprintf(qFmt, r.tableName)

	_, err := r.db.Exec(ctx, q, user.ID, user.Name, user.Age, user.Created)
	if err != nil {
		return fmt.Errorf("insert user error: %w", err)
	}

	return nil
}

func NewPgUsersRepository(db *pgxpool.Pool, table string) *PgUsersRepository {
	return &PgUsersRepository{
		tableName: table,
		db:        db,
	}
}
