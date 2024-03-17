package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, closeCh, err := usersCreatedChannel("users-created")
	if err != nil {
		closeCh()
		log.Fatal(err)
	}

	c, err := ch.Consume("users-created", "", true, false, false, false, nil)
	if err != nil {
		log.Println("error consuming:", err)
	}

	var forever chan struct{}

	go func() {
		log.Println("listening...")
		for d := range c {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
	/*defer func() {
		closeCh()
	}()*/
}

func usersCreatedChannel(queueName string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, func() {}, fmt.Errorf("rabbitmq connection error: %w", err)
	}

	closeConn := func() {
		if connErr := conn.Close(); connErr != nil {
			slog.ErrorContext(context.Background(), "error closing rabbitmq connection",
				"error", connErr)
		}
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, closeConn, fmt.Errorf("rabbitmq channel error: %w", err)
	}

	closeCh := func() {
		if chErr := ch.Close(); chErr != nil {
			slog.ErrorContext(context.Background(), "error closing rabbitmq channel",
				"error", chErr)
		}

		closeConn()
	}

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return ch, closeConn, fmt.Errorf("rabbit queue declare error: %w", err)
	}

	return ch, closeCh, nil
}
