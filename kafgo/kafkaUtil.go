package kafgo

import (
	"github.com/astaxie/beego/logs"
	"github.com/optiopay/kafka/v2"
	"github.com/optiopay/kafka/v2/proto"
)

const OFFSET_INIT string = "oldest"
// const CONSUME_ONE_PARTITION int = 10
const CONSUME_TOTAL_MESSAGE_NUM int32 = 200

var Broker *kafka.Broker

// 记录每个topic分区的大小
var TopicPartition map[string]int32 = make(map[string]int32)

// get kafka broker list and controller id
func GetKafkaBroker() ([]proto.MetadataRespBroker, int32) {
	resp, err := Broker.Metadata()
	if err != nil {
		panic(err)
	}
    return resp.Brokers, resp.ControllerID
}

// get kafka topic list
func GetKafkaTopic () ([]proto.MetadataRespTopic) {
	resp, err := Broker.Metadata()
	if err != nil {
		panic(err)
	}
	return resp.Topics
}

// get topic available msg count
func GetTopicMsgNum(topic string, partitionSize int32) int64 {
	var sum int64
	var count int64
	
	for i:=int32(0); i<partitionSize; i++{
		offsetEarliest, _ := Broker.OffsetEarliest(topic, i)
	    offsetLatest, _ := Broker.OffsetLatest(topic, i)
		count = offsetLatest - offsetEarliest
		// logs.Info("partition: %d, msg num: %d", i, count)
		sum += count
	}
    
    return sum
}

// get kafka topic msg
func GetKafkaMsg(topic string) (map[int32][]*proto.Message, int32) {
	// 获取分区个数
	var mapblock map[int32][]*proto.Message
	partitionCount, err := Broker.PartitionCount(topic)    
    if (err == nil) {
        logs.Info("topic: %s, partition size: %d\n", topic, partitionCount)
    } else {
        logs.Error("don't have this topic: ", topic)
        return make(map[int32][]*proto.Message), 0
    }

    /**
     * 记录所有分区可消费记录数
     * 排序+总数
     */
    var minOffset,totalOffset int64 = 0,0
    // consumed := int64(0)
    offsetSlice := make([]int64,0)
	result := make([]int32, 0)
	offsetEarliest, _ := Broker.OffsetEarliest(topic, 0)
	offsetLatest, _ := Broker.OffsetLatest(topic, 0)
	minOffset = offsetLatest - offsetEarliest
    for i := int32(0); i < partitionCount; i++ {
		offsetEarliest, _ := Broker.OffsetEarliest(topic, i)
		offsetLatest, _ := Broker.OffsetLatest(topic, i)
        offsetSlice = append(offsetSlice, offsetLatest - offsetEarliest)
        totalOffset += offsetLatest - offsetEarliest
        if (minOffset > offsetLatest - offsetEarliest) {
            minOffset = offsetLatest - offsetEarliest
        }
        result = append(result, i)
        logs.Info("partition: %d, diff: %d", i, offsetLatest - offsetEarliest)
	}	
	logs.Info("minOffset: %d, totalOffset: %d", minOffset, totalOffset)
    
    /**
     * 根据所有分区可消费记录数，开始消费
     * 方式1：每个分区的消息足够多，消费每个分区最近的消息
     * 方式2：每个分区加起来消息不够多，消费每个分区所有的消息
     * 方式3：消息分布不均匀，排序后，从消息最少的分区开始消费
     */
    if (minOffset >= int64(CONSUME_TOTAL_MESSAGE_NUM/partitionCount)) {        
        mapblock = getTopicMsg1(topic, partitionCount)
    } else if (totalOffset < int64(CONSUME_TOTAL_MESSAGE_NUM)) {        
        mapblock = getTopicMsg2(topic, partitionCount)
    } else {        
        // 记录从小到大排序后的分区号
        sortMaoPao(offsetSlice, result)
        logs.Info("order result: ", result)
        mapblock = getTopicMsg3(topic, partitionCount, result)
    }	

	return mapblock, partitionCount	
}

/**
 * 方式1
 *
 */
func getTopicMsg1(topic string, partitionCount int32) map[int32][]*proto.Message {    
    mapblock := make(map[int32][]*proto.Message)

    // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
    for i:=int32(0);i<partitionCount;i++ {
        // create new consumer
		conf := kafka.NewConsumerConf(topic, i)
		offsetLatest, _ := Broker.OffsetLatest(topic, i)
		conf.StartOffset = offsetLatest - int64(CONSUME_TOTAL_MESSAGE_NUM/partitionCount)
		consumer, err := Broker.Consumer(conf)
		if err != nil {
			panic(err)
		}
		
		for j:=int32(0);j<CONSUME_TOTAL_MESSAGE_NUM/partitionCount;j++{
			msg, err := consumer.Consume()
			if err != nil {
				if err == kafka.ErrNoData {
					break
				}
				panic(err)
			}
			// str := string(msg.Value)			
			mapblock[i] = append(mapblock[i], msg)			
		}
    }
    return mapblock
}

/**
 * 方式2
 *
 */
func getTopicMsg2(topic string, partitionCount int32) map[int32][]*proto.Message {    
    mapblock := make(map[int32][]*proto.Message)

    // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
    for i:=int32(0);i<partitionCount;i++ {
        // create new consumer
		conf := kafka.NewConsumerConf(topic, i)		
		offsetEarliest, _ := Broker.OffsetEarliest(topic, i)
		offsetLatest, _ := Broker.OffsetLatest(topic, i)
		conf.StartOffset = offsetEarliest
		consumer, err := Broker.Consumer(conf)
		if err != nil {
			panic(err)
		}
		
		for j:=int32(0);j<int32(offsetLatest - offsetEarliest);j++{
			msg, err := consumer.Consume()
			if err != nil {
				if err == kafka.ErrNoData {
					break
				}
				panic(err)
			}
			// str := string(msg.Value)			
			mapblock[i] = append(mapblock[i], msg)			
		}
    }
    return mapblock
}

/**
 * 方式3
 *
 */
 func getTopicMsg3(topic string, partitionCount int32, result []int32) map[int32][]*proto.Message {    
    mapblock := make(map[int32][]*proto.Message)
	var total int32

    // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
    for i:=int32(0);i<partitionCount;i++ {
        // create new consumer
		conf := kafka.NewConsumerConf(topic, result[i])		
		offsetEarliest, _ := Broker.OffsetEarliest(topic, result[i])
		offsetLatest, _ := Broker.OffsetLatest(topic, result[i])
		conf.StartOffset = offsetEarliest
		consumer, err := Broker.Consumer(conf)
		if err != nil {
			panic(err)
		}
		
		for j:=int32(0);j<int32(offsetLatest - offsetEarliest);j++{			
			msg, err := consumer.Consume()
			if err != nil {
				if err == kafka.ErrNoData {
					break
				}
				panic(err)
			}
			total++;
			// str := string(msg.Value)			
			mapblock[i] = append(mapblock[i], msg)
			if (total > CONSUME_TOTAL_MESSAGE_NUM){
				break
			}
		}
    }
    return mapblock
}

/**
 * num：按分区号0,1,2...顺序记录每个分区可消费消息数；排序后按result中分区号记录分区可消费数
 * result：返回排序后分区号，和num中的消息数对应
 */
func sortMaoPao(num []int64, result []int32){
    var tmp int64
    var tmp1 int32
    for i:=int32(len(num)-1);i>0;i-- {
        for j:=int32(0);j<i;j++ {
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