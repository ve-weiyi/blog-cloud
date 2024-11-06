package rabbitmqx

import (
	"errors"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqChannel struct {
	connection *amqp.Connection
	rawChannel *amqp.Channel

	mu sync.RWMutex
}

func newRabbitChannel(conn *amqp.Connection) (*rabbitmqChannel, error) {
	rabbitCh := &rabbitmqChannel{
		connection: conn,
	}

	if err := rabbitCh.Connect(); err != nil {
		return nil, err
	}

	return rabbitCh, nil
}

func (r *rabbitmqChannel) Connect() (err error) {
	r.rawChannel, err = r.connection.Channel()
	return err
}

func (r *rabbitmqChannel) Close() (err error) {
	return r.rawChannel.Close()
}

func (r *rabbitmqChannel) Declare(q *QueueOptions, e *ExchangeOptions, b *BindingOptions) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error

	err = declareQueue(r.rawChannel, q)
	if err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}

	err = declareExchange(r.rawChannel, e)
	if err != nil {
		return fmt.Errorf("declare exchange failed: %w", err)
	}

	err = declareBindings(r.rawChannel, q, e, b)
	if err != nil {
		return fmt.Errorf("declare bindings failed: %w", err)
	}

	return nil
}

func (r *rabbitmqChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args map[string]interface{}) (<-chan amqp.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.rawChannel == nil {
		return nil, amqp.ErrClosed
	}

	return r.rawChannel.Consume(
		queue,                  // queue
		consumer,               // consumer
		autoAck,                // autoAck
		exclusive,              // exclusive
		noLocal,                // noLocal
		noWait,                 // noWait
		tableToAMQPTable(args), // args
	)
}

func (r *rabbitmqChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.rawChannel == nil {
		return errors.New("rawChannel is nil")
	}

	return r.rawChannel.Publish(
		exchange,  // exchange
		key,       // key
		mandatory, // mandatory
		immediate, // immediate
		msg)
}
