package bootstrap

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-play-app/features/ping"
	userscreate "go-play-app/features/users/create"
	"go-play-app/infra/http/middleware"
	"go-play-app/infra/uuid"
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

	createUsersRepo := userscreate.NewPgUsersRepository(db, "users")
	createUserCmdHandler := userscreate.NewCommandHandler(logger, createUsersRepo, uuid.New)

	app := &MainApp{
		PingHandler:       ping.NewEndpointHandler(),
		CreateUserHandler: userscreate.NewEndpointHandler(createUserCmdHandler.Handle),
	}

	app.onShutdown = append(app.onShutdown, func() {
		log.Println("closing db connection...")
		db.Close()
	})

	return app, nil
}

func defaultLogger() *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, nil)
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
