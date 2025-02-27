package models

import (
	"log"
	"math"
)

// Commodities 商品表
type Commodities struct {
	ID             int64   `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:商品ID"`
	Name           string  `json:"name" gorm:"type:varchar(255);comment:名称;not null"`
	Price          float64 `json:"price" gorm:"type:decimal(10,2);comment:价格;default:0"`
	OriginalPrice  float64 `json:"original_price" gorm:"type:decimal(10,2);comment:原价;default:0"`
	Enable         int8    `json:"enable" gorm:"type:tinyint;comment:是否启用;default:0;not null"`
	Type           int8    `json:"type" gorm:"type:tinyint;comment:商品类型;default:0;not null"`
	Unit           int8    `json:"unit" gorm:"type:tinyint;comment:单位(1-GB;2-TB;3-PB);default:1;not null"`
	Discount       float64 `json:"discount" gorm:"type:decimal(10,2);comment:折扣;default:1"`
	TotalCount     int64   `json:"total_count" gorm:"type:bigint;comment:对应资产数量(流量);default:0;not null"`
	DurationTypeId int64   `json:"duration_type_id" gorm:"type:bigint;comment:时长类型ID;default:0;not null"`
	Currency       string  `json:"currency" gorm:"type:varchar(255);comment:货币;not null"`
	Description    string  `json:"description" gorm:"type:varchar(500);comment:描述"`
	Weight         int64   `json:"weight" gorm:"type:bigint;comment:权重;default:0;not null"`
	CreatedTime    int64   `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdatedTime    int64   `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
}

// Commodities 商品表 关联查询时长类型表
type CommoditiesDuration struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name" `
	Price          float64 `json:"price" `
	OriginalPrice  float64 `json:"original_price" `
	Enable         int8    `json:"enable" `
	Type           int8    `json:"type"`
	Unit           int8    `json:"unit" `
	Discount       float64 `json:"discount" gorm:"type:decimal(10,2);comment:折扣;default:0"`
	TotalCount     int64   `json:"total_count" `
	DurationTypeId int64   `json:"duration_type_id" `
	DurationName   string  `json:"duration_name" `
	DurationType   int8    `json:"duration_type" `
	DurationCount  int64   `json:"duration_count" `
	Currency       string  `json:"currency" `
	Description    string  `json:"description" `
	Weight         int64   `json:"weight"`
}

var CommoditiesFieid = []string{
	"name",
	"label",
	"price",
	"term",
	"total_count",
	"weight",
	"enable",
	"duration_type_id",
	"description",
}

// TableName 设置表名
func (Commodities) TableName() string {
	return "ip_commodities"
}

// SelectByEnable 查询启用商品和时长类型连表
func (model Commodities) SelectByDurationAndEnable() (commoditiesDatas []CommoditiesDuration, err error) {
	err = DB.Table("ip_commodities AS ic").
		Select(`
		ic.id,
		ic.name,
		ROUND(ic.price * ic.discount, 2) AS price,
		ic.original_price,
		ic.unit,
		ic.type,
		ic.total_count,
		ic.currency,
		ic.description,
		ic.weight,
		ic.duration_type_id,
		idt.type as duration_type,
		idt.name as duration_name,
		idt.count as duration_count
		`).
		Joins("LEFT JOIN ip_duration_types AS idt ON ic.duration_type_id = idt.id").
		Where("ic.enable = ?", 1).  // 商品表启用
		Where("idt.c_type = ?", 1). // 1: 商品类型 只查询动态IP类型
		Where("ic.type = ?", 1).    // 1: 商品类型 只查询动态IP类型
		Order("ic.weight").
		Scan(&commoditiesDatas).Error
	return
}

// SelectByEnable 查询启用商品和时长类型连表
func (model Commodities) SelectByCountryId(countryId int64) (commoditiesDatas []CommoditiesDuration, err error) {
	err = DB.Table("ip_commodities AS ic").
		Select(`
		ic.id,
		ic.name,
		ROUND(ic.price * ic.discount, 2) AS price,
		ic.original_price,
		ic.unit,
		ic.type,
		ic.total_count,
		ic.currency,
		ic.description,
		ic.weight,
		ic.duration_type_id,
		idt.type as duration_type,
		idt.name as duration_name,
		idt.count as duration_count
		`).
		Joins("LEFT JOIN ip_duration_types AS idt ON ic.duration_type_id = idt.id").
		Where("ic.enable = ?", 1). // 商品表启用
		Where("idt.type = ?", 1).  // 1: 商品类型 只查询动态IP类型
		Where("ic.type = ?", 1).   // 1: 商品类型 只查询动态IP类型
		Order("ic.weight").
		Scan(&commoditiesDatas).Error
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// SelectByEnable 查询启用商品和时长类型连表
func (model Commodities) SelectByDurationAndEnableSinger() (commoditiesData CommoditiesDuration, err error) {
	err = DB.Table("ip_commodities AS ic").
		Select(`
		ic.id,
		ic.name,
		ROUND(ic.price * ic.discount, 2) AS price,
		ic.original_price,
		ic.unit,
		ic.type,
		ic.total_count,
		idt.count as duration_count,
		ic.currency,
		ic.description,
		ic.weight,
		idt.name as duration_name
		`).
		Joins("LEFT JOIN ip_duration_types AS idt ON ic.duration_type_id = idt.id").
		Where("ic.enable = ?", 1). // 商品表启用
		Where("ic.type = ?", 1).   // 1: 商品类型 只查询动态IP类型
		Where("ic.id = ?", model.ID).
		Order("ic.weight").
		First(&commoditiesData).Error
	return
}

// SelectByEnable 查询启用商品
func (model Commodities) SelectByEnable() (commoditiesDatas []Commodities, err error) {
	err = DB.Where("enable = 1").Find(&commoditiesDatas).Error
	for i, v := range commoditiesDatas {
		commoditiesDatas[i].Price = math.Round(v.Price * v.Discount)
	}
	return
}

// 根据商品ID查询商品
func (model Commodities) SelectByID() (commoditiesData Commodities, err error) {
	err = DB.Where("id = ?", model.ID).First(&commoditiesData).Error
	return
}
