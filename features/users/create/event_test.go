package create_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/royhq/go-play-app/features/users/create"
	"github.com/royhq/go-play-app/internal/mocks"
	"github.com/royhq/go-play-app/shared/infra/rabbitmq"
)

func TestRabbitEventPublisher_Publish(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Parallel()

	container := initRabbitMqContainer(t, true)

	// GIVEN
	getChInput := rabbitmq.GetChannelInput{
		URL:       fmt.Sprintf("amqp://guest:guest@%s:%s/", container.host, container.port),
		QueueName: "test-queue",
	}

	sender, err := rabbitmq.GetChannel(getChInput)
	if err != nil {
		t.Fatalf("could not get rabbitmq channel (sender): %s", err)
	}

	t.Cleanup(sender.Close)

	receiver, err := rabbitmq.GetChannel(getChInput)
	if err != nil {
		t.Fatalf("could not get rabbitmq channel (receiver): %s", err)
	}

	t.Cleanup(receiver.Close)

	event := create.CreatedUserEvent{
		Date:   time.Date(2024, time.March, 17, 17, 0, 0, 0, time.UTC),
		UserID: "test-user-id",
	}

	handler := mocks.NewCreatedUserEventHandlerMock(t)
	p := create.NewRabbitEventPublisher(sender.Queue.Name, sender.Channel, noLogger())
	l := create.NewRabbitEventListener(receiver.Queue.Name, receiver.Channel, handler, noLogger())

	listenerCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	handler.EXPECT().Handle(mock.Anything, event).Return(nil).Once().
		Run(func(_ mock.Arguments) {
			t.Log("cancelling listener context...")
			cancel()
		})

	// WHEN
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		p.Publish(context.Background(), event)
	}()

	go func() {
		defer wg.Done()

		listenErr := l.Listen(listenerCtx)
		if listenErr != nil {
			t.Logf("listen error: %s", listenErr)
		}
	}()

	wg.Wait()
}

type rabbitMqContainer struct {
	host string
	port string
}

func initRabbitMqContainer(t *testing.T, terminate bool) rabbitMqContainer {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.13-alpine",
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete"),
		Name:         "test-rabbit",
	}

	rabbitC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Logf("start rabbitmq container error: %s", err)
	}

	if terminate {
		t.Cleanup(func() {
			t.Log("terminating rabbitmq container...")

			if termErr := rabbitC.Terminate(context.Background()); termErr != nil {
				t.Logf("could not terminate rabbitmq container: %s", termErr)
			}
		})
	}

	host, err := rabbitC.Host(ctx)
	if err != nil {
		t.Logf("could not get rabbitmq host: %s", err)
	}

	port, err := rabbitC.MappedPort(ctx, "5672")
	if err != nil {
		t.Logf("could not get rabbitmq port: %s", err)
	}

	t.Logf("rabbitmq url %s:%s", host, port.Port())

	return rabbitMqContainer{
		host: host,
		port: port.Port(),
	}
}
