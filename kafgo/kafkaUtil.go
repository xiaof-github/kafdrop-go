package kafgo

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

const OFFSET_INIT string = "oldest"
const CONSUME_NUM int = 200

var Client sarama.Client
var Broker *sarama.Broker
var Consumer sarama.Consumer
// 记录每个topic分区的大小
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
	// 获取分区数量	
	partitionsNum,ok := TopicPartiton[topic]
	if (ok) {
		fmt.Printf("topic: %s, partition size: %d\n", topic, partitionsNum)
	} else {
		logs.Error("don't have this topic", topic)
		return make(map[int][]*sarama.ConsumerMessage), 0
	}
	// 取分区可消费的起始偏移量
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

	// 取分区可消费的终止偏移量
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

	/**
	 * 记录所有分区可消费记录数
	 * 排序+总数
	 */
	var minOffset,totalOffset int64 = 0,0
	// consumed := int64(0)
	offsetSlice := make([]int64,0)
	result := make([]int, 0)
	for i:=0;i<partitionsNum;i++ {
		block1 := res1.GetBlock(topic, int32(i))
		block2 := res2.GetBlock(topic, int32(i))
		offsetSlice = append(offsetSlice, block2.Offset - block1.Offset)
		totalOffset += block2.Offset - block1.Offset
		if (minOffset > block2.Offset - block1.Offset) {
			minOffset = block2.Offset - block1.Offset
		}
		result = append(result, i)
		logs.Info("block1.Offset: %d, block2.Offset: %d, count: %d, partition: %d", 
			block1.Offset, block2.Offset, block2.Offset - block1.Offset, i)
	}
	// 记录从小到大排序后的分区号
	sortMaoPao(offsetSlice, result)
	logs.Info("order result: ", result)
	
	/**
	 * 根据所有分区可消费记录数，确定消费策略，开始消费
	 * 策略1：可消费分区每个分区足够消费，平均消费
	 * 策略2：可消费分区累加不到200条消息，每个分区消费最大消息数，直到每个分区消费完
	 * 策略3：可消费分区累加到200条消息，不满足策略1条件，从最大可消费分区开始消费，直到满200条，需要
	 * 排序，从最少的分区开始消费
	 */
	

	// 返回200条消息
	// partitionConsumer, err := Consumer.ConsumePartition(topic, int32(i), offset)
	// if err != nil {
	// 	panic(err)    
	// }		
	// for {
	// 	select {
	// 	case msg := <-partitionConsumer.Messages():
	// 		log.Printf("Consumed message partition: %d, offset: %d, key: %s, value: %s\n", i, msg.Offset, msg.Key, msg.Value)
	// 		consumed++;
	// 		mblock[i] = append(mblock[i], msg)
	// 		time.Sleep(time.Second)
	// 	default :
	// 		log.Printf("consumed: %d", consumed)
	// 		time.Sleep(time.Second)
	// 	}
	// 	if (consumed >= len) {
	// 		log.Printf("partition: %d, consumed: %d", i, consumed)
	// 		break
	// 	}                
	// }

	return mblock, partitionsNum
}
/**
 * num是排序后数组内容，result记录初始下标，返回排序后下标
 */
func sortMaoPao(num []int64, result []int){
	var tmp int64
	var tmp1 int
	for i:=len(num)-1;i>0;i-- {
		for j:=0;j<i;j++ {
			if (num[j] > num[j+1]){
				// 值交换
				tmp = num[j]
				num[j] = num[j+1]
				num[j+1] = tmp
				// 下标交换
				tmp1 = result[j]
				result[j] = result[j+1]
				result[j+1] = tmp1
			}
		}
	}
}