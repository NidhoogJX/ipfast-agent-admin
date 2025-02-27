package scheduler

import (
	"ipfast_server/internal/db/library"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/cronscheduler"
	"ipfast_server/pkg/util/log"
)

var scheduler *cronscheduler.Scheduler

func init() {
	scheduler = cronscheduler.NewScheduler()
}

// StartScheduler 启动定时任务调度器
func StartScheduler() {
	// UpdateDurationTypeListTask()
	// UpdateStaticCountryTask()
	// UpdateCommoditiesListTask()
	// UpdateCountryCommoditiesTask()
	// UpdateLocationInfoTask()
	ClearVerificationCodeIPListTask()
	scheduler.Start()
}

func SyncFlow() {
	UpdateUserFlowTask()
	scheduler.Start()
}

// 定时任务 缓存动态IP商品列表
func UpdateCommoditiesListTask() {
	services.UpdateCommoditiesList()
	scheduler.AddJob("*/10 * * * * *", func() {
		log.Trace("缓存动态IP商品列表,提供查询")
		services.UpdateCommoditiesList()
	})
}

// 定时任务 缓存静态IP商品列表
func UpdateCountryCommoditiesTask() {
	services.UpdateCountryCommodities()
	services.UpdateDataCenterCommodities()
	scheduler.AddJob("*/10 * * * * *", func() {
		log.Trace("缓存静态和数据中心IP商品列表,提供查询")
		services.UpdateCountryCommodities()
		services.UpdateDataCenterCommodities()
	})
}

// 定时任务 缓存静态IP国家列表
func UpdateStaticCountryTask() {
	services.UpdateStaticCountry()
	scheduler.AddJob("*/10 * * * * *", func() {
		log.Trace("缓存静态IP国家列表,提供查询")
		services.UpdateStaticCountry()
	})
}

// 定时任务 缓存时长类型列表
func UpdateDurationTypeListTask() {
	services.UpdateDurationTypeList()
	scheduler.AddJob("*/10 * * * * *", func() {
		log.Trace("缓存静态IP国家列表,提供查询")
		services.UpdateDurationTypeList()
	})
}

// 定时任务 清除验证码IP限制列表
func ClearVerificationCodeIPListTask() {
	scheduler.AddJob("0 0 0 * * *", func() {
		log.Trace("清除验证码IP限制列表")
		library.ClearVerificationCodeIPList()
	})
}

// 定时任务 缓存用户当前流量信息
func UpdateUserFlowTask() {
	services.GetUserCurrentFlowData()
	scheduler.AddJob("*/60 * * * * *", func() {
		log.Trace("缓存子用户流量信息")
		services.GetUserCurrentFlowData()
	})
}
