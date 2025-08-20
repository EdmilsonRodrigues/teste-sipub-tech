package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/messaging/dtos"
)


type ConsumerFunction func(ctx context.Context, body any) error
type ProducerFunction func(ctx context.Context, body any) error
type ContextKey string

const CorrelationIdKey ContextKey = "correlationId"
const MetadataKey ContextKey = "metadata"

func NewRabbitMqServer(connectionUrl, nodeId string) *RabbitMqServer {
	return &RabbitMqServer{
		connectionUrl: connectionUrl,
		nodeId: nodeId,
	}
}

type RabbitMqServer struct {
	connectionUrl string
	nodeId string
	
	conn *amqp.Connection
	ch   *amqp.Channel

	consumers []*consumerData
}

func (rmqServer *RabbitMqServer) Open() {
	rmqServer.connect()
	rmqServer.openChannel()
	rmqServer.setQoS()
}

func (rmqServer *RabbitMqServer) Close() {
	if rmqServer.conn != nil {
		if err := rmqServer.conn.Close(); err != nil {
			log.Printf("Error ocurred closing RabbitMq Connection: %v", err)
		}
	}
	if rmqServer.ch != nil {
		if err := rmqServer.ch.Close(); err != nil {
			log.Printf("Error ocurred closing RabbitMq Channel: %v", err)
		}
	}
}

func (rmqServer *RabbitMqServer) Listen(ctx context.Context) {
	for _, consumer := range rmqServer.consumers {
		go rmqServer.consumeForever(ctx, consumer)
	}
}

func (rmqServer *RabbitMqServer) RegisterConsumer(
	queueName string,
	queueConfig *QueueConfig,
	consumerConfig *ConsumerConfig,
	consumerFunction ConsumerFunction,
) {
	if queueConfig == nil {
		queueConfig = StandardQueueConfig()
	}
	if consumerConfig == nil {
		consumerConfig = StandardConsumerConfig()
	}

	queue := rmqServer.declareQueue(queueName, queueConfig)
	consumer := rmqServer.registerConsumer(queue, consumerConfig)

	consumerData := &consumerData{
		queue: queue,
		consumer: consumer,
		consumerFunction: consumerFunction,
	}

	rmqServer.consumers = append(rmqServer.consumers, consumerData)
}

func (rmqServer *RabbitMqServer) CreateProducer(
	queueName string,
	queueConfig *QueueConfig,
	producerConfig *ProducerConfig,
) (amqp.Queue, ProducerFunction) {
	if queueConfig == nil {
		queueConfig = StandardQueueConfig()
	}
	if producerConfig == nil {
		producerConfig = StandardProducerConfig()
	}
	
	queue := rmqServer.declareQueue(queueName, queueConfig)

	return queue, func(ctx context.Context, body any) error {
		correlationId, ok := ctx.Value(CorrelationIdKey).(string)
		if ok {
			correlationId = correlationId + "-"
		} else {
			correlationId = ""
		}
		newCorrelationId := fmt.Sprintf(
			"%s%s[%s-%s]",
			correlationId,
			rmqServer.nodeId,
			queue.Name,
			uuid.New().String(),
		)
		messageBody := dtos.Message{
			Metadata: dtos.MessageMetadata{CorrelationId: newCorrelationId},
			Data: body,
		}
		bytes, err := json.Marshal(messageBody)
		if err != nil {
			return fmt.Errorf("error marshalling body: %w", err)
		}

		err = rmqServer.ch.PublishWithContext(
			ctx,
			producerConfig.exchange,
			queue.Name,
			producerConfig.mandatory,
			producerConfig.immediate,
			amqp.Publishing{
				DeliveryMode: producerConfig.deliveryMode,
				ContentType: "application/json",
				Body: bytes,
			},
		)
		if err != nil {
			return fmt.Errorf("error publishing message: %w", err)
		}

		return nil
	}
}


func (rmqServer *RabbitMqServer) connect() {
	conn, err := amqp.Dial(rmqServer.connectionUrl)
	rmqServer.failOnError(err, "Failed to connect to RabbitMQ")
	rmqServer.conn = conn
}

func (rmqServer *RabbitMqServer) openChannel() {
	ch, err := rmqServer.conn.Channel()
	rmqServer.failOnError(err, "Failed to open a channel")
	rmqServer.ch = ch
}

func (rmqServer *RabbitMqServer) declareQueue(queueName string, queueConfig *QueueConfig) amqp.Queue {
	q, err := rmqServer.ch.QueueDeclare(
		queueName, 
		queueConfig.durable,   
		queueConfig.deleteWhenUnused,   
		queueConfig.exclusive,   
		queueConfig.noWait,   
		queueConfig.arguments,     
	)
	rmqServer.failOnError(err, fmt.Sprintf("Failed to declare %q queue", queueName))
	return q
}

func (rmqServer *RabbitMqServer) setQoS() {
	err := rmqServer.ch.Qos(
		4,
		0,
		false,
	)
	rmqServer.failOnError(err, "Failed to set QoS")
}

func (rmqServer *RabbitMqServer) registerConsumer(queue amqp.Queue, consumerConfig *ConsumerConfig) <-chan amqp.Delivery {
	msgs, err := rmqServer.ch.Consume(
		queue.Name, 
		rmqServer.nodeId,     
		consumerConfig.autoAcknowledge,  
		consumerConfig.exclusive,
		consumerConfig.noLocal,
		consumerConfig.noWait,
		consumerConfig.arguments,
	)
	rmqServer.failOnError(err, fmt.Sprintf("Failed to register %q consumer", queue.Name))
	return msgs
}

func (rmqServer *RabbitMqServer) consumeForever(ctx context.Context, consumer *consumerData) {
	for delivery := range consumer.consumer {
		var message dtos.Message
		
		if err := json.Unmarshal(delivery.Body, &message); err != nil {
			log.Printf("Error unmarshalling body: %+v", delivery.Body)
			return
		}

		newCorrelationId := fmt.Sprintf(
			"%s-%s[%s-%s]",
			message.Metadata.CorrelationId,
			rmqServer.nodeId,
			consumer.queue.Name,
			uuid.New().String(),
		)

		internalContext := context.WithValue(ctx, CorrelationIdKey, newCorrelationId)
		internalContext = context.WithValue(internalContext, MetadataKey, message.Metadata)

		if err := consumer.consumerFunction(internalContext, message.Data); err != nil {
			log.Printf("Error ocurred on consumerFunction for queue %q: %v", consumer.queue.Name, err)
		}
		if err := delivery.Ack(false); err != nil {
			log.Printf("Error deliveryng acknowledgemennt for message %+v", message)
		}
	}
}

func (rmqServer *RabbitMqServer) failOnError(err error, msg string) {
	if err != nil {
		rmqServer.Close()
		log.Panicf("%s: %s", msg, err)
	}
}


type consumerData struct {
	queue amqp.Queue
	consumer <- chan amqp.Delivery
	consumerFunction ConsumerFunction
}

