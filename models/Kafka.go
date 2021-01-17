package models

import (
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/xiaof-github/kafdrop-go/kafgo"
)

type KafkaBroker struct {
	Id         int32  // id
	Addr       string // 地址
	Controller bool   // 是否controller
}

type KafkaTopic struct {
	Topic         string
	PartitionSize int32
	AvailableCount int64
}

type MessageBlock struct {
	PartitionId int32
	Txt         string
	Offset      int64
}

type TopicMessages struct {
	Topic   string
	Message []MessageBlock
}

// 获取Topic消息
func GetTopicMessages(topic string) (topicMessage TopicMessages, err error) {
	var mb []MessageBlock
	mb = make([]MessageBlock, 0)
	msg, partitionNum := kafgo.GetKafkaMsg(topic)	
	for i:=int32(0);i<partitionNum;i++ {
		for j:=0;j<len(msg[i]);j++ {
			mb = append(mb, MessageBlock{
				PartitionId: msg[i][j].Partition,
				Txt:         string(msg[i][j].Value),
				Offset:		 msg[i][j].Offset,
			})
		}
	}	
	topicMessage = TopicMessages{
		Topic:   topic,
		Message: mb,
	}
	return topicMessage, nil
}

// GetTopics: 获取topic列表
func GetTopics() (dataList []interface{}, err error) {		
	// 获取topic列表
	dataList = make([]interface{}, 0)
	topics := kafgo.GetKafkaTopic()	
	for _, v := range topics {
		topic := new(KafkaTopic)
		topic.Topic = v.Name
		count, err := kafgo.Broker.PartitionCount(v.Name)
		if err != nil {
			logs.Info("err :", err)
		}
		topic.PartitionSize = count
		topic.AvailableCount = kafgo.GetTopicMsgNum(v.Name, count)
		// 缓存topic分区
		kafgo.TopicPartition[v.Name] = count
		dataList = append(dataList, topic)	
    }
	
	return dataList, nil
}

// 获取broker列表
func GetBrokers() (dataList []interface{}, err error) {	
	dataList = make([]interface{}, 0)
	
	brokerL, controllerId := kafgo.GetKafkaBroker()
	for _, br := range brokerL {
		broker := &KafkaBroker{}
		broker.Id = br.NodeID
		broker.Addr = br.Host + ":" + strconv.Itoa(int(br.Port))
		if(broker.Id == controllerId) {
			broker.Controller = true
		}
		dataList = append(dataList, broker)
	}

	return dataList, nil
}