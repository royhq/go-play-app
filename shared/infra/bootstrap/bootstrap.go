package bootstrap

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/royhq/go-play-app/features/ping"
	userscreate "github.com/royhq/go-play-app/features/users/create"
	"github.com/royhq/go-play-app/shared/commons/clock"
	"github.com/royhq/go-play-app/shared/infra/http/middleware"
	"github.com/royhq/go-play-app/shared/infra/uuid"
)

type MainApp struct {
	onShutdown []func()

	PingHandler       http.Handler
	CreateUserHandler http.Handler
}

func (a *MainApp) ListenAndServe(addr string) error {
	mux := http.NewServeMux()

	pingHandler := middleware.WithRecover(a.PingHandler)
	createUserHandler := middleware.WithRecover(middleware.WithRequestID(a.CreateUserHandler))

	// routing
	mux.Handle("GET /ping", pingHandler)
	mux.Handle("POST /users", createUserHandler)

	return http.ListenAndServe(addr, mux)
}

func (a *MainApp) Shutdown() {
	for _, fn := range a.onShutdown {
		fn()
	}
}

func NewMainApp() (*MainApp, error) {
	logger := defaultLogger()

	db, err := dbConnection(true)
	if err != nil {
		return nil, err
	}

	ch, closeCh, err := usersCreatedChannel("users-created")

	// create user
	createUsersRepo := userscreate.NewPgUsersRepository(db, "users")
	createUserCmdHandler := userscreate.NewCommandHandler(
		logger,
		clock.Default(),
		createUsersRepo,
		userscreate.NewRabbitEventPublisher(ch, logger),
		uuid.New,
	)
	createUserEndpointHandler := userscreate.NewEndpointHandler(
		createUserCmdHandler.Handle, userscreate.NewEndpointErrorHandler(logger),
	)

	app := &MainApp{
		PingHandler:       ping.NewEndpointHandler(),
		CreateUserHandler: createUserEndpointHandler,
	}

	app.onShutdown = append(app.onShutdown, func() {
		log.Println("closing db connection...")
		db.Close()

		if closeCh != nil {
			closeCh()
		}
	})

	return app, nil
}

func defaultLogger() *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(h)

	slog.SetDefault(logger)

	return logger
}

func dbConnection(ping bool) (*pgxpool.Pool, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	if ping {
		if err = conn.Ping(ctx); err != nil {
			return nil, fmt.Errorf("ping db error: %w", err)
		}

		log.Println("ping db success")
	}

	return conn, nil
}

func usersCreatedChannel(queueName string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitmq connection error: %w", err)
	}

	closeConn := func() {
		if connErr := conn.Close(); connErr != nil {
			slog.ErrorContext(context.Background(), "error closing rabbitmq connection",
				"error", connErr)
		}
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, closeConn, fmt.Errorf("rabbitmq channel error: %w", err)
	}

	closeCh := func() {
		if chErr := ch.Close(); chErr != nil {
			slog.ErrorContext(context.Background(), "error closing rabbitmq channel",
				"error", chErr)
		}

		closeConn()
	}

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return ch, closeConn, fmt.Errorf("rabbit queue declare error: %w", err)
	}

	return ch, closeCh, nil
}
