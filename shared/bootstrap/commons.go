package bootstrap

import (
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/royhq/go-play-app/shared/infra/postgres"
	"github.com/royhq/go-play-app/shared/infra/rabbitmq"
)

func defaultLogger() *slog.Logger {
	// TODO: maybe set logger level and options by LOG_LEVEL env var or similar.
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(h)

	slog.SetDefault(logger)

	return logger
}

func dbConnection() (*pgxpool.Pool, error) {
	// TODO: this should be obtained from env var or similar.
	return postgres.CreatePool(postgres.CreatePoolInput{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Ping:     true,
	})
}

func usersCreatedChannel() (rabbitmq.GetChannelOutput, error) {
	// TODO: this should be obtained from env var or similar.
	return rabbitmq.GetChannel(rabbitmq.GetChannelInput{
		URL:       "amqp://guest:guest@localhost:5672/",
		QueueName: "users-created",
	})
}

func execute(f ...func()) {
	for _, fn := range f {
		fn()
	}
}
