package main

import (
	"fmt"
	"gin_demo/mq/rabbitmq_demo"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}


func SubMessage()  {
	conn := rabbitmq_demo.DefaultConn()
	ch, err := conn.Channel()
	failOnErr(err, "fail to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	failOnErr(err, "declare exchange logs error")
	q, err := ch.QueueDeclare("", true, false, false, false, nil)
	failOnErr(err, "Failed to declare a queue")
	fmt.Printf("declare queue[%s] success\n", q.Name)
	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	failOnErr(err, "queue bind error")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnErr(err, "Consume err")
	done := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("recv msg: %s from queue[%s]\n", string(d.Body), q.Name)
		}
	}()
	<-done
}

func main() {
	SubMessage()
}



