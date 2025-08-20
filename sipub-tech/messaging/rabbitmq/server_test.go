package server_test

import (
    "context"
	"fmt"
	"strings"
    "testing"
	"time"

    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/rabbitmq"
)

const (
	rabbitMqImage = "rabbitmq:3"
	rabbitMqConnectionPort = "5672"
	rabbitMqUser = "guest"
	rabbitMqPassword = "guest"
	rabbitMqProtocol = "amqp"
)

var (
	rabbitMqExposedPorts = []string{"4369/tcp", "5672/tcp"}
)

func TestRabbitMqServer(t *testing.T) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        rabbitMqImage,
        ExposedPorts: rabbitMqExposedPorts,
        WaitingFor:   wait.ForLog("Ready to start client connection listeners"),
    }
    rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    testcontainers.CleanupContainer(t, rabbitmqC)
    require.NoError(t, err)

	endpoint, err := rabbitmqC.PortEndpoint(ctx, rabbitMqConnectionPort, rabbitMqProtocol)
	if err != nil {
		t.Error(err)
	}

	authInfo := fmt.Sprintf("%s:%s", rabbitMqUser, rabbitMqPassword)
	connectionUrl := insertAuthInfo(endpoint, authInfo)
	
	nodeId := "teste"

	t.Run("should open and close without errors", func(t *testing.T) {
		rabbitmqServer := server.NewRabbitMqServer(connectionUrl, nodeId)
		rabbitmqServer.Open()
		rabbitmqServer.Close()
	})

	t.Run("should be able to send and listen to messages and add correlationId and metadata to context", func(t *testing.T) {
		rabbitmqServer := server.NewRabbitMqServer(connectionUrl, nodeId)
		receivedChan := make(chan bool, 10)
		consumerFunction := func(ctx context.Context, body any) error {
			_, ok := ctx.Value(server.MetadataKey).(dtos.MessageMetadata)
			if !ok {
				t.Error("metadata is not set.")
			}
			_, ok = ctx.Value(server.CorrelationIdKey).(string)
			if !ok {
				t.Error("correlationId is not set.")
			}
			boolean, ok := body.(bool)
			if !ok {
				t.Error("body uncorrectly parsed")
			}
			receivedChan <- boolean
			return nil
		}
		const queueName = "testQueue"
		const producedValue = true
		
		rabbitmqServer.Open()
		defer rabbitmqServer.Close()
		_, producerFunction := rabbitmqServer.CreateProducer(queueName, nil, nil)

		if err := producerFunction(ctx, producedValue); err != nil {
			t.Errorf("Error found when running producer function: %v", err)
		}
		rabbitmqServer.RegisterConsumer(queueName, nil, nil, consumerFunction)
		rabbitmqServer.Listen(context.Background())

		select {
		case received := <- receivedChan:
			if received != producedValue {
				t.Errorf("Expected %+v, got %+v", received, producedValue)
			}
		case <- time.After(1 * time.Second):
			t.Errorf("Consumer Function was not called after 1 second.")
		}
		
	})
}

func insertAuthInfo(endpoint, authInfo string) string {
	parts := strings.Split(endpoint, "//")
	return parts[0] + "//" + authInfo + "@" + parts[1]
}
