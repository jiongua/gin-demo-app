package mq_serivces

import (
	glog "gin_demo/interal/log"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var log = glog.Log

const (
	NotifyTaskExchangeName = "zhihu_notify_exchange"
	AnswerNotifyQueueName = "answer_notify_queue"
	VoteNotifyQueueName= "vote_notify_queue"
	CommentNotifyQueueName= "comment_notify_queue"

)

func init() {
	DeclareDefault()
}

func NewConnOrNil() *amqp.Connection {
	url := GetRabbitEnv()
	conn, err := amqp.Dial(url)
	if err != nil {
		log.WithFields(logrus.Fields{
			"url": url,
		}).Error("Fail to connect rabbitmq_demo")
		return nil
	}
	return conn
}

func NewChannelOrNil(conn *amqp.Connection) *amqp.Channel{
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Fail to open channel: %s", err.Error())
		return nil
	}
	return ch
}

func DeclareDefault()  {
	conn := NewConnOrNil()
	if conn == nil {
		log.Fatal("init connect to rabbitmq fails")
	}
	ch := NewChannelOrNil(conn)
	if conn == nil {
		log.Fatal("init channel rabbitmq fails")
	}
	defer conn.Close()
	defer ch.Close()
	DeclareExchange(ch, NotifyTaskExchangeName)
	DeclareAndBindQueue(ch, NotifyTaskExchangeName, AnswerNotifyQueueName)
	DeclareAndBindQueue(ch, NotifyTaskExchangeName, CommentNotifyQueueName)
	DeclareAndBindQueue(ch, NotifyTaskExchangeName, VoteNotifyQueueName)
}

func DeclareExchange(c *amqp.Channel, exchangeName string) {
	if err := c.ExchangeDeclare(exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil); err != nil {
		log.WithFields(logrus.Fields{
			"exchange_name": exchangeName,
			"error":         err.Error(),
		}).Error("declare exchange error")
		panic("declare exchange error")
	}
}

func DeclareAndBindQueue(c *amqp.Channel, exchangeName, queueName string) *amqp.Queue{
	q, err := c.QueueDeclare(queueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"queue_name": queueName,
			"error": err.Error(),
		}).Error("declare queue error")
		panic("declare queue error")
	}
	err = c.QueueBind(queueName,
		queueName,
		exchangeName,
		false,
		nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"exchange_name": exchangeName,
			"queue_name": queueName,
			"error": err.Error(),
		}).Error("bind queue error")
		panic("bind queue error")
	}
	return &q
}
