package models

import "gorm.io/gorm"

// Assets 资产表
type Assets struct {
	ID            int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:资产ID"`
	UserID        int64  `json:"user_id" gorm:"type:bigint;comment:用户ID;not null;index"`
	CommodityName string `json:"commodity_name" gorm:"type:varchar(255);comment:套餐(商品)名称;not null"`
	ExpireTime    int64  `json:"expire_time" gorm:"type:bigint;comment:到期时间戳;not null;default:0"`
	Type          int8   `json:"type" gorm:"type:tinyint;comment:类型;default:1;not null"` //ToDo : 动态/静态  后续可拓展数据中心
	Status        int8   `json:"status" gorm:"type:tinyint;comment:状态:未使用/已使用/使用完毕;default:0;not null"`
	TotalCount    int64  `json:"total_count" gorm:"type:bigint;comment:资产资产数量;default:0;not null"`
	UsedCount     int64  `json:"used_asset" gorm:"type:bigint;comment:资产已分配数量/已使用数量;default:0;not null"`
	Unit          int8   `json:"unit" gorm:"type:tinyint;comment:单位;default:1;not null"` //ToDo : 流量GB/IP个数  后续可拓展数据中心
	AreaId        int64  `gorm:"not null;type:bigint;default:0;comment:国家ID;"`
	CreatedTime   int64  `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"` //资产购买时间
	UpdatedTime   int64  `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"` //资产更新时间
}

// Assets 资产表
type AssetsRes struct {
	ID            int64   `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:资产ID"`
	CommodityName string  `json:"commodity_name" gorm:"type:varchar(255);comment:套餐(商品)名称;not null"`
	ExpireTime    int64   `json:"expire_time" gorm:"type:bigint;comment:到期时间戳;not null;default:0"`
	Type          int8    `json:"type" gorm:"type:tinyint;comment:类型;default:1;not null"` //ToDo : 动态/静态  后续可拓展数据中心
	TotalCount    float64 `json:"total_count" gorm:"type:bigint;comment:资产资产数量;default:0;not null"`
	UsedCount     float64 `json:"used_asset" gorm:"type:bigint;comment:资产已分配数量/已使用数量;default:0;not null"`
	Unit          int8    `json:"unit" gorm:"type:tinyint;comment:单位;default:1;not null"`          //ToDo : 流量GB/IP个数  后续可拓展数据中心
	CreatedTime   int64   `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"` //资产购买时间
}

func (Assets) TableName() string {
	return "ip_assets"
}

func (AssetsRes) TableName() string {
	return "ip_assets"
}

func (model Assets) GetAssetsByUserID() (assets []AssetsRes, err error) {
	err = DB.Where("user_id = ? AND status != 2", model.UserID).Find(&assets).Error
	return
}

func (model Assets) GetUsedAssetsByUserID() (assets []AssetsRes, err error) {
	err = DB.Where("user_id = ? AND status !=2", model.UserID).Order("created_time desc").Find(&assets).Error
	return
}

func (model Assets) InsertBatch(models []Assets, tx *gorm.DB) error {
	return tx.Create(&models).Error
}
