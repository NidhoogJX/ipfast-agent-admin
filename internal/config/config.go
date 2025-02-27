package config

import (
	"ipfast_server/pkg/util/log"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var serverDebug bool

func init() {
	serverDebug = false
}

/*
多实例
*/
func NewViper() *viper.Viper {
	return viper.New()
}

/*
监听配置文件变化数据结构
*/
type WatchConfigData struct {
	key            string
	handler        func(interface{}, interface{})
	oldConfigValue interface{}
}

/*
监听配置文件变化数据
*/
var wactching []WatchConfigData

/*
添加监听
*/
func SetWatching(key string, handler func(interface{}, interface{}), oldConfigValue interface{}) {
	wactching = append(wactching, WatchConfigData{
		key:            key,
		handler:        handler,
		oldConfigValue: oldConfigValue,
	})
}

/*
加载配置文件

	param:
		configName: 配置文件名称
		suffix: 配置文件后缀
		path: 配置文件路径
	return:
		error: 错误信息
*/
func LoadConfig(configName, suffix, path string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType(suffix)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	serverDebug = viper.GetBool("server.debug")
	return nil
}

/*
检查当前是否为debug模式
*/
func Debug() bool {
	return serverDebug
}

/*
监听配置文件变化
*/
func WatchingConfig() {
	viper.WatchConfig()
	for index, watch := range wactching {
		wactching[index].oldConfigValue = viper.Get(watch.key)
	}
	var debounceTimer *time.Timer
	viper.OnConfigChange(func(e fsnotify.Event) {
		serverDebug = viper.GetBool("server.debug")
		if serverDebug {
			log.Info("[Debug模式开启]会打印大量日志可能对服务的性能有一定损耗,线上谨慎开启!!!")
		}
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(3*time.Second, func() {
			for index, watch := range wactching {
				newConfigValue := viper.Get(watch.key)
				if !reflect.DeepEqual(wactching[index].oldConfigValue, newConfigValue) {
					watch.handler(wactching[index].oldConfigValue, newConfigValue)
					wactching[index].oldConfigValue = newConfigValue
				}
			}
		})
	})
}
