package kafgo

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

const OFFSET_INIT string = "oldest"

var Client sarama.Client

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

// get kafka broker list
func GetKafkaBroker() []*sarama.Broker {
    brList := make([]*sarama.Broker, 0)
    // 获取controller信息    
    controller, ok := Client.Controller()
    if ok != nil {
        logs.Error("controller")
        return nil
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
    return brList
}