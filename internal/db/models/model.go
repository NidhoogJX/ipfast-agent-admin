package models

import "ipfast_server/internal/db/core/gorm"

var DB = gorm.MasterDb

/*
注册模型
这里的所有模型都会在程序启动的时候自动迁移
*/
func AutoMigrateAllModels() error {
	return gorm.MasterDb.AutoMigrate(
		// &Assets{},
		// &Account{},
		// &Announcement{},
		// &Commodities{},
		// &DurationTypes{},
		// &TrafficRegion{},
		// &TrafficCountry{},
		// &TrafficCountryCommodites{},
		// &TransactionOrders{},
		// &OrderItem{},
		&User{},
		// &VerificationCode{},
		// &IPIPGoAccount{},
		// &IPIPGoAccountStaticIp{},
		// &IpWhiteList{},
		&UserFlow{},
		&SubUser{},
		// &Country{},
		// &Province{},
		// &City{},
		// &ProxyServer{},
		// &PayPlatform{},
		&FlowRecord{},
		&IpRecord{},
		&FlowRecord{},
		&Admin{},
		&Agent{},
		&Recharge{},
	) // 自动迁移数据库
}
