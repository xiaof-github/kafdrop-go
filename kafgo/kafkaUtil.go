package kafgo

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

const OFFSET_INIT string = "oldest"

var Client sarama.Client
var Broker *sarama.Broker

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
func GetKafkaMsg(topic string) map[int][]string {
	/* 
	 * 从第一个分区开始消费，消费到200条消息停止；如果没有200条消息，接着消费下一个分区的消息，
	 * 直到消费满200条消息为止	 
	 */ 
	return make(map[int][]string, 20)
}
