package main

import (
    "fmt"
    "github.com/streadway/amqp"
    "log"
)

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    tag     string
    done    chan error
}

func AMQPConsumer(uri, key, ctag string, conversation *Conversation) (*Consumer, error) {
    exchangeName := "encrypted-messaging.exchange"
    queueName := "encrypted-messaging.incoming"

    c := &Consumer{
        conn:    nil,
        channel: nil,
        tag:     ctag,
        done:    make(chan error),
    }

    var err error

    c.conn, err = amqp.Dial(uri)
    if err != nil {
        return nil, fmt.Errorf("Dial: %s", err)
    }

    go func() {
        fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
    }()

    log.Printf("got Connection, getting Channel")
    c.channel, err = c.conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("Channel: %s", err)
    }

    log.Printf("got Channel, declaring Exchange (%q)", exchangeName)
    if err = c.channel.ExchangeDeclare(
        exchangeName,     // name of the exchange
        "direct",         // type
        true,             // durable
        false,            // delete when complete
        false,            // internal
        false,            // noWait
        nil,              // arguments
    ); err != nil {
        return nil, fmt.Errorf("Exchange Declare: %s", err)
    }

    log.Printf("declared Exchange, declaring Queue %q", queueName)
    queue, err := c.channel.QueueDeclare(
        queueName,                     // name of the queue
        true,                          // durable
        false,                         // delete when usused
        false,                         // exclusive
        false,                         // noWait
        nil,                           // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Declare: %s", err)
    }

    log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
        queue.Name, queue.Messages, queue.Consumers, key)

    if err = c.channel.QueueBind(
        queue.Name,                     // name of the queue
        key,                            // bindingKey
        exchangeName,                   // sourceExchange
        false,                          // noWait
        nil,                            // arguments
    ); err != nil {
        return nil, fmt.Errorf("Queue Bind: %s", err)
    }

    log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
    deliveries, err := c.channel.Consume(
        queue.Name, // name
        c.tag,      // consumerTag,
        false,      // noAck
        false,      // exclusive
        false,      // noLocal
        false,      // noWait
        nil,        // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("Queue Consume: %s", err)
    }

    go handle(deliveries, conversation, c.done)
    select{}
}

func (c *Consumer) Shutdown() error {
    // will close() the deliveries channel
    if err := c.channel.Cancel(c.tag, true); err != nil {
        return fmt.Errorf("Consumer cancel failed: %s", err)
    }

    if err := c.conn.Close(); err != nil {
        return fmt.Errorf("AMQP connection close error: %s", err)
    }

    defer log.Printf("AMQP shutdown OK")

    // wait for handle() to exit
    return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, c *Conversation,done chan error) {
    for d := range deliveries {
        c.ReceiveMessage( []byte(d.Body) )
        d.Ack(true)
    }
    log.Printf("handle: deliveries channel closed")
    done <- nil
}

func AMQPPublisher(amqpUri, key, body string) error {
    log.Println("go")
    exchangeName := "encrypted-messaging.exchange"
//    queueName := "encrypted-messaging.incoming"

    connection, err := amqp.Dial(amqpUri)
    if err != nil {
        return fmt.Errorf("Dial: %s", err)
    }
    defer connection.Close()

    channel, err := connection.Channel()
    if err != nil {
        return fmt.Errorf("Channel: %s", err)
    }

    if err := channel.ExchangeDeclare(
        exchangeName,     // name
        "direct",         // type
        true,             // durable
        false,            // auto-deleted
        false,            // internal
        false,            // noWait
        nil,              // arguments
    ); err != nil {
        return fmt.Errorf("Exchange Declare: %s", err)
    }

    log.Printf("enabling publishing confirms.")
    if err := channel.Confirm(false); err != nil {
        return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
    }

    confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

    defer confirmOne(confirms)

    if err = channel.Publish(
        exchangeName,   // publish to an exchange
        key,            // routing to 0 or more queues
        false,          // mandatory
        false,          // immediate
        amqp.Publishing{
            Headers:         amqp.Table{},
            ContentType:     "text/plain",
            ContentEncoding: "",
            Body:            []byte(body),
            DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
            Priority:        0,              // 0-9
        },
    ); err != nil {
        return fmt.Errorf("Exchange Publish: %s", err)
    }

    return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
    if confirmed := <-confirms; confirmed.Ack {
        log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
    } else {
        log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
    }
}
