package create

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitEventPublisher struct {
	queue string
	ch    *amqp.Channel
	log   *slog.Logger
}

func (p *RabbitEventPublisher) Publish(ctx context.Context, event CreatedUserEvent) {
	body, err := json.Marshal(event)
	if err != nil {
		p.log.ErrorContext(ctx, "marshal event error", slog.Any("error", err))
		return
	}

	err = p.ch.PublishWithContext(ctx,
		"",
		p.queue,
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

func NewRabbitEventPublisher(queue string, ch *amqp.Channel, log *slog.Logger) *RabbitEventPublisher {
	return &RabbitEventPublisher{
		queue: queue,
		ch:    ch,
		log:   log,
	}
}
