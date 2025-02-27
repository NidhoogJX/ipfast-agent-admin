package models

/*
流量静态住宅区域管理
*/
type TrafficRegion struct {
	Id   int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:区域ID"`
	Name string `json:"name" gorm:"type:varchar(50);not null;comment:区域名称"`
	Code string `json:"code" gorm:"type:varchar(50);not null;comment:区域代码"`
	// 是否热门区域
	IsHot int8 `json:"is_hot" gorm:"type:tinyint;not null;default:0;comment:是否热门区域"`
	// 是否为推荐区域
	IsRecommend int8 `json:"is_recommend" gorm:"type:tinyint;not null;default:0;comment:是否为推荐区域"`
	// 是否启用
	IsEnabled int8 `json:"is_enabled" gorm:"type:tinyint;not null;default:0;comment:是否启用"`
	// 权重
	Weight      int8  `json:"weight" gorm:"type:tinyint;not null;default:0;comment:权重"`
	CreatedTime int64 `json:"created_time" gorm:"type:bigint;not null;comment:创建时间"`
	UpdatedTime int64 `json:"updated_time" gorm:"type:bigint;not null;comment:更新时间"`
}

var TrafficRegionFiled = []string{
	"name",
	// "code",
	// "is_hot",
	// "is_recommend",
	"is_enabled",
	// "weight",
	"updated_time",
}

// TableName 表名
func (a TrafficRegion) TableName() string {
	return "ip_traffic_region"
}

// SelectByEnable 查询启用区域 根据权重排序
func (model TrafficRegion) SelectByEnable() (trafficRegions []TrafficRegion, err error) {
	err = DB.Where("is_enabled = ?", 1).Order("weight").Find(&trafficRegions).Error
	return
}

// 添加地区信息
func (model TrafficRegion) Create() (err error) {
	return DB.Create(&model).Error
}

// 修改地区信息
func (model TrafficRegion) Update() (err error) {
	return DB.Model(&model).Select(TrafficRegionFiled).Updates(model).Error
}

// 删除地区信息
func (model TrafficRegion) Delete() (err error) {
	return DB.Model(&model).Where("id = ?", model.Id).Delete(&model).Error
}
