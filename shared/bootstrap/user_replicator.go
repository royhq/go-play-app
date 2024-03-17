package bootstrap

import (
	"context"

	"github.com/royhq/go-play-app/features/users/create"
)

type UserReplicatorApp struct {
	onShutdown []func()

	eventListener *create.RabbitEventListener
}

func (a *UserReplicatorApp) Start() error {
	ctx := context.Background()
	return a.eventListener.Listen(ctx)
}

func (a *UserReplicatorApp) Shutdown() {
	execute(a.onShutdown...)
}

func NewUserReplicator() (*UserReplicatorApp, error) {
	logger := defaultLogger()
	onShutdown := make([]func(), 0, 1)

	usersCh, err := usersCreatedChannel()
	if err != nil {
		usersCh.Close()
		return nil, err
	}

	onShutdown = append(onShutdown, usersCh.Close)

	replicator := create.NewUserReplicator(nil, nil, logger)
	eventListener := create.NewRabbitEventListener(usersCh.Queue.Name, usersCh.Channel, replicator, logger)

	app := &UserReplicatorApp{
		onShutdown:    onShutdown,
		eventListener: eventListener,
	}

	return app, nil
}
