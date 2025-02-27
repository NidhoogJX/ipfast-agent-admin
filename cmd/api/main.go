package main

import (
	_ "ipfast_server/internal/config/api"
	"ipfast_server/internal/scheduler"
	"ipfast_server/pkg/util/log"
	// "net/http"
	// _ "net/http/pprof"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Application crashed with error: %v", r)
		}
	}()
	scheduler.StartScheduler() // 启动定时任务 调度器
	select {}
}
