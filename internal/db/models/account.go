package models

import (
	"gorm.io/gorm"
)

/*
流量账号管理
*/
type Account struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:账号ID"`
	Name        string `json:"name" gorm:"type:varchar(150);default:'';not null;comment:账号名称"`
	StopTime    int64  `json:"stop_time" gorm:"type:bigint unsigned;default:0;not null;comment:账号到期时间"`
	TotalFlow   int64  `json:"total_flow" gorm:"type:bigint unsigned;default:0;not null;comment:总流量"`
	UsedTraffic int64  `json:"used_traffic" gorm:"type:bigint unsigned;default:0;not null;comment:已使用流量"`
	Status      int8   `json:"status" gorm:"type:tinyint unsigned;default:0;not null;comment:账号状态"`
	UserLevel   int8   `json:"user_level" gorm:"type:tinyint unsigned;default:0;not null;comment:用户会员等级"`
	CreateTime  int64  `json:"create_time" gorm:"type:bigint unsigned;default:0;not null;comment:注册时间"`
	UpdateTime  int64  `json:"update_time" gorm:"type:bigint unsigned;default:0;not null;comment:会员信息上次更新时间"`
}

var Field = []string{
	"name",
	"stop_time",
	"total_flow",
	"used_traffic",
	"status",
	"user_level",
	"create_time",
	"update_time",
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (Account) TableName() string {
	return "ip_account"
}

func (model Account) SelectByAccountIdAndName() ([]Account, error) {
	var account []Account
	err := DB.Select("name,id").Where("used_traffic < total_flow and status = ?", model.Status).Find(&account).Error
	return account, err
}

func (model Account) Create() error {
	return DB.Select(Field).Create(model).Error
}

func (model Account) Update() error {
	return DB.Select(Field).Updates(model).Error
}

func (model Account) Find() (accountData Account, err error) {
	err = DB.Where("id = ?", model.Id).Find(&accountData).Error
	return
}

func (model Account) UpdateFlows(updates []Account) error {
	// 开始事务
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 遍历更新 map
	for _, account := range updates {
		if err := tx.Model(&account).
			UpdateColumn("used_traffic", gorm.Expr("used_traffic + ?", account.UsedTraffic)).Error; err != nil {
			// 如果有错误，回滚事务
			tx.Rollback()
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error
}
