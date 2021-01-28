package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type WriterClient struct {
	*kafka.Writer
}

type ReaderClient struct {
	*kafka.Reader

}

func NewDefaultWriter(brokerList []string, topic string) *WriterClient{
	return &WriterClient{Writer: &kafka.Writer{
		Addr: kafka.TCP(brokerList...),
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
		RequiredAcks: 1,
	}}
}

func NewDefaultReader(brokerList []string, topic string, groupID string) *ReaderClient{
	return &ReaderClient{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokerList,
			Topic: topic,
			GroupID: groupID,
			//MinBytes: 10e3,
			//MaxBytes: 1e6,
			//MaxWait: time.Second * 5,
			StartOffset: kafka.LastOffset,
		}),
	}
}

//Push 发布消息到主题
//Each writer is bound to a single topic, to write to multiple topics,
//a program must create multiple writers.
//如果消息顺序不重要，key可以设置为nil，每条消息将rr到各分区
func (k *WriterClient) Produce(ctx context.Context, key, value []byte) error {
	//block until batch reach
	return k.WriteMessages(ctx,
		kafka.Message{
			Key: key,
			Value: value,
		},
	)
}

//Consume 消费者组其中的一个消费者从kafka读取消息, 并将消息写入out chan供worker goroutine处理
func (k *ReaderClient) Consume(ctx context.Context, out chan<- kafka.Message) {
	defer k.Close()
	for {
		m, err := k.ReadMessage(ctx)
		//TODO error handler
		if err != nil {
			log.Printf("err when read message: %s\n", err.Error())
			break
		}
		out<-m
		//fmt.Printf("who[%d] message at topic/partition/offset %v/%v/%v: %s = %s\n", who, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}


