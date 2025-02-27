package models

// 套餐时长类型表 秒 分 时 日  月  年
// const (
// 	DurationTypeSecond = 1
// 	DurationTypeMinute = 2
// 	DurationTypeHour   = 3
// 	DurationTypeDay    = 4
// 	DurationTypeMonth  = 5
// 	DurationTypeYear   = 6
// )

// DurationTypes 时长类型表
type DurationTypes struct {
	ID               int64   `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:时长类型ID"`
	Name             string  `json:"name" gorm:"type:varchar(50);not null;comment:时长类型名称"`
	Type             int8    `json:"type" gorm:"type:tinyint;not null;comment:时长类型(1-秒;2-分;3-时;4-日;5-月;6-年)"`
	CType            int8    `json:"commodities_type" gorm:"type:tinyint;not null;comment:商品类型"`
	Count            int64   `json:"count" gorm:"type:bigint;not null;comment:时长数量"`
	MultiplyingPower float64 `json:"multiplying_power" gorm:"type:decimal(10,2);not null;comment:价格倍率"`
	Status           int8    `json:"status" gorm:"type:tinyint;not null;default:0;comment:状态"`
	Weight           int8    `json:"weight" gorm:"type:tinyint;not null;default:0;comment:权重"`
	CreateTime       int64   `json:"create_time" gorm:"type:bigint;not null;default:0;comment:创建时间"`
	UpdateTime       int64   `json:"update_time" gorm:"type:bigint;not null;default:0;comment:更新时间"`
}

// TableName 表名
func (a DurationTypes) TableName() string {
	return "ip_duration_types"
}

// SelectByEnable 查询指定状态的时长类型 根据权重排序
func (model DurationTypes) SelectByStatus() (durationTypes []DurationTypes, err error) {
	err = DB.Where("status = ?", model.Status).Order("weight").Find(&durationTypes).Error
	return
}

// SelectByEnable 查询时长类型 根据权重排序
func (model DurationTypes) SelectAll() (durationTypes []DurationTypes, err error) {
	err = DB.Where("status = ?", 1).Find(&durationTypes).Error
	return
}

// 根据ID查询时长类型
func (model DurationTypes) SelectByID() (durationType DurationTypes, err error) {
	err = DB.Where("id = ?", model.ID).First(&durationType).Error
	return
}
