package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

var brokerList = []string{"192.168.3.101:32775", "192.168.3.101:32774"}

func connectAndWrite(topic string, partition int) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "192.168.3.101:32775", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("ddddd")},
		kafka.Message{Value: []byte("eeeeee")},
		kafka.Message{Value: []byte("ffffff")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func connectAndRead(topic string, partition int) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "192.168.3.101:32775", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	conn.SetReadDeadline(time.Now().Add(5*time.Second))
	batch := conn.ReadBatch(10e3, 1e6)
	for {
		b := make([]byte, 10e3) //10kb per message
		n, err := batch.Read(b)
		if err != nil {
			fmt.Printf("read error: %s\n", err.Error())
			fmt.Printf("current b:%s\tn:%d\n", string(b), n)
			break
		}
		fmt.Printf("read: %s\tlen=%d\n", string(b), n)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connect:", err)
	}
}

func ConsumerByGroup(ctx context.Context, topic string, who int) {
	//l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerList,
		Topic: topic,
		GroupID: "consumer-group-id2",
		//MinBytes: 10e3,
		//MaxBytes: 1e6,
		//MaxWait: time.Second * 5,
		StartOffset: kafka.LastOffset,
		//Logger: l,
	})
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("err when read message: %s\n", err.Error())
			break
		}
		fmt.Printf("who[%d] message at topic/partition/offset %v/%v/%v: %s = %s\n", who, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
func Consumer(ctx context.Context, topic string, partition int, offset int64) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerList,
		Topic: topic,
		Partition: partition,
		MinBytes: 10e3,
		MaxBytes: 1e6,

	})
	defer r.Close()

	r.SetOffset(offset)
	for {
		m, err := r.ReadMessage(ctx) //BLOCKING
		if err != nil {
			fmt.Printf("err when read message: %s\n", err.Error())
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}

func producer(ctx context.Context, topic string, addr string) {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	w := &kafka.Writer{
		Addr: kafka.TCP(addr),
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
		BatchBytes: 100,
		BatchTimeout: time.Second * 3,
		RequiredAcks: 1,
		Logger: l,
	}

	err := w.WriteMessages(ctx,
		kafka.Message{Key: []byte("key-A"), Value: []byte("hello 111")},
		kafka.Message{Key: []byte("key-B"), Value: []byte("hello 222")},
		kafka.Message{Key: []byte("key-C"), Value: []byte("hello 333")},
		kafka.Message{Key: []byte("key-D"), Value: []byte("hello 444")},
	)

	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	fmt.Println("producer done")
}
func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	//defer cancel()
	//connectAndWrite("test001", 1)
	producer(ctx, "test001", "192.168.3.101:32774")
	//producer(ctx, "test001", "192.168.3.101:32776")
	//start 4 consumer in one Group
	for i := 0; i < 4; i++ {
		go ConsumerByGroup(ctx, "test001", i)
	}
	select {
	case <-ctx.Done():
	}
}