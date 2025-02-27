package syncData

import (
	"ipfast_server/internal/config"

	"ipfast_server/internal/db/core/gorm"
	"ipfast_server/internal/db/core/kafka"
	"ipfast_server/pkg/util/log"
)

/*
初始化配置
*/
func init() {
	config.LoadConfig("config", "yaml", ".")
	loggingInit()
	databaseInit()
	kafkaMessageQueueInit() // 初始化kafka消息队列
	config.WatchingConfig()
}

/*
初始化日志配置
*/
func loggingInit() {
	err := log.Setup()
	if err != nil {
		log.Fatalln("程序启动失败:%s", err.Error())
	}
	config.SetWatching("logging", func(oldWebConfig, newWebConfig interface{}) {
		log.Info("日志配置发生变化,将应用最新配置:\n旧配置:%+v,\n新配置:%+v", oldWebConfig, newWebConfig)
		err := log.Setup()
		if err != nil {
			log.Fatalln("更新日志配置失败:%s", err.Error())
		}
	}, nil)
}

/*
初始化数据库
*/
func databaseInit() {
	gorm.Setup()
	config.SetWatching("database", func(oldDatabaseConfig, newDatabaseConfig interface{}) {
		log.Info("数据库配置发生变化,将重新连接数据库:\n旧配置:%+v,\n新配置:%+v", oldDatabaseConfig, newDatabaseConfig)
		gorm.Setup()
	}, nil)

}

/*
初始化kafka消息队列
*/
func kafkaMessageQueueInit() {
	kafka.Setup()
	// go func() {
	// // 初始化同步子账户生产者
	// kafka.SyncSubUserSetup()
	// log.Info("同步子账户消息队列初始化成功")

	// // 初始化同步流量生产者
	// kafka.SyncTrafficProducerSetup()
	// log.Info("同步流量消息队列初始化成功")

	// 初始化同步流量消费者
	// }()
}
