package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

// Sarama configuration options
var (
    version  = "2.1.1"
    group    = ""
    topics   = ""
    assignor = "range"
    oldest   = true
    verbose  = false
)

type PartitionOffsets struct {
    partition []string
    firstOffsets []int
}

// func init() {
//     flag.StringVar(&brokers, "brokers", "10.155.200.120:9092", "Kafka bootstrap brokers to connect to, as a comma separated list")
//     flag.StringVar(&group, "group", "sab", "Kafka consumer group definition")
//     flag.StringVar(&version, "version", "2.1.1", "Kafka cluster version")
//     flag.StringVar(&topics, "topics", "test", "Kafka topics to be consumed, as a comma separated list")
//     flag.StringVar(&assignor, "assignor", "range", "Consumer group partition assignment strategy (range, roundrobin, sticky)")
//     flag.BoolVar(&oldest, "oldest", true, "Kafka consumer consume initial offset from oldest")
//     flag.BoolVar(&verbose, "verbose", false, "Sarama logging")
//     flag.Parse()	

//     if len(brokers) == 0 {
//         panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
//     }

//     if len(topics) == 0 {
//         panic("no topics given to be consumed, please set the -topics flag")
//     }

//     if len(group) == 0 {
//         panic("no Kafka consumer group defined, please set the -group flag")
//     }
// }

func getTopicMsgNum(broker *sarama.Broker, partition map[string]int, topic string) int64 {
    var i int32
    var sum int64
    // 当前Topic, partition可消费的最小偏移量
    offsr := sarama.OffsetRequest{
        Version: 1,		
    }
    offsrEnd := sarama.OffsetRequest{
        Version: 1,		
    }
    len := partition[topic]
    
    for i=0;i<int32(len);i++{
        offsr.AddBlock(topic, i, sarama.OffsetOldest, 999999999)		
        offsrEnd.AddBlock(topic, i, sarama.OffsetNewest, 999999999)
    }

    // offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetNewest, 999999999)
    res1, err1 := broker.GetAvailableOffsets(&offsr)
    if err1 != nil {
        panic("broker offset error")
    }
    res2, err2 := broker.GetAvailableOffsets(&offsrEnd)
    if err2 != nil {
        panic("broker offsetEnd error")
    }
    for i=0;i<int32(len);i++{
        r2 := res2.GetBlock(topic, i)
        r1 := res1.GetBlock(topic, i)
        sum += r2.Offset-r1.Offset
    }        
    
    return sum
}

func consumeTopic(topic string, block *sarama.OffsetResponseBlock, partition int32) {
    partitionConsumer, err := consumer.ConsumePartition(topic, partition, block.Offset)
    if err != nil {
        panic(err)
    }

}

func consumeMsg(broker *sarama.Broker, addr []string, topic string, partitions map[string]int) ([]string, error) {
    consumer, err := sarama.NewConsumer(addr, nil)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err := consumer.Close(); err != nil {
            log.Fatalln(err)
        }
    }()

    offsr := sarama.OffsetRequest{
        Version: 1,		
    }

    len := partitions[topic]

    for i:=0;i<len;i++{
        offsr.AddBlock(topic, int32(i), sarama.OffsetOldest, 999999999)
    }

    // offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetNewest, 999999999)
    res1, err1 := broker.GetAvailableOffsets(&offsr)
    if err1 != nil {
        panic("broker offset error")
    }

    var block *sarama.OffsetResponseBlock
    for i:=0;i<len;i++ {
        block = res1.GetBlock(topic, int32(i))
        go consumeTopic(topic, block, int32(i))        
    }  
       

    defer func() {
        if err := partitionConsumer.Close(); err != nil {
            log.Fatalln(err)
        }
    }()

    

    log.Printf("Consumed: \n")
    return []string{"a", "b", "c"} , nil
}


