package kafka

import (
	"ipfast_server/pkg/util/log"
	kk "ipfast_server/pkg/util/messageKafka"

	"github.com/spf13/viper"
)

var Reader *kk.KafkaReader

func init() {
	Reader = &kk.KafkaReader{}
}

var Topic = ""

/*
初始化kafka加载配置

	param:
		readKafka bool 是消费者还是生产者
	return:
		error 错误信息
*/
func Setup() (err error) {
	Topic = viper.GetString("kafka.topic")
	brokers := viper.GetStringSlice("kafka.brokers")

	Reader = &kk.KafkaReader{}
	err = Reader.Init(brokers, Topic, viper.GetInt("kafka.minBytes"), viper.GetInt("kafka.maxBytes"), viper.GetInt("kafka.partition"), viper.GetString("kafka.group"))
	if err != nil {
		return
	}
	log.Info("[Kafka Connect Url:%+v]:SUCCESS", brokers)
	return
}
