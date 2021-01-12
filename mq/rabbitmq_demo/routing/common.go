package routing

import (
	"gin_demo/mq/rabbitmq_demo"
	"github.com/streadway/amqp"
	"log"
)

func HandlerError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s\n", msg, err.Error())
	}
}

func Declare() *amqp.Channel {
	conn := rabbitmq_demo.DefaultConn()
	ch, err := conn.Channel()
	HandlerError(err, "fail open channel")
	err = ch.ExchangeDeclare("logs_direct", "direct", true, false, false, false,nil)
	HandlerError(err, "fail declare logs_direct exchange")
	return ch
}
