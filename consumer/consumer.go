package main

import (
    "fmt"
    "os"

    "github.com/optiopay/kafka/v2"	
)

func main() {

    if len(os.Args) < 3 {
        fmt.Printf("Usage: %s addr:port topic\n", os.Args[0])
        os.Exit(1)
    }
    addr := os.Args[1]
    topic := os.Args[2]

    // connect to kafka cluster
    // addresses := []string{"10.155.200.106:9092", "10.155.200.107:9092", "10.155.200.108:9092"}
    addresses := []string{addr}
    broker, err := kafka.Dial(addresses, kafka.NewBrokerConf("test"))
    if err != nil {
        panic(err)
    }
    defer broker.Close()
    
    resp, err := broker.Metadata()
    for i := 0; i < len(resp.Topics); i++ {
        println("resp topic: ", resp.Topics[i].Name)
        tname, _ := broker.PartitionCount(resp.Topics[i].Name)
        println("topic partition: ", tname)
    }

    count, err := broker.PartitionCount(topic)
    if err != nil {
        panic(err)
    }
    fmt.Println("topic: ", count, "count: ", count)

    // read all messages
    for j:=0;j<10;j++{		
        for i:=int32(0);i<count;i++ {
            
            conf := kafka.NewConsumerConf(topic, i)
            conf.StartOffset = kafka.StartOffsetNewest
            consumer, err := broker.Consumer(conf)
            if err != nil {
                panic(err)
            }
            msg, err := consumer.Consume()
            if err != nil {
                panic(err)
            }
            fmt.Println("Partition: ", msg.Partition, " Offset: ", msg.Offset, " Value: ", string(msg.Value))
            
        }					
    }

}


