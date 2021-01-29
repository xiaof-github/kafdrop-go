package main

import (
	"fmt"

	"github.com/optiopay/kafka/v2"
)

func main() {

	// connect to kafka cluster
	// addresses := []string{"10.155.200.106:9092", "10.155.200.107:9092", "10.155.200.108:9092"}
	addresses := []string{"152.136.200.213:9092"}
	broker, err := kafka.Dial(addresses, kafka.NewBrokerConf("test"))
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// create new consumer
	conf := kafka.NewConsumerConf("MSG_EXAMPLE", 0)
	conf.StartOffset = kafka.StartOffsetOldest
	consumer, err := broker.Consumer(conf)
	if err != nil {
		panic(err)
	}
	
	resp, err := broker.Metadata()
	for i := 0; i < len(resp.Topics); i++ {
		println("resp topic: ", resp.Topics[i].Name)
		tname, _ := broker.PartitionCount(resp.Topics[i].Name)
		println("topic partition: ", tname)
	}

	// for i := 0; i < len(resp.Brokers); i++ {
	// 	println("resp broker: ", resp.Brokers[i].NodeID, resp.Brokers[i].Host, resp.Brokers[i].Port)
	// 	println("resp controller id: ", resp.ControllerID)
	// }
	

	// get offsets
	// for i := int32(0); i < 20; i++ {
	// 	offsetEarliest, _ := broker.OffsetEarliest("MSG_EXAMPLE", i)
	// 	offsetLatest, _ := broker.OffsetLatest("MSG_EXAMPLE", i)
	// 	println("partition: ", i, " start: ", offsetEarliest, " end: ", offsetLatest)
	// }


	
	// read all messages
	for {
		msg, err := consumer.Consume()
		if err != nil {
			if err == kafka.ErrNoData {
				break
			}
			panic(err)
		}
		str := string(msg.Value)
		fmt.Println("Partition: ", msg.Partition, " Offset: ", msg.Offset, " Value: ", str)
	}

	

}


