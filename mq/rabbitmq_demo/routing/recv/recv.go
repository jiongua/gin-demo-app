package main

import (
	"fmt"
	"gin_demo/mq/rabbitmq_demo/routing"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func Sub(channel *amqp.Channel, bindKey string)  {
	q, err := channel.QueueDeclare("", false, false, false, false, nil)
	err = channel.QueueBind(q.Name, bindKey, "logs_direct", false, nil)
	routing.HandlerError(err, fmt.Sprintf("bind queue[%s] with key[%s] to logs_direct error", q.Name, bindKey))

	fmt.Printf("declare queue[%s] with key[%s] to logs_direct\n", q.Name, bindKey)

	forever := make(chan bool)
	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	go func() {
		for d := range msgs {
			fmt.Printf("get message %s\n", string(d.Body))
		}
	}()
	<-forever
}

var validKey = []string{"error", "info", "warning"}

func Customer() string {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./cmd [error|info|warning]")
	}
	if !strings.Contains(strings.Join(validKey, " "), os.Args[1]) {
		log.Fatal("usage: ./cmd [error|info|warning]")
	}
	return os.Args[1]
}

func main() {

	key := Customer()
	ch := routing.Declare()
	Sub(ch, key)
}