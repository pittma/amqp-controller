Go amqp-controller
==================

## Purpose

To provide a useful abstraction to the amqp lib.  This controller provides a simple API for listening and publishing messages on an AMQP exchange.

## WIP

There currently no tests, and only a few config items available.  More to come.

## Config file based

The controller loada a config via [gcfg](https://code.google.com/p/gcfg/) (which is based loosely on gitconfig).

```
[connection]
ip = 127.0.0.1
port = 5672
user = guest
pass = guest

[bindings]
exchange = test.exch
routes = test.key.1
routes = test.key.2
queuename = simple-consumer
```

## Examples

### As a Consumer

```go
package main

import (
    "log"
    "os/signal"
    "os"
    "github.com/danielscottt/amqp-controller"
    "github.com/streadway/amqp"
)

func catchInterrupt(c *controller.Controller) {
    sigint := make(chan os.Signal, 10)
    signal.Notify(sigint, os.Interrupt)
    <-sigint
    c.Shutdown()
    os.Exit(0)
}

func handleMsgs(msg amqp.Delivery) {
    log.Printf("Message received! [%s]", string(msg.Body))
    msg.Ack(false)
}

func main() {
    c := controller.NewController()
    go catchInterrupt(c)

    c.Bind()
    c.Consume(handleMsgs)

    defer c.Shutdown()
    // listen forever
    select {}
}
```

### Publishing Messages

```go
package main

import (
    "github.com/danielscottt/amqp-controller"
    "flag"
)

var body = flag.String("body", "foobar", "Body of message")

func init() {
    flag.Parse()
}

func main() {
    c := controller.NewController()
    var headers map[string]interface{}

    c.Publish("test.key.1", headers, "text/plain", "", *body, 2, 0)
}
```
