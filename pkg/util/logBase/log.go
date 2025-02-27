package logBase

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

/*
日志级别

	LevelFatalln: 0 uint32 致命错误 影响程序运行 直接退出
	LevelError: 1 uint32 错误
	LevelWarning: 2 uint32 警告
	LevelInfo: 3 uint32 信息
	LevelDebug: 4 uint32 调试
	LevelTrace: 5 uint32 追踪
	LevelSilent: 6  uint32 不打印日志
*/
const (
	LevelFatalln uint32 = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
	LevelSilent
)

/*
ANSI 颜色转义码

	ResetColor: string 重置颜色
	Red: string 红色
	Green: string 绿色
	Yellow: string 黄色
	Blue: string 蓝色
	Purple: string 紫色
	Cyan: string 青色
	White: string 白色
*/
const (
	ResetColor = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
)

/*
日志级别名称映射到其对应的数字值
*/
var LevelMap = map[string]uint32{
	"fatalln": LevelFatalln,
	"error":   LevelError,
	"warning": LevelWarning,
	"info":    LevelInfo,
	"debug":   LevelDebug,
	"trace":   LevelTrace,
	"silent":  LevelSilent,
}

/*
将日志级别映射到其对应的日志描述和颜色
*/
var LevelInfoMap = map[uint32]LogColorInfo{
	LevelFatalln: {"致命错误", Red},
	LevelError:   {"错误", Red},
	LevelWarning: {"警告", Yellow},
	LevelInfo:    {"信息", Green},
	LevelDebug:   {"调试", Purple},
	LevelTrace:   {"追踪", Blue},
}

/*
存储日志描述和颜色结构体

	Text: string 日志描述
	Color: string 日志颜色
*/
type LogColorInfo struct {
	Text  string
	Color string
}

type Log struct {
	GlobalLevel    uint32                  //全局日志级别
	PrintToConsole bool                    //是否在控制台打印日志
	LogFilePath    string                  //日志文件路径
	LevelInfoMap   map[uint32]LogColorInfo //日志级别描述和颜色
	LevelMap       map[string]uint32       //日志级别映射
	Logger         *log.Logger             //日志记录器
}

// 修改日志级别 默认为DEBUG
func (loginstance *Log) SetLevel(level interface{}) (err error) {
	switch v := level.(type) {
	case uint32:
		if v > LevelSilent || v < LevelFatalln {
			err = fmt.Errorf(Red+"[Error] Unsupported LogLevel Value Level:%v(%T)"+ResetColor, level, v)
		} else {
			loginstance.GlobalLevel = v
		}
	case string:
		newGlobalLevel, ok := loginstance.LevelMap[strings.ToLower(v)]
		if !ok {
			err = fmt.Errorf(Red+"[Error] Invalid LogLevel Value :%s"+ResetColor, v)
		} else {
			loginstance.GlobalLevel = newGlobalLevel
		}
	default:
		// 报告错误
		err = fmt.Errorf(Red+"[Error] Unsupported LogLevel Type :%v(%T)"+ResetColor, level, v)
	}
	return
}

/*
根据日志级别记录对应级别颜色日志 到控制台和文件
如果全局日志级别大于或等于 LevelNone，或全局日志级别大于调用等级level,函数将直接返回，不会打印任何内容。

	param:
		level uint32 日志级别
		format string 日志格式
*/
func (loginstance *Log) LogWithColor(level uint32, format string, a ...interface{}) {
	if (loginstance.GlobalLevel >= LevelSilent || loginstance.GlobalLevel < level) && loginstance.GlobalLevel != LevelFatalln {
		return
	}
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	levelName := loginstance.LevelInfoMap[level].Text
	colorCode := loginstance.LevelInfoMap[level].Color
	tmpStr := fmt.Sprintf(time.Now().Format("[01-02 15:04:05]") + "[" + levelName + "]" + format)
	// 打开日志文件
	logFile, err := os.OpenFile(loginstance.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer logFile.Close()

	// 设置日志记录器的输出
	if loginstance.PrintToConsole {
		loginstance.Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	} else {
		loginstance.Logger.SetOutput(logFile)
	}

	// 打印日志
	loginstance.Logger.Println(colorCode + tmpStr + ResetColor)
	if level == LevelFatalln {
		os.Exit(1)
	}
}
