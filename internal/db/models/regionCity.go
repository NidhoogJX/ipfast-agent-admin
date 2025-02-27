package models

type City struct {
	CityID       int    `json:"city_id" gorm:"type:int; primary_key; AUTO_INCREMENT; not null; comment:城市id"`
	CountryCode  string `json:"country_code" gorm:"type:varchar(50); not null; comment:国家代码"`
	ProvinceCode string `json:"province_code" gorm:"type:varchar(50); not null; comment:省份代码"`
	CityName     string `json:"city_name" gorm:"type:varchar(150); not null; comment:城市名称"`
}

type CityRes struct {
	CityID   int    `json:"city_id"`
	CityName string `json:"city_name"`
}

// 设置表名
func (City) TableName() string {
	return "ip_region_city"
}

// 根据省份编码获取城市列表
func (model City) GetCityListByProvince(countryCode, provinceCode string) (cityList []CityRes, err error) {
	err = DB.Model(&model).
		Select("city_id,city_name").
		Where("country_code = ? and province_code = ?", countryCode, provinceCode).
		Find(&cityList).Error
	return
}

// 根据省份编码获取城市列表
func (model City) GetCityList() (cityList []City, err error) {
	err = DB.Model(&model).
		Find(&cityList).Error
	return
}
