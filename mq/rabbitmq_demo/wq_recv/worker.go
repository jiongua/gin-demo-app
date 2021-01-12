package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}

func StartWorker() {
	conn ,err := amqp.Dial("amqp://root:root@192.168.3.101:5672/")
	failOnErr(err, "Fail to connect rabbitmq_demo")
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	failOnErr(err, "Failed to declare a queue")
	tasks, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // 消费者手动ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	err = ch.Qos(1,0,false)
	failOnErr(err, "set qos error")
	forever := make(chan bool)
	go func() {
		for task := range tasks {
			fmt.Printf("woker Received a message: %s\n", task.Body)
			//do worker

			t := bytes.Count(task.Body, []byte("."))
			time.Sleep(time.Second * time.Duration(t))
			fmt.Printf("worker done\n")
			task.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	StartWorker()
}