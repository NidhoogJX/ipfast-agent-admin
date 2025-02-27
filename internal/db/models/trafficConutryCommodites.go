package models

// TrafficCountryCommodites 国家商品 关联表
type TrafficCountryCommodites struct {
	CommoditesId int64 `json:"commodites_id" gorm:"type:bigint;not null;comment:国家商品ID;index"`
	CountryId    int64 `json:"country_id" gorm:"type:bigint;not null;comment:国家ID;index"`
}

// TableName 表名
func (a TrafficCountryCommodites) TableName() string {
	return "ip_traffic_country_commodites"
}

type TrafficCountryWithCommodities struct {
	RegionId      int64   `json:"region_id" `
	CountryId     int64   `json:"country_id"`
	CountryName   string  `json:"country_name" `
	CountryCode   string  `json:"country_code"`
	IsHot         int8    `json:"is_hot"`
	IsRecommend   int8    `json:"is_recommend" `
	Weight        int8    `json:"weight" `
	CommodityId   int64   `json:"commodity_id"`
	CommodityName string  `json:"commodity_name"`
	Price         float64 `json:"price"`
	Currency      string  `json:"currency"`
	OriginalPrice float64 `json:"original_price"`
	StockQuantity int64   `json:"stock_quantity"`
	UsedQuantity  int64   `json:"used_quantity" `
	Flag          string  `json:"flag"`
}

func (TrafficCountryWithCommodities) TableName() string {
	return "ip_traffic_country"
}

func (model TrafficCountryWithCommodities) SelectAllByType(ctype int8) (trafficCountryDatas []TrafficCountryWithCommodities, err error) {
	err = DB.Table("ip_traffic_country AS itc").
		Select(`
		itc.region_id AS RegionId,
		itc.id AS CountryId,
		itc.code AS CountryCode,
		itc.name AS CountryName,
		itc.is_hot AS IsHot,
		itc.is_recommend AS IsRecommend,
		itc.weight AS Weight,
		ic.id AS CommodityId,
		ic.name AS CommodityName,
		ic.original_price AS OriginalPrice,
		ic.currency AS Currency,
		ROUND(ic.price * ic.discount, 2) AS price,
		itc.stock_quantity AS StockQuantity,
		itc.used_quantity AS UsedQuantity
		`).
		Joins("LEFT JOIN ip_traffic_country_commodites AS itcc ON itc.id = itcc.country_id").
		Joins("LEFT JOIN ip_commodities AS ic ON itcc.commodites_id = ic.id").
		Where("itc.is_enabled = ?", 1). // 国家表启用
		Where("ic.enable = ?", 1).      // 商品表启用
		Where("ic.type = ?", ctype).    // 2: 商品类型 只查询静态住宅类型
		Order("itc.weight").
		Scan(&trafficCountryDatas).Error
	return
}

func (model TrafficCountryWithCommodities) SelectByCountryIdByType(ctype int8) (trafficCountryDatas TrafficCountryWithCommodities, err error) {
	err = DB.Table("ip_traffic_country AS itc").
		Select(`
		CONCAT(itc.name, '|', itc.city_name) AS CountryName,
		ic.name AS CommodityName,
		ROUND(ic.price * ic.discount, 2) AS price,
		ic.original_price AS OriginalPrice,
		ic.unit AS Unit,
		ic.total_count AS TotalCount,
		ic.duration_type_id AS DurationTypeId,
		ic.currency AS Currency,
		ic.description AS Description,
		itc.id AS CountryId,
		itc.region_id AS RegionId,
		itc.code AS CountryCode,
		itc.city_name AS CityName,
		itc.city_code AS CityCode,
		itc.is_hot AS IsHot,
		itc.is_recommend AS IsRecommend,
		itc.weight AS Weight,
		ic.id AS CommodityId,
		itc.stock_quantity AS StockQuantity
		`).
		Joins("LEFT JOIN ip_traffic_country_commodites AS itcc ON itc.id = itcc.country_id").
		Joins("LEFT JOIN ip_commodities AS ic ON itcc.commodites_id = ic.id").
		Where("itc.is_enabled = ?", 1).       // 国家表启用
		Where("ic.enable = ?", 1).            // 商品表启用
		Where("itc.id = ?", model.CountryId). //查询对应国家ID
		Where("ic.type = ?", ctype).          // 2: 商品类型
		Order("itc.weight").
		Find(&trafficCountryDatas).Error
	return
}

type RenewCommodities struct {
	CommodityName  string  `json:"commodity_name"`
	Price          float64 `json:"price"`
	DurationTypeId int64   `json:"duration_type_id"`
	Currency       string  `json:"currency"`
	CountryId      int64   `json:"country_id"`
}

func (model TrafficCountryCommodites) SelectRenewCommodities() (renewCommodities []RenewCommodities, err error) {
	err = DB.Table("ip_traffic_country_commodites AS itcc").
		Select(`itcc.country_id,ic.name AS CommodityName,ROUND(ic.price * ic.discount, 2) AS price,ic.duration_type_id AS DurationTypeId,ic.currency AS Currency`).
		Joins("LEFT JOIN ip_commodities AS ic ON ic.id = itcc.commodites_id").
		Where("itcc.country_id = ?", model.CountryId).
		Order("ic.weight").
		Scan(&renewCommodities).Error
	return
}

// 根据国家id获取商品id
func (model TrafficCountryCommodites) SelectCommoditesIdByCountryID() (commoditesId int64, err error) {
	err = DB.Model(model).
		Select("commodites_id").
		Where("country_id = ?", model.CountryId).
		First(&commoditesId).Error
	return commoditesId, err
}
