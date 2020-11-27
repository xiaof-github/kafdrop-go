package demo

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// Sarama configuration options
var (
    version  = "2.1.1"
    group    = ""    
    assignor = "range"
    oldest   = true
    verbose  = false
)

var wg sync.WaitGroup
var messages map[string][][]byte

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

    // offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
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

func consumeTopic(consumer sarama.Consumer, topic string, block *sarama.OffsetResponseBlock, partition int32, len int) {
    
    partitionConsumer, err := consumer.ConsumePartition(topic, partition, block.Offset)
    if err != nil {
        panic(err)    
    }

    defer func() {
        if err := partitionConsumer.Close(); err != nil {
            log.Fatalln(err)
        }
    }()

    defer wg.Done()

    consumed := 0
    for {
        select {
        case msg := <-partitionConsumer.Messages():
            log.Printf("Consumed message partition %d\n offset %d\n key %s\n value %s\n", partition, msg.Offset, msg.Key, msg.Value)
            consumed++;
            messages[topic] = append(messages[topic], msg.Value)
            time.Sleep(time.Second)            
        default :
            log.Printf("consumed: %d", consumed)
            time.Sleep(10*time.Second)
        }
        if (consumed > len) {
            log.Printf("consumed: %d, want: %d", consumed, len)
            break
        }                
    }
    
}

func consumeMsg(broker *sarama.Broker, addr []string, topic string, partitions map[string]int, len int) ([]string, error) {
    consumer, err := sarama.NewConsumer(addr, nil)
    if err != nil {
        panic(err)
    }

    // defer func() {
    //     if err := consumer.Close(); err != nil {
    //         log.Fatalln(err)
    //     }
    // }()

    offsr := sarama.OffsetRequest{
        Version: 1,		
    }

    partitionsNum := partitions[topic]

    for i:=0;i<partitionsNum;i++{
        offsr.AddBlock(topic, int32(i), sarama.OffsetOldest, 999999999)
    }

    // offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
    res1, err1 := broker.GetAvailableOffsets(&offsr)
    if err1 != nil {
        panic("broker offset error")
    }

    var block *sarama.OffsetResponseBlock
    
    for i:=0;i<partitionsNum;i++ {
        block = res1.GetBlock(topic, int32(i))
        wg.Add(1)
        go consumeTopic(consumer, topic, block, int32(i), len)
    }  
    
    
    
    wg.Wait()
    log.Printf("Consumed: \n")
    return []string{"a", "b", "c"} , nil
}


func main() {
    log.Println("Starting kafdrop-go")

    messages = make(map[string][][]byte)

    addrs := []string{"10.155.200.108:9092"}
    kafkaHost := addrs[0]
    topic1 := "my_topic"

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

    client, err := sarama.NewClient(addrs, config)
    if err != nil {
        panic("client create error")
    }
    defer client.Close()
    // 获取broker信息
    brokers := client.Brokers()
    for _, bro := range brokers {
        fmt.Println(bro.Addr())
    }
    // 获取controller信息
    brok, ok := client.Controller()
    if ok != nil {
        panic("controller")
    }
    fmt.Printf("controller: %s\n", brok.Addr())

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
        fmt.Printf("v.Name: %+v, v.Partitions: %d\n", v.Name, partitions[v.Name])
    }

    // 计算topic可消费消息总量
    sum := getTopicMsgNum(broker, partitions, topic1)
    fmt.Printf("topic[%s], sum: %d\n", topic1, sum)

    // tp := map[string][]int32{
    //     topic1:[]int32{0,1},
    // }

    // ConsumerGroup已消费的Topic, partition的偏移量
    offr := sarama.OffsetFetchRequest{
        Version: 5,
        ConsumerGroup: "group4",
    }
    offr.AddPartition(topic1, 1)	

    res, err := broker.FetchOffset(&offr)
    if err != nil {
        panic("broker fetchoffset error")
    }
    fmt.Printf("res: %+v\n", res)
    v1 := res.GetBlock(topic1, 1)
    fmt.Printf("v1: %+v\n", v1)	

    // 当前Topic, partition可消费的最小偏移量
    offsr := sarama.OffsetRequest{
        Version: 1,		
    }

    offsr.AddBlock(topic1, 1, sarama.OffsetOldest, 999999999)
    offsr.AddBlock(topic1, 2, sarama.OffsetOldest, 999999999)
    offsr.AddBlock(topic1, 3, sarama.OffsetOldest, 999999999)
    // offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
    res1, err := broker.GetAvailableOffsets(&offsr)
    if err != nil {
        panic("broker offset error")
    }
    fmt.Printf("res: %+v\n", res1)
    o1 := res1.GetBlock(topic1, 1)
    fmt.Printf("o1: %+v\n", o1)
    o2 := res1.GetBlock(topic1, 2)
    fmt.Printf("o2: %+v\n", o2)
    o3 := res1.GetBlock(topic1, 3)
    fmt.Printf("o3: %+v\n", o3)

    // 当前Topic, partition可消费的最大偏移量
    offsr.AddBlock(topic1, 1, sarama.OffsetNewest, 999999999)
    offsr.AddBlock(topic1, 2, sarama.OffsetNewest, 999999999)
    offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
    // offsr.AddBlock(topic1, 3, sarama.OffsetNewest, 999999999)
    res1, err = broker.GetAvailableOffsets(&offsr)
    if err != nil {
        panic("broker offset error")
    }
    fmt.Printf("res: %+v\n", res1)
    n1 := res1.GetBlock(topic1, 1)
    fmt.Printf("n1: %+v\n", n1)
    n2 := res1.GetBlock(topic1, 2)
    fmt.Printf("n2: %+v\n", n2)
    n3 := res1.GetBlock(topic1, 3)
    fmt.Printf("n3: %+v\n", n3)

    // partition lag, lag = lastOffset - curOffset
    // p1Lag := n1.Offset - v1.Offset
    // fmt.Printf("Topic: %s, Partition: 1, lag: %d\n", topic1, p1Lag)




    // 获取主题的名称集合
    topics, err := client.Topics()
    if err != nil {
        panic("get topics err")
    }
    
    for _, tv1 := range topics {
        // fmt.Println(e)
        _ = tv1		
    }
    abc, e := client.GetOffset(topic1, 1, sarama.OffsetOldest)
    fmt.Printf("first offset: %+v,",abc)
    fmt.Println(e)
    abc, e = client.GetOffset(topic1, 1, sarama.OffsetNewest)
    fmt.Printf("last offset: %+v,",abc)
    fmt.Println(e)

    
    // ConsumerGroup已消费的Topic, partition的偏移量，last offset减去这个偏移量就是lag
    // clusterAdmin, err := sarama.NewClusterAdminFromClient(client)
    // if err != nil {
    //     panic("get clusterAdmin err")
    // }
    
    // ofr, err := clusterAdmin.ListConsumerGroupOffsets("dns_response_slbsdc", tp)
    // if err != nil {
    //     panic("get ofr err")
    // }
    // b := ofr.GetBlock(topic1, 1)
    // fmt.Printf("b: %+v\n", b)


    // 消费指定topic,1000个消息保存到本地，返回给前端
    count := 1
    
    msg, err := consumeMsg(broker, addrs, topic1, partitions, count)
    if err != nil {
        fmt.Printf("consume msg err, %+v", msg)
        panic("consume msg err")
    }


    if err = broker.Close(); err != nil {
        panic(err)
    }
    

    return
}