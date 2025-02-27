package models

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

/*
流量静态住宅国家区域管理
*/
type TrafficCountry struct {
	Id            int64                  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:国家ID"`
	RegionId      int64                  `json:"region_id" gorm:"type:bigint;not null;comment:区域ID"`
	Name          string                 `json:"name" gorm:"type:varchar(50);not null;comment:国家名称"`
	Code          string                 `json:"code" gorm:"type:varchar(50);not null;comment:国家代码"`
	CityName      string                 `json:"city_name" gorm:"type:varchar(100);comment:城市名称"`
	CityCode      string                 `json:"city_code" gorm:"type:varchar(100);comment:城市代码"`
	PrvanceName   string                 `json:"prvance_name" gorm:"type:varchar(100);comment:省份名称"`
	PrvanceCode   string                 `json:"prvance_code" gorm:"type:varchar(100);comment:省份代码"`
	IsHot         int8                   `json:"is_hot" gorm:"type:tinyint;not null;default:0;comment:是否热门区域"`
	IsRecommend   int8                   `json:"is_recommend" gorm:"type:tinyint;not null;default:0;comment:是否为推荐区域"`
	IsEnabled     int8                   `json:"is_enabled" gorm:"type:tinyint;not null;default:0;comment:是否启用"`
	Weight        int8                   `json:"weight" gorm:"type:tinyint;not null;default:0;comment:权重"`
	StockQuantity int64                  `json:"stock_quantity" gorm:"type:bigint;not null;default:0;comment:IP总库存数量"`
	UsedQuantity  int64                  `json:"used_quantity" gorm:"type:bigint;not null;default:0;comment:已使用IP库存数量"`
	CreateTime    int64                  `json:"create_time" gorm:"type:bigint;not null;comment:创建时间"`
	UpdateTime    int64                  `json:"update_time" gorm:"type:bigint;not null;comment:更新时间"`
	Version       optimisticlock.Version // 引入乐观锁版本号
}

// 更新库存数量
func (model *TrafficCountry) UpdateStockQuantity() error {
	return DB.Instance.Transaction(func(tx *gorm.DB) error {
		// 使用事务更新库存数量
		if err := tx.Table(model.TableName()).Where("id = ? AND stock_quantity >= used_quantity + ?", model.Id, model.UsedQuantity).Update("used_quantity", gorm.Expr("used_quantity + ?", model.UsedQuantity)).Error; err != nil {
			// 如果更新失败，返回错误，事务将回滚
			return err
		}
		// 如果更新成功，返回 nil，事务将提交
		return nil
	})
}

// SelectByEnable 查询启用的国家
func (model TrafficCountry) SelectByEnable() (data []TrafficCountry, err error) {
	err = DB.Table(model.TableName()).Where("is_enabled = 1").Find(&data).Error
	return
}

// SelectByEnable 查询区域代码指定国家信息
func (model TrafficCountry) SelectByCountryCodeAndCityCode() (data *TrafficCountry, err error) {
	data = &TrafficCountry{}
	err = DB.Table(model.TableName()).Where("is_enabled = 1 AND code = ? AND city_code = ?", model.Code, model.CityCode).Find(data).Error
	return
}

// TableName 表名
func (a TrafficCountry) TableName() string {
	return "ip_traffic_country"
}

// 根据ID获取信息
func (model TrafficCountry) GetById(id int64) (trafficCountry TrafficCountry, err error) {
	err = DB.Table(model.TableName()).Where("id = ?", id).First(&trafficCountry).Error
	return
}

// 更新已用库存数量
func (model TrafficCountry) UpdateUsedQuantity(trafficCountry TrafficCountry) error {
	return DB.Table(model.TableName()).Where("id = ?", trafficCountry.Id).Update("used_quantity", trafficCountry.UsedQuantity).Error
}
