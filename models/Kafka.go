package models

import (
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
}

type TopicMessages struct {
	Topic   string
	Message []MessageBlock
}

// 获取Topic消息
func GetTopicMessages(topic string) (topicMessage TopicMessages, err error) {
	var mb []MessageBlock
	mb = make([]MessageBlock, 0)
	msg := kafgo.GetKafkaMsg(topic)
	mb = append(mb, MessageBlock{
		PartitionId: 0,
		Txt:         "test1",
	})
	topicMessage = TopicMessages{
		Topic:   topic,
		Message: mb,
	}
	return topicMessage, nil
}

// GetTopics: 获取topic列表
func GetTopics() (dataList []interface{}, err error) {
	dataList = make([]interface{}, 0)
	topics := kafgo.GetKafkaTopic()
	for _, v := range topics {
		topic := new(KafkaTopic)
		topic.Topic = v.Name
		topic.PartitionSize = int32(len(v.Partitions))
		topic.AvailableCount = kafgo.GetTopicMsgNum(kafgo.Broker, topic.PartitionSize, topic.Topic)
		dataList = append(dataList, topic)	
    }
	
	return dataList, nil
}

// 获取broker列表
func GetBrokers() (dataList []interface{}, err error) {	
	dataList = make([]interface{}, 0)
	
	brList, controllerId := kafgo.GetKafkaBroker()
	for _, br := range brList {
		broker := &KafkaBroker{}
		broker.Id = br.ID()
		broker.Addr = br.Addr()
		if(broker.Id == controllerId) {
			broker.Controller = true
		}
		dataList = append(dataList, broker)
	}

	return dataList, nil
}