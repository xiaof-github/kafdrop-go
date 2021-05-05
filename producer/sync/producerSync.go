package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

// 封装发送消息C
func sendCMsg(msg *sarama.ProducerMessage, producer sarama.SyncProducer, i int) {
	value := "{\"app\":\"SDFA\",\"app_version\":\"V1.1.0\""
	
	value = value + ",\"id\":\"" + strconv.Itoa(i) + "\"}"
	msgTopic := "MSG_EXAMPLE"
	fmt.Println("msgTopic = ",msgTopic,",value = ",value)
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
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForLocal
	// 随机的分区类型：返回一个分区器，该分区器每次轮询一个可用的分区
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"10.155.200.105:9092"}, config)
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
	for {		
		sendCMsg(msg, producer, i)
		i++
		time.Sleep(60*time.Second)
	}
}