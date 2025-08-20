package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func StandardQueueConfig() *QueueConfig {
	return &QueueConfig{
		durable: true,
		deleteWhenUnused: false,
		exclusive: false,
		noWait: false,
		arguments: nil,
	}
}

func NewQueueConfig(durable, deleteWhenUnused, exclusive, noWait bool, arguments amqp.Table) *QueueConfig {
	return &QueueConfig{
		durable,
		deleteWhenUnused,
		exclusive,
		noWait,
		arguments,
	}
}

type QueueConfig struct {
	durable bool
	deleteWhenUnused bool
	exclusive bool
	noWait bool
	arguments amqp.Table
}

func StandardConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		autoAcknowledge: false,
		exclusive: false,
		noLocal: false,
		noWait: false,
		arguments: nil,
	}
}

func NewConsumerConfig(autoAcknowledge, exclusive, noLocal, noWait bool, arguments amqp.Table) *ConsumerConfig {
	return &ConsumerConfig{
		autoAcknowledge,
		exclusive,
		noLocal,
		noWait,
		arguments,
	}
}

type ConsumerConfig struct {
	autoAcknowledge bool
	exclusive bool
	noLocal bool
	noWait bool
	arguments amqp.Table
}

func StandardProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		exchange: "",
		mandatory: false,
		immediate: false,
		deliveryMode: amqp.Persistent,
	}
}

func NewProducerConfig(exchange string, mandatory, immediate bool, deliveryMode uint8) *ProducerConfig {
	return &ProducerConfig{
		exchange,
		mandatory,
		immediate,
		deliveryMode,
	}
}

type ProducerConfig struct {
	exchange string
	mandatory bool
	immediate bool
	deliveryMode uint8
}


