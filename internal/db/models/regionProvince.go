package models

type Province struct {
	ProvinceID   int    `json:"province_id" gorm:"type:int; primary_key; AUTO_INCREMENT; not null; comment:省份id"`
	CountryCode  string `json:"country_code" gorm:"type:varchar(50); not null; comment:国家代码"`
	ProvinceCode string `json:"province_code" gorm:"type:varchar(50); not null; comment:省份代码"`
	ProvinceName string `json:"province_name" gorm:"type:varchar(150); not null; comment:省份名称"`
}

type ProvinceRes struct {
	ProvinceID   int    `json:"province_id"`
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

// 设置表名
func (Province) TableName() string {
	return "ip_region_province"
}

// 根据国家编码获取省份列表
func (model Province) GetProvinceListByCountry(countryCode string) (provinceList []ProvinceRes, err error) {
	err = DB.Model(&model).
		Select("province_id, province_code, province_name").
		Where("country_code =?", countryCode).
		Find(&provinceList).Error
	return
}

// 根据国家编码获取省份列表
func (model Province) GetProvinceList() (provinceList []Province, err error) {
	err = DB.Model(&model).
		Find(&provinceList).Error
	return
}
