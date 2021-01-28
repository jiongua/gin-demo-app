package task
//
//import (
//	"context"
//	logger "gin_demo/interal/log"
//	kafka "gin_demo/third_party/kafka"
//)
//
//var log = logger.Log
//
//var UserActionChan = make(chan []byte)
//var kafkaWriter = kafka.NewDefaultWriter(kafka.GetKafkaBrokers(), "user_action")
//
////StartUserActionReporter 处理用户行为上报
//func StartUserActionReporter(parent context.Context) {
//	defer kafkaWriter.Close()
//	for i := 0; i < 4; i++ {
//		go func() {
//			select {
//			case data := <-UserActionChan:
//				if err := kafkaWriter.Produce(parent, nil, data); err != nil {
//					log.Errorf("produce error: %s\n", err.Error())
//					break
//				}
//			}
//		}()
//	}
//}
