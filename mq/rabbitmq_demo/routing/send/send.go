package main

import (
	"fmt"
	"gin_demo/mq/rabbitmq_demo/routing"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

func generateKey(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "info"
	}else {
		s = args[1]
	}
	return s
}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	}else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func Send(channel *amqp.Channel)  {
	body := bodyFrom(os.Args)
	err := channel.Publish("logs_direct", generateKey(os.Args), false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(body),
		})
	routing.HandlerError(err, "publish error")
	fmt.Printf("sent message %s\n", body)
}

func main() {
	Send(routing.Declare())
}
