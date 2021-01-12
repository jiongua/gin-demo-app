package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func SendMessage()  {
	conn ,err := amqp.Dial("amqp://root:root@192.168.3.101:5672/")
	failOnErr(err, "Fail to connect rabbitmq_demo")
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnErr(err, "Failed to declare a queue")
	body := "Hello World456!"
	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErr(err, "publish err")
}

func main() {
	SendMessage()
}



