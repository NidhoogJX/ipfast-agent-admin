package main

import (
	_ "ipfast_server/internal/config/api"
	"ipfast_server/internal/scheduler"
	"ipfast_server/pkg/util/log"
	"testing"
	// "net/http"
	// _ "net/http/pprof"
)

func TestMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Application crashed with error: %v", r)
		}
	}()
	// go func() {
	// 	log.Info("启动pprof服务 0.0.0.0:6060")
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()
	log.Info("启动定时任务 调度器")
	scheduler.StartScheduler() // 启动定时任务 调度器
	select {}
}
