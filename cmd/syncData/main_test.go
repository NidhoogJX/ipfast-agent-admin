package main

import (
	_ "ipfast_server/internal/config/syncData"
	"ipfast_server/internal/scheduler"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/log"
	"testing"
)

func TestMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln("Application crashed with error: %v", r)
		}
	}()
	services.StartReceiveFlowData()
	scheduler.SyncFlow()
	select {}
}
