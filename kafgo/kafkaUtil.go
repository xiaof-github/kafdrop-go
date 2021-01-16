package kafgo

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/optiopay/kafka/v2"
	"github.com/optiopay/kafka/v2/proto"
)

const OFFSET_INIT string = "oldest"
// const CONSUME_ONE_PARTITION int = 10
const CONSUME_TOTAL_MESSAGE_NUM int = 200

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
func GetKafkaMsg(topic string) (map[int][]*sarama.ConsumerMessage, int) {	
    // var mapblock map[int][]*sarama.ConsumerMessage
    // 获取分区数量	
    partitionsNum,ok := TopicPartition[topic]
    if (ok) {
        fmt.Printf("topic: %s, partition size: %d\n", topic, partitionsNum)
    } else {
        logs.Error("don't have this topic", topic)
        return make(map[int][]*sarama.ConsumerMessage), 0
    }
    // // 取分区可消费的起始偏移量
    // offsr1 := sarama.OffsetRequest{
    //     Version: 1,		
    // }
    // for i:=0;i<partitionsNum;i++{
    //     offsr1.AddBlock(topic, int32(i), sarama.OffsetOldest, 999999999)
    // }
    // res1, err1 := Broker.GetAvailableOffsets(&offsr1)
    // if err1 != nil {
    //     panic("broker offset error")
    // }	

    // // 取分区可消费的终止偏移量
    // offsr2 := sarama.OffsetRequest{
    //     Version: 1,
    // }
    // for i:=0;i<partitionsNum;i++{
    //     offsr2.AddBlock(topic, int32(i), sarama.OffsetNewest, 999999999)
    // }
    // res2, err1 := Broker.GetAvailableOffsets(&offsr2)
    // if err1 != nil {
    //     panic("broker offset error")
    // }	

    // /**
    //  * 记录所有分区可消费记录数
    //  * 排序+总数
    //  */
    // var minOffset,totalOffset int64 = 0,0
    // // consumed := int64(0)
    // offsetSlice := make([]int64,0)
    // result := make([]int, 0)
    // for i:=0;i<partitionsNum;i++ {
    //     block1 := res1.GetBlock(topic, int32(i))
    //     block2 := res2.GetBlock(topic, int32(i))
    //     offsetSlice = append(offsetSlice, block2.Offset - block1.Offset)
    //     totalOffset += block2.Offset - block1.Offset
    //     if (minOffset > block2.Offset - block1.Offset) {
    //         minOffset = block2.Offset - block1.Offset
    //     }
    //     result = append(result, i)
    //     logs.Info("block1.Offset: %d, block2.Offset: %d, count: %d, partition: %d", 
    //         block1.Offset, block2.Offset, block2.Offset - block1.Offset, i)
    // }	
    
    /**
     * 根据所有分区可消费记录数，确定消费策略，开始消费
     * 策略1：可消费分区每个分区足够消费，平均消费最近的消息
     * 策略2：可消费分区累加不到消费消息数，每个分区消费最大消息数，直到每个分区消费完
     * 策略3：可消费分区累加超过消费消息数，不满足策略1条件。对分区可消费数排序，从最小可消费分区开始消费最大消息数，直到满足总消费消息数
     */
    // if (minOffset >= int64(CONSUME_TOTAL_MESSAGE_NUM/partitionsNum)) {
    //     // 策略1
    //     mapblock = getTopicMsg1(topic, partitionsNum, res2)
    // } else if (totalOffset < int64(CONSUME_TOTAL_MESSAGE_NUM)) {
    //     // 策略2
    //     mapblock = getTopicMsg2(topic, partitionsNum, res1, res2)
    // } else {
    //     // 策略3
    //     // 记录从小到大排序后的分区号
    //     sortMaoPao(offsetSlice, result)
    //     logs.Info("order result: ", result)
    //     mapblock = getTopicMsg3(topic, result, res1, res2)
    // }	

	// return mapblock, partitionsNum
	return make(map[int][]*sarama.ConsumerMessage), 0
}

/**
 * 策略1方式
 *
 */
// func getTopicMsg1(topic string, partitionsNum int, res2 *sarama.OffsetResponse) map[int][]*sarama.ConsumerMessage {
//     var consumed int;

//     mapblock := make(map[int][]*sarama.ConsumerMessage)

