package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if len(args) < 2 || os.Args[1] == "" {
		s = "hello"
	}else {
		s = strings.Join(os.Args[1:], " ")
	}
	return s
}

func PubMessage()  {
	conn ,err := amqp.Dial("amqp://root:root@192.168.3.101:5672/")
	failOnErr(err, "Fail to connect rabbitmq_demo")
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	failOnErr(err, "declare exchange logs error")
	//q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	//failOnErr(err, "Failed to declare a queue")
	body := bodyFrom(os.Args)
	err = ch.Publish("logs", "", false, false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErr(err, "publish err")
}

func main() {
	PubMessage()
}
