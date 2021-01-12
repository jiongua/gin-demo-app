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

func NewTask()  {
	conn ,err := amqp.Dial("amqp://root:root@192.168.3.101:5672/")
	failOnErr(err, "Fail to connect rabbitmq_demo")
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	failOnErr(err, "Failed to declare a queue")
	body := bodyFrom(os.Args)
	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing {
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErr(err, "publish err")
	fmt.Printf("[x] sent %s\n", body)
}

func main() {
	NewTask()
}



