package messageKafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

/*
kafka配置结构体

	Brokers []string Kafka 服务器地址
	Topic string 主题(读写消息使用)

	GroupId string 消费者组ID 设置消费者组ID时,会自动均衡读取消息队列(平均分配分区数量 每个分区最多只能被一个消费者读取,分区数量少于消费者时,会导致消费者空闲)
	Partition int 分区号 有 消费者组ID时 Partition设置 无效
	MinBytes int 限制最少读取字节消息 10e3 10KB
	MaxBytes int 限制最多读取字节消息 10e6 10MB
*/
type Config struct {
	// common
	Brokers []string
	Topic   string

	// reader
	GroupId   string
	Partition int
	MinBytes  int
	MaxBytes  int
}

/*
写入策略

	RoundRobin 轮询分区 简单的轮询分区策略。每个消息都会被依次写入每个分区。
	LeastBytes 最少字节分区 根据每个分区的字节数量选择分区。这样可以保证每个分区的负载基本均衡。
	Hash 哈希分区策略。生产者会根据消息的键的哈希值选择分区。这样可以保证具有相同键的消息总是发送到同一个分区。
*/
type (
	RoundRobin = kafka.RoundRobin
	LeastBytes = kafka.LeastBytes
	Hash       = kafka.Hash
)

/*
Kafka 消息生产者结构体

	Brokers: []string Kafka 服务器地址(可以配置多个地址,自动发现集群中的其他节点)
	Topic: string 主题
	Balancer: kafka.Balancer 写入策略
	Instance: *kafka.Writer Kafka Writer 生产者实例对象
*/
type KafkaWriter struct {
	Instance *kafka.Writer
}

/*
Kafka Writer(消息生产者) 初始化实例
*/
func (w *KafkaWriter) Init(brokers []string, topic string, balancer kafka.Balancer) error {
	if len(brokers) <= 0 {
		return fmt.Errorf("kafkaReader.Brokers is empty")
	}
	if topic == "" {
		return fmt.Errorf("kafkaReader.Topic is empty")
	}
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("kafkaWriter.Init panic:", r)
		}
	}()
	w.Instance = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: balancer,
	})
	return nil
}

/*
写入消息

	param:
		key: string 消息的 key
		data: []byte 消息的数据
	return:
		error
*/
func (w *KafkaWriter) WriteMessage(key string, data []byte) error {
	return w.Instance.WriteMessages(context.Background(), kafka.Message{Key: []byte(key), Value: data})
}

/*
关闭 Kafka Writer生产者连接
*/
func (w *KafkaWriter) CloseWriter() {
	w.Instance.Close()
}

/*
Kafka 消息消费者结构体

	Url: string Kafka 服务器地址
	Topic: string 主题
	GroupId: string 消费者组ID 用于实现消息偏移量自动管理 和自动均衡读取消息队列(平均分配分区数量 每个分区最多只能被一个消费者读取,分区数量少于消费者时,会导致消费者空闲)
	Partition int 分区号
	Instance: *kafka.Reader Kafka Writer 消费者实例对象
	MinBytes: int 限制最少读取字节消息 10e3 10KB
	MaxBytes: int 限制最多读取字节消息 10e6 10MB
*/
type KafkaReader struct {
	Instance *kafka.Reader
}

/*
Kafka Reader(消息消费者者) 初始化实例 并订阅主题
*/
func (r *KafkaReader) Init(brokers []string, topic string, minBytes, maxBytes, partition int, groupId string) error {
	if len(brokers) <= 0 {
		return fmt.Errorf("kafkaReader.Brokers is empty")
	}
	if topic == "" {
		return fmt.Errorf("kafkaReader.Topic is empty")
	}
	if minBytes <= 0 {
		minBytes = 10e3
	}
	if maxBytes <= 0 {
		maxBytes = 10e6
	}
	config := kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: minBytes,
		MaxBytes: maxBytes,
	}
	if groupId == "" {
		if partition < 0 {
			partition = 0
		}
		config.Partition = partition
	} else {
		config.GroupID = groupId
	}
	r.Instance = kafka.NewReader(config)
	return nil
}

/*
从 Kafka 读取消息

	param:
		topic: string 主题
	return:
		kafka.Message kafka消息
		error
*/
func (r *KafkaReader) ReadMessage(topic string) (kafka.Message, error) {
	return r.Instance.ReadMessage(context.Background())
}

/*
关闭 Kafka Reader消费者连接
*/
func (r *KafkaReader) CloseReader() {
	r.Instance.Close()
}

/*
向 Kafka 手动提交偏移量

	param:
		offset: int64 已读取的消息偏移量
	return:
		error 错误信息
*/
func (r *KafkaReader) SetOffset(offset int64) error {
	return r.Instance.SetOffset(offset)
}

/*
创建一个Kafka主题

	param:
		topic: string 主题
		url: string Kafka 服务器地址
		numPartitions: int 分区数
		replicationFactor: int 副本数(取决于broker数量 集群部署时 不能超过broker数量,否则会报错,单机部署时可以设置为1)
*/
func CreateKafkaTopic(topic, url string, numPartitions, replicationFactor int) (err error) {
	conn, err := kafka.Dial("tcp", url)
	if err != nil {
		return
	}
	defer conn.Close()
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}
	err = conn.CreateTopics(topicConfig)
	return
}
