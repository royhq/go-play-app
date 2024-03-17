package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreatePoolInput struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Ping     bool
}

func CreatePool(input *CreatePoolInput) (*pgxpool.Pool, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connectionString(input))
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	if input.Ping {
		if err = conn.Ping(ctx); err != nil {
			return nil, fmt.Errorf("ping db error: %w", err)
		}

		slog.InfoContext(ctx, "ping db success")
	}

	return conn, nil
}

func connectionString(i *CreatePoolInput) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		i.Host, i.Port, i.User, i.Password, i.Database)
}
