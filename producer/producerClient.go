package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	//设置配置
    config := sarama.NewConfig()
    //等待服务器所有副本都保存成功后的响应
    config.Producer.RequiredAcks = sarama.NoResponse
    //随机的分区类型
    config.Producer.Partitioner = sarama.NewRandomPartitioner
    //是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
    config.Producer.Return.Successes = true
    config.Producer.Return.Errors = true
    //设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
    config.Version = sarama.V2_1_0_0

    //使用配置,新建一个异步生产者    
    producer, err := sarama.NewAsyncProducer([]string{"10.155.200.106:9092"}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			errors++
		}
	}()

	ProducerLoop:
	for {
		message := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder("testing 456")}
		select {
		case producer.Input() <- message:
			enqueued++

		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break ProducerLoop
		}
	}

	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, errors)
}