//     // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
//     for i:=0;i<partitionsNum;i++ {
//         block2 := res2.GetBlock(topic, int32(i))
//         partitionConsumer, err := Consumer.ConsumePartition(topic, int32(i), block2.Offset - int64(CONSUME_TOTAL_MESSAGE_NUM/partitionsNum))
//         if err != nil {
//             panic(err)    
//         }
//         consumed = 0
//         for {
//             select {
//                 case msg := <-partitionConsumer.Messages():
//                     log.Printf("Consumed message partition: %d, offset: %d, key: %s, value: %s\n", i, msg.Offset, msg.Key, msg.Value)
//                     consumed++;
//                     mapblock[i] = append(mapblock[i], msg)
//                     // time.Sleep(time.Second)
//                 default :
//                     log.Printf("consumed: %d", consumed)
//                     // time.Sleep(time.Second)
//             }
//             if (consumed >= CONSUME_TOTAL_MESSAGE_NUM/partitionsNum) {
//                 log.Printf("partition: %d finished, consumed: %d", i, consumed)
//                 break
//             }                
//         }
//     }
//     return mapblock
// }

/**
 * 策略2
 *
 */
//  func getTopicMsg2(topic string, partitionsNum int, res1 *sarama.OffsetResponse, res2 *sarama.OffsetResponse) map[int][]*sarama.ConsumerMessage {
//     var consumed int;

//     mapblock := make(map[int][]*sarama.ConsumerMessage)

//     // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
//     for i:=0;i<partitionsNum;i++ {
//         block1 := res1.GetBlock(topic, int32(i))
//         block2 := res2.GetBlock(topic, int32(i))
//         if (block1.Offset == block2.Offset){
//             continue
//         }
//         partitionConsumer, err := Consumer.ConsumePartition(topic, int32(i), block1.Offset)
//         if err != nil {
//             panic(err)    
//         }
//         consumed = 0
//         for {
//             select {
//                 case msg := <-partitionConsumer.Messages():
//                     log.Printf("Consumed message partition: %d, offset: %d, key: %s, value: %s\n", i, msg.Offset, msg.Key, msg.Value)
//                     consumed++;
//                     mapblock[i] = append(mapblock[i], msg)
//                     // time.Sleep(time.Second)
//                 default :
//                     log.Printf("consumed: %d", consumed)
//                     // time.Sleep(time.Second)
//             }
//             if (consumed >= int(block2.Offset - block1.Offset)) {
//                 log.Printf("partition: %d finished, consumed: %d", i, consumed)
//                 break
//             }                
//         }
//     }
//     return mapblock
// }

/**
 * 策略3 
 */
// func getTopicMsg3(topic string, result []int, res1 *sarama.OffsetResponse, res2 *sarama.OffsetResponse) map[int][]*sarama.ConsumerMessage {
//     var consumed int;
//     var total int = 0;

//     mapblock := make(map[int][]*sarama.ConsumerMessage)

//     // 返回最多CONSUME_TOTAL_MESSAGE_NUM条消息
//     for i:=0;i<len(result);i++ {
//         block1 := res1.GetBlock(topic, int32(result[i]))
//         block2 := res2.GetBlock(topic, int32(result[i]))
//         if (block1.Offset == block2.Offset){
//             continue
//         }
//         partitionConsumer, err := Consumer.ConsumePartition(topic, int32(i), block1.Offset)
//         if err != nil {
//             panic(err)    
//         }
//         consumed = 0
//         for {
//             select {
//                 case msg := <-partitionConsumer.Messages():
//                     log.Printf("Consumed message partition: %d, offset: %d, key: %s, value: %s\n", i, msg.Offset, msg.Key, msg.Value)
//                     consumed++;
//                     total++;
//                     mapblock[i] = append(mapblock[i], msg)
//                     // time.Sleep(time.Second)                    
//                 default :
//                     log.Printf("consumed: %d", consumed)
//                     // time.Sleep(time.Second)
//             }
//             if ((consumed >= int(block2.Offset - block1.Offset)) || total >= CONSUME_TOTAL_MESSAGE_NUM) {
//                 log.Printf("partition: %d finished, consumed: %d", i, consumed)
//                 break
//             }                
//         }
//     }
//     return mapblock
// }

/**
 * num：按分区号0,1,2...顺序记录每个分区可消费消息数；排序后按result中分区号记录分区可消费数
 * result：返回排序后分区号，和num中的消息数对应
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