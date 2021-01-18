package task

import (
	"encoding/json"
	"gin_demo/interal"
	"gin_demo/interal/entity"
	"gin_demo/mq/rabbitmq"
	"github.com/sirupsen/logrus"
)


var mqClient mq_serivces.IMessageClient =
	mq_serivces.NewMessageClient(
		interal.GetRabbitEnv(),
		mq_serivces.NotifyTaskExchangeName,
		nil)

var AnswerNotifyChan = make(chan entity.AnswerNotify)
var voteAnswerNotifyChan = make(chan entity.VoteAnswerNotify)

// StartAnswerNotifyWorker 启动代理goroutine接受API-handler的消息发布
//目的是避免在API-handler中启动多个消息发布routine
func StartAnswerNotifyWorker() {
	go func() {
		for  {
			select {
			case message:=<-AnswerNotifyChan:
				data, err := json.Marshal(message)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"queue": mq_serivces.AnswerNotifyQueueName,
						"message": message,
						"error": err.Error(),
					}).Error("marshal message error")
				}
				mqClient.PublishToQueue(data, mq_serivces.AnswerNotifyQueueName)
			}
		}
	}()
}

func StartVoteAnswerNotifyWorker() {
	go func() {
		for  {
			select {
			case message:=<-voteAnswerNotifyChan:
				data, _ := json.Marshal(message)
				mqClient.PublishToQueue(data, mq_serivces.VoteNotifyQueueName)
			}
		}
	}()
}
