package kafgo

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

const OFFSET_INIT string = "oldest"

var Client sarama.Client
var Broker *sarama.Broker
var Consumer sarama.Consumer
var TopicPartiton map[string]int

// get kafka client
func GetClient(addrs []string, version string, offInit string) (sarama.Client, error) {
    config := sarama.NewConfig()
    ver, err := sarama.ParseKafkaVersion(version)
    if err != nil {
        logs.Error("Error parsing Kafka version: %v", err)
        panic("version error")
    }
    config.Version = ver
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
    if offInit == OFFSET_INIT {
        config.Consumer.Offsets.Initial = sarama.OffsetOldest
    } else {
        config.Consumer.Offsets.Initial = sarama.OffsetNewest
    }

    client, err := sarama.NewClient(addrs, config)
    if err != nil {
        panic("kafka client create error")
    }
    return client, err
}

// get kafka broker list and controller id
func GetKafkaBroker() ([]*sarama.Broker, int32) {
    brList := make([]*sarama.Broker, 0)
    // 获取controller信息    
    controller, ok := Client.Controller()
    if ok != nil {
        logs.Error("controller")
        return brList, -1
    }    
    brList = append(brList, controller)
    // 获取broker节点信息，去除controller
    brokers := Client.Brokers()
    for _, br := range brokers {
        if br.Addr() == controller.Addr() {
            continue
        }        
        brList = append(brList, br)
    }
    return brList, controller.ID()
}

// get kafka topic list
func GetKafkaTopic () ([]*sarama.TopicMetadata) {
    request := sarama.MetadataRequest{ /*Topics: []string{"abba"}*/ }
    response, err := Broker.GetMetadata(&request)
    if err != nil {
        _ = Broker.Close()
        panic(err)
    }

    return response.Topics
}

// get topic available msg count
func GetTopicMsgNum(broker *sarama.Broker, partitionSize int32, topic string) int64 {
    var i int32
    var sum int64
    // 当前Topic, partition可消费的最小偏移量
    offsr := sarama.OffsetRequest{
        Version: 1,		
    }
    offsrEnd := sarama.OffsetRequest{
        Version: 1,		
    }
    len := partitionSize
    
    for i=0;i<len;i++{
        offsr.AddBlock(topic, i, sarama.OffsetOldest, 999999999)		
        offsrEnd.AddBlock(topic, i, sarama.OffsetNewest, 999999999)
    }

    // offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
    res1, err1 := broker.GetAvailableOffsets(&offsr)
    if err1 != nil {
        panic("broker offset error")
    }
    res2, err2 := broker.GetAvailableOffsets(&offsrEnd)
    if err2 != nil {
        panic("broker offsetEnd error")
    }
    for i=0;i<len;i++{
        r2 := res2.GetBlock(topic, i)
        r1 := res1.GetBlock(topic, i)
        sum += r2.Offset-r1.Offset
    }        
    
    return sum
}

// get kafka topic msg
func GetKafkaMsg(topic string) (map[int][]*sarama.ConsumerMessage, int) {
	/* 
	 * 从第一个分区开始消费，消费到200条消息停止；如果所有分区消费完之后，没有200条消息，则退出
	 */ 
	// 分区数量
	partitionsNum,ok := TopicPartiton[topic]
	if (ok) {
		fmt.Println("topic: %s, partition size: %d", topic, partitionsNum)
	} else {
		logs.Error("don't have this topic", topic)
		return make(map[int][]*sarama.ConsumerMessage), 0
	}
	// 取第一个分区最近可消费的偏移量
	offsr1 := sarama.OffsetRequest{
        Version: 1,		
	}
	for i:=0;i<partitionsNum;i++{
        offsr1.AddBlock(topic, int32(i), sarama.OffsetOldest, 999999999)
	}
	res1, err1 := Broker.GetAvailableOffsets(&offsr1)
    if err1 != nil {
        panic("broker offset error")
	}	

	// 取最新可消费的偏移量
	offsr2 := sarama.OffsetRequest{
        Version: 1,
	}
	for i:=0;i<partitionsNum;i++{
        offsr2.AddBlock(topic, int32(i), sarama.OffsetNewest, 999999999)
	}
	res2, err1 := Broker.GetAvailableOffsets(&offsr2)
    if err1 != nil {
        panic("broker offset error")
	}
	mblock := make(map[int][]*sarama.ConsumerMessage)

	var offset int64
	for i:=0;i<partitionsNum;i++ {
		block1 := res1.GetBlock(topic, int32(i))
		block2 := res2.GetBlock(topic, int32(i))
		len := block2.Offset - block1.Offset
		if (len>int64(200/partitionsNum)) {
			offset = block2.Offset - int64(200/partitionsNum)
		} else {
			offset = block1.Offset
		}
		partitionConsumer, err := Consumer.ConsumePartition(topic, int32(i), offset)
		if err != nil {
			panic(err)    
		}
		consumed := int64(0)
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("Consumed message partition %d\n offset %d\n key %s\n value %s\n", i, msg.Offset, msg.Key, msg.Value)
				consumed++;
				mblock[i] = append(mblock[i], msg)
				// time.Sleep(time.Second)            
			default :
				log.Printf("consumed: %d", consumed)				
			}
			if (consumed >= int64(200/partitionsNum)) {
				log.Printf("partition: %d, consumed: %d", i, consumed)
				break
			}                
		}


	}
	

	// 取第n个分区最近可消费的偏移量

	// 从第一个分区开始消费，使用协程，多个分区消费的个数进行累加

	// 返回200条消息

	return mblock, partitionsNum
}