func main() {
    log.Println("Starting kafdrop-go")	

    brokers := []string{"10.155.200.120:9092"}
    kafkaHost := brokers[0]

    if verbose {
        sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
    }

    version, err := sarama.ParseKafkaVersion(version)
    if err != nil {
        log.Panicf("Error parsing Kafka version: %v", err)
    }

    /**
     * Construct a new Sarama configuration.
     * The Kafka cluster version has to be defined before the consumer/producer is initialized.
     */
    config := sarama.NewConfig()
    config.Version = version

    switch assignor {
    case "sticky":
        config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
    case "roundrobin":
        config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    case "range":
        config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
    default:
        log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
    }

    if oldest {
        config.Consumer.Offsets.Initial = sarama.OffsetOldest
    }

    
    broker := sarama.NewBroker(kafkaHost)
    err1 := broker.Open(config)
    if err1 != nil {
        panic(err1)
    }
    // 当前topic和partition元数据
    request := sarama.MetadataRequest{ /*Topics: []string{"abba"}*/ }
    response, err := broker.GetMetadata(&request)
    if err != nil {
        _ = broker.Close()
        panic(err)
    }
    // resjson, _ := json.Marshal(response)
    // fmt.Println(string(resjson))
    

    //topics := make([]string, len(response.Topics))
    topics1 := []string{}
    partitions := make(map[string]int, 100)

    for _, v := range response.Topics {
        topics1 = append(topics1, v.Name)
        partitions[v.Name] = len(v.Partitions)
        // fmt.Printf("v.Name: %+v, v.Partitions: %d\n", v.Name, partitions[v.Name])
    }

    // 计算topic可消费消息总量
    sum := getTopicMsgNum(broker, partitions, "PACKET_DNS_RESPONSE")
    fmt.Printf("topic[%s], sum: %d\n", "PACKET_DNS_RESPONSE", sum)

    tp := map[string][]int32{
        "PACKET_DNS_RESPONSE":[]int32{0,1},
    }

    // ConsumerGroup已消费的Topic, partition的偏移量
    offr := sarama.OffsetFetchRequest{
        Version: 5,
        ConsumerGroup: "dns_response_slbsdc",
    }
    offr.AddPartition("PACKET_DNS_RESPONSE", 1)	

    res, err := broker.FetchOffset(&offr)
    if err != nil {
        panic("broker fetchoffset error")
    }
    fmt.Printf("res: %+v\n", res)
    v1 := res.GetBlock("PACKET_DNS_RESPONSE", 1)
    fmt.Printf("v1: %+v\n", v1)	

    // 当前Topic, partition可消费的最小偏移量
    offsr := sarama.OffsetRequest{
        Version: 1,		
    }

    offsr.AddBlock("PACKET_DNS_RESPONSE", 1, sarama.OffsetOldest, 999999999)
    offsr.AddBlock("PACKET_DNS_RESPONSE", 2, sarama.OffsetOldest, 999999999)
    offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetOldest, 999999999)
    // offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetNewest, 999999999)
    res1, err := broker.GetAvailableOffsets(&offsr)
    if err != nil {
        panic("broker offset error")
    }
    fmt.Printf("res: %+v\n", res1)
    o1 := res1.GetBlock("PACKET_DNS_RESPONSE", 1)
    fmt.Printf("o1: %+v\n", o1)
    o2 := res1.GetBlock("PACKET_DNS_RESPONSE", 2)
    fmt.Printf("o2: %+v\n", o2)
    o3 := res1.GetBlock("PACKET_DNS_RESPONSE", 3)
    fmt.Printf("o3: %+v\n", o3)

    // 当前Topic, partition可消费的最大偏移量
    offsr.AddBlock("PACKET_DNS_RESPONSE", 1, sarama.OffsetNewest, 999999999)
    offsr.AddBlock("PACKET_DNS_RESPONSE", 2, sarama.OffsetNewest, 999999999)
    offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetNewest, 999999999)
    // offsr.AddBlock("PACKET_DNS_RESPONSE", 3, sarama.OffsetNewest, 999999999)
    res1, err = broker.GetAvailableOffsets(&offsr)
    if err != nil {
        panic("broker offset error")
    }
    fmt.Printf("res: %+v\n", res1)
    n1 := res1.GetBlock("PACKET_DNS_RESPONSE", 1)
    fmt.Printf("n1: %+v\n", n1)
    n2 := res1.GetBlock("PACKET_DNS_RESPONSE", 2)
    fmt.Printf("n2: %+v\n", n2)
    n3 := res1.GetBlock("PACKET_DNS_RESPONSE", 3)
    fmt.Printf("n3: %+v\n", n3)

    // partition lag, lag = lastOffset - curOffset
    p1Lag := n1.Offset - v1.Offset
    fmt.Printf("Topic: PACKET_DNS_RESPONSE, Partition: 1, lag: %d\n", p1Lag)


    client, err := sarama.NewClient([]string{"10.155.200.120:9092"}, config)
    if err != nil {
        panic("client create error")
    }
    defer client.Close()
    // 获取主题的名称集合
    topics, err := client.Topics()
    if err != nil {
        panic("get topics err")
    }
    
    for _, tv1 := range topics {
        // fmt.Println(e)
        _ = tv1		
    }
    abc, e := client.GetOffset("PACKET_DNS_RESPONSE", 1, sarama.OffsetOldest)
    fmt.Printf("first offset: %+v,",abc)
    fmt.Println(e)
    abc, e = client.GetOffset("PACKET_DNS_RESPONSE", 1, sarama.OffsetNewest)
    fmt.Printf("last offset: %+v,",abc)
    fmt.Println(e)

    
    // ConsumerGroup已消费的Topic, partition的偏移量，last offset减去这个偏移量就是lag
    clusterAdmin, err := sarama.NewClusterAdminFromClient(client)
    if err != nil {
        panic("get clusterAdmin err")
    }
    
    ofr, err := clusterAdmin.ListConsumerGroupOffsets("dns_response_slbsdc", tp)
    if err != nil {
        panic("get ofr err")
    }
    b := ofr.GetBlock("PACKET_DNS_RESPONSE", 1)
    fmt.Printf("b: %+v\n", b)


    // 消费指定topic,消息个数的数据
    // msg, err := consumeMsg(topic, partition, offset, len)
    // if err != nil {
    //     panic("consume msg err")
    // }


    if err = broker.Close(); err != nil {
        panic(err)
    }
    

    return
}