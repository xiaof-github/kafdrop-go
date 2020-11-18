package kafgo

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"github.com/xiaof-github/kafdrop-go/models"
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
func GetKafkaBroker() ([]*models.KafkaBroker) {
    brList := make([]*models.KafkaBroker, 0)
    // 获取controller信息
    bc := new(models.KafkaBroker)
    brk, ok := Client.Controller()
    if ok != nil {
        logs.Error("controller")
        return nil
    }
    bc.Addr = brk.Addr()
    bc.Controller = true
    bc.Id = brk.ID()
    brList = append(brList, bc)
    // 获取broker节点信息    
    brokers := Client.Brokers()
    for _, brk := range brokers {
        if brk.Addr() == bc.Addr {
            continue
        }
        kb := new(models.KafkaBroker)
        kb.Addr = brk.Addr()
        kb.Controller = false
        kb.Id = brk.ID()
        brList = append(brList, kb)
    }  

    return brList
}