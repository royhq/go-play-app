package create

import (
	"context"
	"errors"
	"log/slog"

	"github.com/royhq/go-play-app/shared/commons/clock"
	"github.com/royhq/go-play-app/shared/domain"
)

type (
	UUIDGenerator func() string

	UserInserter interface {
		Insert(context.Context, User) error
	}

	UserCreatedEventPublisher interface {
		Publish(context.Context, CreatedUserEvent)
	}
)

type Command struct {
	Name string
	Age  int
}

type CommandOutput struct {
	CreatedUser User
}

type CommandError struct {
	Msg   string
	Code  string
	Cause error
}

func (e *CommandError) Error() string { return e.Msg }
func (e *CommandError) Unwrap() error { return e.Cause }

type CommandHandler struct {
	log          *slog.Logger
	clock        clock.Clock
	inserter     UserInserter
	publisher    UserCreatedEventPublisher
	generateUUID UUIDGenerator
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) (CommandOutput, error) {
	if err := h.validate(cmd); err != nil {
		return CommandOutput{}, NewValidationError(err.Error())
	}

	user := User{
		ID:      domain.UserID(h.generateUUID()),
		Name:    cmd.Name,
		Age:     cmd.Age,
		Created: h.clock.Now(),
	}

	if err := h.inserter.Insert(ctx, user); err != nil {
		return CommandOutput{}, &CommandError{Msg: "create user error", Code: "users_error", Cause: err}
	}

	h.log.DebugContext(ctx, "user inserted successfully")

	event := CreatedUserEvent{Date: h.clock.Now(), UserID: string(user.ID)}

	h.publisher.Publish(ctx, event)

	h.log.InfoContext(ctx, "user created successfully")

	return CommandOutput{CreatedUser: user}, nil
}

func (h *CommandHandler) validate(cmd Command) error {
	if cmd.Name == "" {
		return errors.New("name could not be empty")
	}

	if cmd.Age <= 0 {
		return errors.New("age should be greater than zero")
	}

	return nil
}

func NewCommandHandler(
	log *slog.Logger,
	clock clock.Clock,
	inserter UserInserter,
	publisher UserCreatedEventPublisher,
	uuidGen UUIDGenerator,
) *CommandHandler {
	return &CommandHandler{
		log:          log,
		clock:        clock,
		inserter:     inserter,
		publisher:    publisher,
		generateUUID: uuidGen,
	}
}
