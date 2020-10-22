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

func GetTopics() (dataList []interface{}, err error) {
	dataList = make([]interface{}, 0)
	topic := KafkaTopic{
		Topic:         "PACKET",
		PartitionSize: 10,
	}

	dataList = append(dataList, topic)
	return dataList, nil
}

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