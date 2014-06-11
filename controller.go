package controller

import (
    "log"
    "github.com/streadway/amqp"
    "fmt"
)

type Controller struct {
    conn *amqp.Connection
    Ch *amqp.Channel
    connected bool
    bound bool
}

type messageHandler func(amqp.Delivery)

func haltOnError(err error, msg string) {
    if err != nil {
      log.Fatalf("%s: %s", msg, err)
      panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func NewController() (*Controller) {

    var err error
    LoadConfig()

    c := &Controller{
        connected: false,
        bound:     false,
    }

    uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", Cfg.Connection.User, Cfg.Connection.Pass, Cfg.Connection.Ip, Cfg.Connection.Port)
    log.Printf("Connecting to server [%s]", uri)
    c.conn, err = amqp.Dial(uri)
    haltOnError(err, "connection error")

    log.Printf("Establishing channel")
    c.Ch, err = c.conn.Channel()
    haltOnError(err, "error establishing channel")

    log.Printf("Declaring exchange [%s]", Cfg.Bindings.Exchange)
    err = c.Ch.ExchangeDeclare(Cfg.Bindings.Exchange, "topic", true, false, false, false, nil)
    haltOnError(err, "exchange declare error")

    c.connected = true

    return c
}

func (c *Controller) Bind() {
    if c.connected {
        log.Printf("Binding to routing keys:")
        for _, route := range Cfg.Bindings.Routes {
            _, err := c.Ch.QueueDeclare(Cfg.Bindings.QueueName, true, false, false, false, nil)
            haltOnError(err, "queue declare error")

            log.Printf("  - %s", route)
            c.Ch.QueueBind(Cfg.Bindings.QueueName, route, Cfg.Bindings.Exchange, false, nil)
            haltOnError(err, "queue binding error")
        }
        c.bound = true
    } else {
        err := NotConnectedError{"NotConnected"}
        haltOnError(err, "Cannot bind to queue, not yet connected to server")
    }
}

func (c *Controller) Consume(msgHandler messageHandler) {
    if c.bound {
        log.Printf("listening for messages...")
        msgs, err := c.Ch.Consume(Cfg.Bindings.QueueName, "", false, false, false, false, nil)
        haltOnError(err, "error on message receipt")

        go handle(msgs, msgHandler)

    } else {
        err := NotBoundError{"QueueNotBound"}
        haltOnError(err, "Cannot consume, not yet bound to any queue")
    }
}

func (c *Controller) Shutdown() {
    log.Printf("Shutting down")

    err := c.Ch.Close()
    haltOnError(err, "Error closing channel")

    err = c.conn.Close()
    haltOnError(err, "Error closing connection")

    log.Printf("clean shutdown successful")
}

func (c *Controller) Publish(routingKey string, headers map[string]interface{}, contentType, contentEncoding, body string, mode, priority uint8) error {
    log.Printf("Publishing message with routing key [%s]", routingKey)
    err := c.Ch.Publish(
        Cfg.Bindings.Exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            Headers:         headers,
            ContentType:     contentType,
            ContentEncoding: contentEncoding,
            Body:            []byte(body),
            DeliveryMode:    mode,
            Priority:        priority,
        },
    )
    haltOnError(err, "Failed to publish")
    log.Printf("message sent")
    return nil
}

func handle(msgs <-chan amqp.Delivery, msgHandler messageHandler) {
    for msg := range msgs {
      msgHandler(msg)
    }
}
