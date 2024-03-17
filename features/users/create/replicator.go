package create

import (
	"context"
	"log/slog"

	"github.com/royhq/go-play-app/shared/domain"
)

type UserGetter interface {
	Get(context.Context, domain.UserID) (User, error)
}

type UsersReplica interface {
	Set(context.Context, User) error
}

type UserReplicator struct {
	users   UserGetter
	replica UsersReplica
	log     *slog.Logger
}

func (r *UserReplicator) Handle(ctx context.Context, event CreatedUserEvent) error {
	r.log.DebugContext(ctx, "replicator start handling event...",
		slog.Any("event", event))

	// TODO: implement

	return nil
}

func NewUserReplicator(users UserGetter, replica UsersReplica, log *slog.Logger) *UserReplicator {
	return &UserReplicator{log: log}
}
