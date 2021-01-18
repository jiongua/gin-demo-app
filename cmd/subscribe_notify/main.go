//消费rabbitmq消息
package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"gin_demo/interal/entity"
	"gin_demo/mq/rabbitmq"
	//"os/signal"
	//"runtime"
)

var mqClient mq_serivces.IMessageClient =
	mq_serivces.NewMessageClient(
		mq_serivces.GetRabbitEnv(),
		mq_serivces.NotifyTaskExchangeName,
		nil)


func subscribeAnswerNotify() {
	goChan := mqClient.SubscribeFromQueue(mq_serivces.AnswerNotifyQueueName)
	limit := make(chan struct{}, 10)
	for  {
		select {
		case data:=<-goChan: {
			fmt.Println("get new answer notify...")
			notify := &entity.AnswerNotify{}
			if err := json.Unmarshal(data.Body, notify); err != nil {
				//todo
				fmt.Printf("%s\n", err.Error())
			} else {
				fmt.Printf("consume subscribe_notify: %v\n", notify)
				go func() {
					limit<- struct{}{}
					notify.CopyAndCreateAll()
					<-limit
				}()
			}
		}
			//case <-done:
			//
		}
	}
}

func main() {
	go subscribeAnswerNotify()
	select {

	}
}


