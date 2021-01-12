package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func StartCustomer() {
	conn ,err := amqp.Dial("amqp://root:root@192.168.3.101:5672/")
	failOnErr(err, "Fail to connect rabbitmq_demo")
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnErr(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main()  {
	StartCustomer()
}