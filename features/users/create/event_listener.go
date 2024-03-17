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

func (l *RabbitEventListener) Listen(ctx context.Context) error {
	eventsCh, err := l.ch.ConsumeWithContext(ctx,
		l.queue,
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

	errCh := make(chan error)

	go func() {
		l.log.InfoContext(ctx, "listening...")

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case e := <-eventsCh:
				l.processEvent(ctx, e)
			}
		}
	}()

	err = <-errCh

	l.log.InfoContext(ctx, "listener stopped", "error", err)

	return err
}

func (l *RabbitEventListener) processEvent(ctx context.Context, e amqp.Delivery) {
	l.log.DebugContext(ctx, "event received", slog.String("body", string(e.Body)))

	var evt CreatedUserEvent

	if evtErr := json.Unmarshal(e.Body, &evt); evtErr != nil {
		l.log.ErrorContext(ctx, "unmarshal event error", "error", evtErr)
		return
	}

	if listenErr := l.eventHandler.Handle(ctx, evt); listenErr != nil {
		l.log.ErrorContext(ctx, "error listening event",
			slog.Any("error", listenErr),
			slog.Any("event", evt))
	}
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
		log:          log.With(slog.String("queue", queue)),
	}
}
