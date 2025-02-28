package test

import (
	"ipfast_server/internal/handler/ginHandler"
	"log"
	"testing"
	"time"
)

// 测试动态添加翻译内容
func TestToken(t *testing.T) {
	token, err := ginHandler.GenerateToken("123456")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(token)
	now := time.Now()
	now.Year()
	log.Println("毫秒时间戳：", now.UnixMilli())
	msec := now.UnixMilli() % 1000
	log.Println("当前毫秒：", msec)
}
