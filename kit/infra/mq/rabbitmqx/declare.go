package rabbitmqx

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareQueue(channel *amqp.Channel, q *QueueOptions) error {
	if q == nil {
		return nil
	}

	queue, err := channel.QueueDeclare(
		q.Name,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		tableToAMQPTable(q.Args),
	)
	if err != nil {
		return err
	}

	if q.Name != queue.Name {
		q.Name = queue.Name
	}

	return nil
}

func declareExchange(channel *amqp.Channel, e *ExchangeOptions) error {
	if e == nil {
		return nil
	}

	err := channel.ExchangeDeclare(
		e.Name,
		e.Kind,
		e.Durable,
		e.AutoDelete,
		e.Internal,
		e.NoWait,
		tableToAMQPTable(e.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareBindings(channel *amqp.Channel, q *QueueOptions, e *ExchangeOptions, b *BindingOptions) error {
	if q == nil || e == nil || b == nil {
		return nil
	}

	return channel.QueueBind(
		q.Name,
		b.RoutingKey,
		e.Name,
		b.NoWait,
		tableToAMQPTable(b.Args),
	)
}
