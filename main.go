package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/xiaof-github/kafdrop-go/kafgo"
	_ "github.com/xiaof-github/kafdrop-go/routers"
)

func main() {	
    viper.SetConfigName("app")
    viper.AddConfigPath("conf")
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
    addr := viper.GetString("kafka.addr")
    version := viper.GetString("kafka.version")
    offsetsInitial := viper.GetString("kafka.offsetsInitial")

    addrs := strings.Split(addr, ",")	
    fmt.Println(addrs)
    fmt.Println(kafgo.OFFSET_INIT)
    // 打开调试信息
    sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
    
    // kafka客户端建立连接
    client, err := kafgo.GetClient(addrs, version, offsetsInitial)
    if err != nil {
        panic(fmt.Errorf("Fatal error get client: %s \n", err))
    }
    defer client.Close()
    kafgo.Client = client

    // kafka broker
    kafkaHost := addrs[0]
    broker := sarama.NewBroker(kafkaHost)
    config := initConfig(version, offsetsInitial)
    err1 := broker.Open(config)
    if err1 != nil {
        panic(err1)
	}
	defer broker.Close()
	kafgo.Broker = broker	

	// kafka Consumer
	consumer, err := sarama.NewConsumer(addrs, nil)
    if err != nil {
        panic(err)
	}
	kafgo.Consumer = consumer

	// 缓存topic分区
	topicPartiton := make(map[string]int)
	kafgo.TopicPartiton = topicPartiton

    //注册sqlite3
    orm.RegisterDataBase("default", "sqlite3", "go-admin.db")
    //同步 ORM 对象和数据库
    //这时, 在你重启应用的时候, beego 便会自动帮你创建数据库表。
    orm.Debug = true

    orm.RunSyncdb("default", false, true)

    //增加拦截器。
    var filterAdmin = func(ctx *context.Context) {
        url := ctx.Input.URL()
        logs.Info("##### filter url : %s", url)
        //TODO 如果判断用户未登录。

    }
    beego.InsertFilter("/admin/*", beego.BeforeExec, filterAdmin)

    beego.Run()
}

func initConfig(version string, offsetConfig string) *sarama.Config {
    config := sarama.NewConfig()
    ver, err := sarama.ParseKafkaVersion(version)
    if err != nil {
        logs.Error("Error parsing Kafka version: %v", err)
        panic("version error")
    }	
    config.Version = ver
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
    if (offsetConfig == kafgo.OFFSET_INIT) {
        config.Consumer.Offsets.Initial = sarama.OffsetOldest
    } else {
        config.Consumer.Offsets.Initial = sarama.OffsetNewest
    }	
    return config
}