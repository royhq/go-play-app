package create

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreatedUserEvent struct {
	Date   time.Time
	UserID string
}

type RabbitEventPublisher struct {
	log *slog.Logger
	ch  *amqp.Channel
}

func (p *RabbitEventPublisher) Publish(ctx context.Context, event CreatedUserEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		p.log.ErrorContext(ctx, "marshal event error", slog.Any("error", err))
		return
	}

	err = p.ch.PublishWithContext(ctx,
		"",
		"users-created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		p.log.ErrorContext(ctx, "publish event error", slog.Any("error", err))
		return
	}

	p.log.DebugContext(ctx, "event published", slog.Any("event", event))

}

func NewRabbitEventPublisher(ch *amqp.Channel, log *slog.Logger) *RabbitEventPublisher {
	return &RabbitEventPublisher{
		log: log,
		ch:  ch,
	}
}
