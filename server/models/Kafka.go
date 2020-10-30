package models

type KafkaBroker struct {
	Id         int64  // id
	Addr       string // 地址
	Controller bool   // 是否controller
}

type KafkaTopic struct {
	Topic         string
	PartitionSize int32
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
	mb = append(mb, MessageBlock{
		PartitionId: 0,
		Txt:         "test",
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
	topic := KafkaTopic{
		Topic:         "PACKET",
		PartitionSize: 10,
	}

	dataList = append(dataList, topic)
	return dataList, nil
}

// 获取broker列表
func GetBrokers() (dataList []interface{}, err error) {
	dataList = make([]interface{}, 0)
	broker := KafkaBroker{
		Id:         1,
		Addr:       "192.168.1.1",
		Controller: true,
	}

	dataList = append(dataList, broker)
	return dataList, nil
}