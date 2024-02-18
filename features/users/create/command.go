package create

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrValidation = errors.New("validation error")
)

type (
	UUIDGenerator func() string

	UserInserter interface {
		Insert(context.Context, User) error
	}
)

type Command struct {
	Name string
	Age  int
}

type CommandOutput struct {
	CreatedUser User
}

type CommandHandler struct {
	log          *slog.Logger
	inserter     UserInserter
	generateUUID UUIDGenerator
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) (CommandOutput, error) {
	if err := h.validate(cmd); err != nil {
		return CommandOutput{}, err
	}

	user := User{
		ID:      UserID(h.generateUUID()),
		Name:    cmd.Name,
		Age:     cmd.Age,
		Created: time.Now(),
	}

	if err := h.inserter.Insert(ctx, user); err != nil {
		return CommandOutput{}, fmt.Errorf("create user error: %w", err)
	}

	h.log.InfoContext(ctx, "user inserted successfully")

	// TODO: add logic

	h.log.InfoContext(ctx, "user created successfully")

	return CommandOutput{CreatedUser: user}, nil
}

func (h *CommandHandler) validate(cmd Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("%w: username could not be empty", ErrValidation)
	}

	if cmd.Age <= 0 {
		return fmt.Errorf("%w: age should be greater than 0", ErrValidation)
	}

	return nil
}

func NewCommandHandler(log *slog.Logger, inserter UserInserter, uuidGen UUIDGenerator) *CommandHandler {
	return &CommandHandler{
		log:          log,
		inserter:     inserter,
		generateUUID: uuidGen,
	}
}
