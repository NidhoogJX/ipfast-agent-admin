package models

type Country struct {
	CountryID   int    `json:"country_id" gorm:"type:int; primary_key; AUTO_INCREMENT; not null; comment:地区id"`
	CountryCode string `json:"country_code" gorm:"type:varchar(50); not null; comment:国家代码"`
	CountryName string `json:"country_name" gorm:"type:varchar(150); not null; comment:国家名称"`
}

// TableName 设置表名
func (Country) TableName() string {
	return "ip_region_country"
}

type CountryRes struct {
	CountryID   int    `json:"country_id"`   // 国家id
	CountryCode string `json:"country_code"` // 国家代码
	CountryName string `json:"country_name"` // 国家名称
}

// 获取国家名称列表
func (model Country) GetCountryList() (countryList []CountryRes, err error) {
	err = DB.Model(&model).Select("country_id ,country_code, country_name").Find(&countryList).Error
	return
}

// 根据地区id获取国家列表
func (model Country) SelectCountryListByRegionId(regionId int64) (countries []Country, err error) {
	err = DB.Table("ip_region_country AS irc").
		Select(`
			DISTINCT irc.country_id,
			irc.country_name
		`).
		Joins("LEFT JOIN ip_traffic_country AS itc ON irc.country_name = itc.name").
		Where("itc.region_id = ?", regionId).
		Scan(&countries).Error
	return
}
