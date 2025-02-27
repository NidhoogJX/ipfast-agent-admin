package test

import (
	"ipfast_server/internal/handler/ginHandler"
	"log"
	"testing"
)

// 测试动态添加翻译内容
func TestToken(t *testing.T) {
	token, err := ginHandler.GenerateToken("123456")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(token)
}
