package log

import (
	"io"
	"ipfast_server/pkg/util/logBase"
	syslog "log"

	"github.com/spf13/viper"
)

/*
日志配置结构体

	prarm:
		GlobalLevel uint32 全局日志级别
		PrintToConsole bool 是否在控制台打印日志
		LogFilePath string 日志文件路径
*/
type Config struct {
	GlobalLevel    string
	PrintToConsole bool
	LogFilePath    string
}

/*
全局日志配置
*/
var LogConfig *Config

/*
全局日志实例
*/
var logger logBase.Log

/*
初始化全局日志实例

	prarm:
		GlobalLevel uint32 全局日志级别 默认为LevelDebug
		PrintToConsole bool 是否在控制台打印日志 默认为true
		LogFilePath string 日志文件路径 默认为default.log
		LevelInfoMap map[uint32]LogColorInfo 日志级别描述和颜色 可以修改此配置为自定义颜色
		LevelMap map[string]uint32 日志级别映射 可以修改此配置为自定义日志级别
		Logger *log.Logger 日志记录器 默认为syslog.New(io.Discard, "", 0)
*/
func init() {

	logger = logBase.Log{
		GlobalLevel:    logBase.LevelDebug,
		PrintToConsole: true,
		LogFilePath:    "default.log",
		LevelInfoMap:   logBase.LevelInfoMap,
		LevelMap:       logBase.LevelMap,
		Logger:         syslog.New(io.Discard, "", 0),
	}
}

/*
全局日志实例加载配置
*/
func Setup() (err error) {
	level := viper.GetString("logging.level")
	console := viper.GetBool("logging.console")
	path := viper.GetString("logging.path")
	if level == "" {
		level = "INFO"
		console = true
	}
	if path == "" {
		path = "default.log"
		console = true
	}
	logger.LogFilePath = path
	logger.PrintToConsole = console
	err = logger.SetLevel(level)
	return
}

/*
记录Fatalln等级日志并直接退出程序

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Fatalln(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelFatalln, format, a...)
}

/*
记录Error等级日志

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Error(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelError, format, a...)
}

/*
记录Debug等级日志

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Debug(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelDebug, format, a...)
}

/*
记录Info等级日志

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Info(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelInfo, format, a...)
}

/*
记录Warning等级日志

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Warning(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelWarning, format, a...)
}

/*
记录Trace等级日志

	param:
		format string 日志格式
		a ...interface{} 日志参数
*/
func Trace(format string, a ...interface{}) {
	logger.LogWithColor(logBase.LevelTrace, format, a...)
}
