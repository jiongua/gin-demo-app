//消费rabbitmq消息
//TODO 独立服务
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	//"context"
	"encoding/json"
	"fmt"
	"gin_demo/interal/entity"
	logger "gin_demo/interal/log"
	"gin_demo/interal/query"
	"gin_demo/mq/rabbitmq"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"


	//"os/signal"
	//"runtime"
)
var log = logger.Log

var mqClient mq_serivces.IMessageClient =
	mq_serivces.NewMessageClient(
		mq_serivces.GetRabbitEnv(),
		mq_serivces.NotifyTaskExchangeName,
		nil)


//StartConsumerAnswerNotify 从rabbitmq返回chan读取消息，写入db
func StartConsumerAnswerNotify(ctx context.Context, done chan<- struct{}) {
	goChan := mqClient.SubscribeFromQueue(mq_serivces.AnswerNotifyQueueName)
	notifyChan := make(chan *entity.AnswerNotify)
	w := new(sync.WaitGroup)
	for i := 0; i < 4; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			for data := range notifyChan {
				NotifyAll(data)
			}
			log.Info("worker exit")
		}()
	}

	for  {
		select {
			case <-ctx.Done():
				//parent cancel, wait child goroutines done
				log.Info("parent cancel, wait worker exit")
				close(notifyChan)
				w.Wait()
				log.Info("consumer exit")
				done<- struct{}{}
				return
			case data:=<-goChan: {
				fmt.Println("get new answer notify...")
				notify := &entity.AnswerNotify{}
				if err := json.Unmarshal(data.Body, notify); err != nil {
					//todo
					fmt.Printf("%s\n", err.Error())
				} else {
					fmt.Printf("consume subscribe_notify: %v\n", notify)
					notifyChan<-notify
				}
			}
		}
	}
}

// NotifyAll 生成每个关注人的通知，除发布人和禁用通知的接收者
func NotifyAll(m *entity.AnswerNotify)  {
	sender, err := query.GetUserByID(m.SenderID)
	if err != nil {
		log.Errorf("get SenderName by SenderID[%s] error: %s\n", m.SenderID, err.Error())
		return
	}
	m.SenderName = sender.Name

	receiverIDs := query.GetAllAttentionUsersOrNil(m.QuestionID)
	batchSize := len(receiverIDs)
	notifies := make([]entity.AnswerNotify, 0, batchSize)
	for i, id := range receiverIDs {
		//过滤接收者是SenderID
		if uuid.Equal(m.SenderID, id) {
			continue
		}
		//todo 过滤屏蔽通知的用户
		notifyItem := entity.AnswerNotify{}
		err := copier.Copy(notifyItem, m)
		if err != nil {
			log.Errorf("copier error: %s\n", err.Error())
			continue
		}
		notifyItem.ReceiverID = id
		notifies[i] = notifyItem
	}
	if batchSize > 0 {
		entity.Db().CreateInBatches(&notifies, batchSize)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	done := make(chan struct{})

	go StartConsumerAnswerNotify(ctx, done)
	select {
		case <-sig:
			cancel()
			<-done
	}
	log.Info("main exit")
}


