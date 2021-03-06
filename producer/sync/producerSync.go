package main

import (
    "fmt"    
    "time"
    "strconv"
    "os"	

    "github.com/Shopify/sarama"
)

// 封装发送消息C
func sendCMsg(msg *sarama.ProducerMessage, producer sarama.SyncProducer, msgTopic string, data string, i int) {
    value := data
    fmt.Println("msgTopic = ",msgTopic,",value = ",data, "number: ", i)
    msg.Topic = msgTopic
    //将字符串转换为字节数组
    msg.Value = sarama.ByteEncoder(value)
    //fmt.Println(value)
    //SendMessage：该方法是生产者生产给定的消息
    //生产成功的时候返回该消息的分区和所在的偏移量
    //生产失败的时候返回error
    partition1, offset1, err1 := producer.SendMessage(msg)

    if err1 != nil {
        fmt.Println("Send message Fail")
    }
    fmt.Printf("Partition = %d, offset=%d\n", partition1, offset1)
}



func main() {
    if len(os.Args) < 5 {
        fmt.Printf("Usage: %s addr:port topic num msg\n", os.Args[0])
        os.Exit(1)
    }
    addr := os.Args[1]
    topic := os.Args[2]
    num := os.Args[3]
    data := os.Args[4]

    config := sarama.NewConfig()
    // 等待服务器所有副本都保存成功后的响应
    config.Producer.RequiredAcks = sarama.WaitForLocal
    // 随机的分区类型：返回一个分区器，该分区器每次轮询一个可用的分区
    config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
    // 是否等待成功和失败后的响应
    config.Producer.Return.Successes = true

    // 使用给定代理地址和配置创建一个同步生产者
    producer, err := sarama.NewSyncProducer([]string{addr}, config)
    // producer, err := sarama.NewSyncProducer([]string{"10.155.200.106:9092","10.155.200.107:9092","10.155.200.108:9092"}, config)
    if err != nil {
        panic(err)
    }

    defer producer.Close()

    //构建发送的消息，
    msg := &sarama.ProducerMessage {
        //Topic: "test",//包含了消息的主题
        Partition: int32(20),//
        Key:        nil,//
    }
    var i int = 0
    count,err := strconv.Atoi(num)
    if err != nil {
        panic(err)
    }
    for j:=0;j<count;j++{		
        sendCMsg(msg, producer, topic, data, i)
        i++
        time.Sleep(3*time.Second)
    }
}