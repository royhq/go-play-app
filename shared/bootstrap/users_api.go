package bootstrap

import (
	"net/http"

	"github.com/royhq/go-play-app/features/ping"
	userscreate "github.com/royhq/go-play-app/features/users/create"
	"github.com/royhq/go-play-app/shared/commons/clock"
	"github.com/royhq/go-play-app/shared/infra/http/middleware"
	"github.com/royhq/go-play-app/shared/infra/uuid"
)

type UsersAPIApp struct {
	onShutdown []func()

	PingHandler       http.Handler
	CreateUserHandler http.Handler
}

func (a *UsersAPIApp) ListenAndServe(addr string) error {
	mux := http.NewServeMux()

	pingHandler := middleware.WithRecover(a.PingHandler)
	createUserHandler := middleware.WithRecover(middleware.WithRequestID(a.CreateUserHandler))

	mux.Handle("GET /ping", pingHandler)
	mux.Handle("POST /users", createUserHandler)

	return http.ListenAndServe(addr, mux)
}

func (a *UsersAPIApp) Shutdown() {
	execute(a.onShutdown...)
}

func NewUsersAPI() (*UsersAPIApp, error) {
	logger := defaultLogger()
	onShutdown := make([]func(), 0, 2) //nolint:gomnd // max funcs

	db, err := dbConnection()
	if err != nil {
		return nil, err
	}

	onShutdown = append(onShutdown, db.Close)

	usersCh, err := usersCreatedChannel()
	if err != nil {
		execute(onShutdown...)
		usersCh.Close()

		return nil, err
	}

	onShutdown = append(onShutdown, usersCh.Close)

	// create user
	createUsersRepo := userscreate.NewPgUsersRepository(db, "users")
	createUserCmdHandler := userscreate.NewCommandHandler(
		logger,
		clock.Default(),
		createUsersRepo,
		userscreate.NewRabbitEventPublisher(usersCh.Queue.Name, usersCh.Channel, logger),
		uuid.New,
	)
	createUserEndpointHandler := userscreate.NewEndpointHandler(
		createUserCmdHandler.Handle, userscreate.NewEndpointErrorHandler(logger),
	)

	app := &UsersAPIApp{
		onShutdown:        onShutdown,
		PingHandler:       ping.NewEndpointHandler(),
		CreateUserHandler: createUserEndpointHandler,
	}

	return app, nil
}
