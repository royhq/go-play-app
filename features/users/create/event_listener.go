package create

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreatedUserEventHandler interface {
	Handle(context.Context, CreatedUserEvent) error
}

type RabbitEventListener struct {
	queue        string
	ch           *amqp.Channel
	eventHandler CreatedUserEventHandler
	log          *slog.Logger
}

func (r *RabbitEventListener) Listen(ctx context.Context) error {
	eventsCh, err := r.ch.ConsumeWithContext(ctx,
		r.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("rabbitmq consume error: %w", err)
	}

	log := r.log.With(slog.String("queue", r.queue))

	var forever chan struct{}

	go func() {
		log.InfoContext(ctx, "listening...")

		for e := range eventsCh {
			log.DebugContext(ctx, "event received", slog.String("body", string(e.Body)))

			var evt CreatedUserEvent
			if evtErr := json.Unmarshal(e.Body, &evt); evtErr != nil {
				log.ErrorContext(ctx, "unmarshal event error", "error", evtErr)
				return
			}

			if listenErr := r.eventHandler.Handle(ctx, evt); listenErr != nil {
				log.ErrorContext(ctx, "error listening event",
					slog.Any("error", listenErr),
					slog.Any("event", evt))
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.InfoContext(ctx, "ctx done", "error", ctx.Err())
		return ctx.Err()
	case <-forever:
	}

	return nil
}

func NewRabbitEventListener(
	queue string,
	ch *amqp.Channel,
	handler CreatedUserEventHandler,
	log *slog.Logger,
) *RabbitEventListener {
	return &RabbitEventListener{
		queue:        queue,
		ch:           ch,
		eventHandler: handler,
		log:          log,
	}
}
