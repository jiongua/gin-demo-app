package mq_serivces

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
	"time"
)

type IMessageClient interface {
	ConnectToBroker(connectString string) error
	Publish(data []byte, exchangeName, queueName string)
	PublishToQueue(data []byte, queueName string)
	SubscribeFromQueue(queueName string) <-chan amqp.Delivery
}

//RabbitMQ通知客户端
type MessageClient struct {
	conn *amqp.Connection
	ExchangeName string
	channelMapping map[string]*amqp.Channel
	recoverConnHandler func(err error)
}

func NewMessageClient(connString, exchangeName string, recoverHandler func(err error)) *MessageClient{
	var client = &MessageClient{}
	err := client.ConnectToBroker(connString)
	if err != nil {
		log.Fatalf("connect to rabbitmq in NewMessageClient error: %s", err.Error())
	}
	client.ExchangeName = exchangeName
	client.recoverConnHandler = recoverHandler
	client.channelMapping = make(map[string]*amqp.Channel)
	return client
}

//getOrCreateChannel获取指定消息队列相应的amqp-channel
func (m *MessageClient) getOrCreateChannel(queueName string) *amqp.Channel {
	//根据queueName，同一类消息发布到同一个queue
	//新回答--> AnswerNotifyQueueName
	//回答被点赞--> VoteAnswerNotifyQueueName
	//评论被点赞--> VoteCommentNotifyQueueName
	//回答被评论--> CommentAnswerNotifyQueueName
	if channel, ok := m.channelMapping[queueName]; ok {
		return channel
	} else {
		channel, err := m.conn.Channel()
		m.ErrorHandle(err, "get channel from connect error")
		m.channelMapping[queueName] = channel
		return channel
	}
}

func (m *MessageClient) Publish(data []byte, exchangeName, queueName string) {
	publishing := amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "json",
		DeliveryMode:    0,
		Timestamp:       time.Now(),
		Body:            data,
	}
	channel := m.getOrCreateChannel(queueName)
	err := channel.Publish(exchangeName, queueName, false, false, publishing)
	m.ErrorHandle(err, fmt.Sprintf("Publish to %s:%s error", exchangeName, queueName))
}

//PublishToQueue 将已序列化的数据发布到rabbitmq
//对于每一类消息，使用一个固定的goroutine代理发送，避免在controller使用多个goroutine向
//同一个rabbitmq-channel发布消息
func (m *MessageClient) PublishToQueue(data []byte, queueName string) {
	publishing := amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "json",
		DeliveryMode:    0,
		Timestamp:       time.Now(),
		Body:            data,
	}
	channel := m.getOrCreateChannel(queueName)
	err := channel.Publish(m.ExchangeName, queueName, false, false, publishing)
	if err == nil {
		log.WithFields(logrus.Fields{
			"queue": queueName,
			"message": string(data),
		}).Info("publish message ok!")
	}
	//deal err
	m.ErrorHandle(err, "publish to queue error")
}

func (m *MessageClient) RegisterRecoverConnHandler(handler func(err error)) {
	m.recoverConnHandler = handler
}

func (m *MessageClient) ErrorHandle(err error, msg string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error(msg)
		if strings.HasPrefix(err.Error(), "Exception (504) Reason") &&
			m.recoverConnHandler != nil {
			m.recoverConnHandler(err)
		} else if m.recoverConnHandler == nil {
			panic("Exception (504) Reason")
		}
	}
}

func (m *MessageClient) ConnectToBroker(connectString string) error {
	conn, err := amqp.Dial(connectString)
	if err != nil {
		log.WithFields(logrus.Fields{
			"url": connectString,
		}).Error("Fail to connect rabbitmq")
		return err
	}
	m.conn = conn
	return nil
}

//SubscribeFromQueue 从指定消息队列获取消息
// 避免多个goroutine共用同一个rabbitmq channel, 使用代理goroutine调用SubscribeFromQueue
func (m *MessageClient) SubscribeFromQueue(queueName string) <-chan amqp.Delivery{
	channel := m.getOrCreateChannel(queueName)
	//defer channel.Close()
	goChan, err := channel.Consume(queueName, "", true, false, false,false,nil)
	m.ErrorHandle(err, "consume message error")
	return goChan
}

