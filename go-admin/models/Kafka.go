package models

type KafkaBroker struct {
	Id         int64  // id
	Addr       string // 地址
	Controller bool   // 是否controller
}

func GetBrokers() (dataList []interface{}, err error) {
	dataList = make([]interface{}, 10)
	broker := KafkaBroker{
		Id:         1,
		Addr:       "192.168.1.1",
		Controller: true,
	}

	dataList = append(dataList, broker)
	return dataList, nil
}