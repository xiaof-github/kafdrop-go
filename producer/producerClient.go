package main

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	//设置配置
    config := sarama.NewConfig()
    //等待服务器所有副本都保存成功后的响应
    config.Producer.RequiredAcks = sarama.WaitForAll
    //随机的分区类型
    config.Producer.Partitioner = sarama.NewRandomPartitioner
    //是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    //设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
    config.Version = sarama.V2_1_0_0

    //使用配置,新建一个异步生产者    
    producer, err := sarama.NewAsyncProducer([]string{"kafka-1:9092"}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()



	var enqueued, errors int
	
	for {
		select {
		case producer.Input() <- &sarama.ProducerMessage{Topic: "my_topic", Key: sarama.StringEncoder("test"), Value: sarama.StringEncoder("testing 123")}:
			enqueued++
		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			errors++
		default:			
			log.Printf("enqueued: %d", enqueued)
			time.Sleep(10*time.Second)
			break
		}
	}

	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)
}