package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GetChannelInput struct {
	URL       string
	QueueName string
}

type GetChannelOutput struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func (o *GetChannelOutput) Close() {
	ctx := context.Background()

	log := slog.Default().With(
		slog.String("queue", o.Queue.Name),
	)

	log.DebugContext(ctx, "closing rabbitmq resources...",
		slog.String("queue", o.Queue.Name))

	// close channel
	if o.Channel != nil {
		if err := o.Channel.Close(); err != nil {
			log.ErrorContext(ctx, "rabbitmq closing channel error",
				slog.Any("error", err))
		}
	}

	// close connection
	if o.Connection != nil {
		if err := o.Connection.Close(); err != nil {
			log.ErrorContext(ctx, "rabbitmq closing connection error",
				slog.Any("error", err))
		}
	}

	log.DebugContext(ctx, "rabbitmq resources closed")
}

func GetChannel(input GetChannelInput) (GetChannelOutput, error) {
	conn, err := amqp.Dial(input.URL)
	if err != nil {
		return GetChannelOutput{}, fmt.Errorf("rabbitmq connection error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return GetChannelOutput{}, fmt.Errorf("rabbitmq channel error: %w", err)
	}

	queue, err := ch.QueueDeclare(
		input.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return GetChannelOutput{}, fmt.Errorf("rabbit queue declare error: %w", err)
	}

	return GetChannelOutput{
		Connection: conn,
		Channel:    ch,
		Queue:      queue,
	}, nil
}
