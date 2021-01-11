package main

import (
	"github.com/optiopay/kafka/v2"
)

func main() {

	// connect to kafka cluster
	addresses := []string{"10.155.200.106:9092", "10.155.200.107:9092", "10.155.200.108:9092"}
	broker, err := kafka.Dial(addresses, kafka.NewBrokerConf("test"))
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// create new consumer
	// conf := kafka.NewConsumerConf("MSG_EXAMPLE", 0)
	// conf.StartOffset = kafka.StartOffsetOldest
	// consumer, err := broker.Consumer(conf)
	// if err != nil {
	// 	panic(err)
	// }

	for i := int32(0); i < 20; i++ {
		offsetEarliest, _ := broker.OffsetEarliest("MSG_EXAMPLE", i)
		offsetLatest, _ := broker.OffsetLatest("MSG_EXAMPLE", i)
		println("partition: ", i, " start: ", offsetEarliest, " end: ", offsetLatest)
	}
	
	// read all messages
	// for {
	// 	msg, err := consumer.Consume()
	// 	if err != nil {
	// 		if err == kafka.ErrNoData {
	// 			break
	// 		}
	// 		panic(err)
	// 	}

	// 	fmt.Printf("message: %#v", msg)
	// }

	

}


