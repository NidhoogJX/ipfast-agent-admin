package services

import (
	"errors"
	"fmt"
	"ipfast_server/internal/config/countrycode"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"strings"
	"sync"
)

type Commodities = models.Commodities
type CommoditiesDuration = models.CommoditiesDuration
type TrafficCountryWithCommodities = models.TrafficCountryWithCommodities

// 商品列表响应结构（动态IP）
type CommoditiesRes struct {
	ID            int64    `json:"id" `
	Name          string   `json:"name" `
	Price         float64  `json:"price" `
	OriginalPrice float64  `json:"original_price"`
	TotalCount    int64    `json:"total_count" `
	Unit          string   `json:"unit" `
	Duration      string   `json:"duration" `
	Currency      string   `json:"currency" `
	Description   []string `json:"description" `
	Weight        int64    `json:"weight" `
}

// 静态数据中心 国家库存
var (
	staticCountryList = sync.Map{}
	mutex             sync.Mutex
)

// 商品类型
const (
	DisEnableType     = iota     // 禁用
	DynamicIPType                // 动态IP类型标识
	StaticIPType                 // 静态IP类型标识
	DataIPType                   // 数据中心IP类型标识
	RenewStaticIPType = int8(10) // 续费静态IP类型标识
	RenewDataIPType   = int8(11) // 续费数据中心IP类型标识
)

// 商品单位
var UnitMap = map[int8]map[int8]string{
	DynamicIPType: {
		1: "GB",
		2: "TB",
		3: "PB",
	},
	StaticIPType: {
		1: "个",
	},
	// DataType:{
	// 	1: "个",
	// },
}

// 动态商品列表缓存
var dynamicCommoditiesList = sync.Map{}

// 静态 和数据中心商品列表缓存
var commoditiesList = sync.Map{}

// 定时任务：更新动态IP商品列表缓存
func UpdateCommoditiesList() {
	commoditiesListDatas, err := Commodities{}.SelectByDurationAndEnable()
	dynamicCommoditiesList = sync.Map{}
	if err != nil {
		log.Error("updateCommoditiesList error:%v", err)
		return
	}
	//log.Info("commoditiesListDatas:%v", commoditiesListDatas)
	for _, v := range commoditiesListDatas {
		dynamicCommoditiesList.Store(v.ID, v)
	}
}

// UpdateCountryCommodities 定时任务：更新静态住宅国家商品列表缓存
func UpdateCountryCommodities() {
	data, err := models.TrafficCountryWithCommodities{}.SelectAllByType(StaticIPType)
	if err != nil {
		log.Error("updateCountryCommodities error:%v", err)
		return
	}
	commoditiesList.Store("static_commodities", data)
}

// UpdateCountryCommodities 定时任务：更新数据中心国家商品列表缓存
func UpdateDataCenterCommodities() {
	data, err := models.TrafficCountryWithCommodities{}.SelectAllByType(DataIPType)
	if err != nil {
		log.Error("updateCountryCommodities error:%v", err)
		return
	}
	commoditiesList.Store("data_commodities", data)
}

// GetCommoditiesList 获取全部动态IP商品列表（缓存）
func GetCommoditiesList() (commoditiesRes []CommoditiesRes) {
	dynamicCommoditiesList.Range(func(key, value interface{}) bool {
		v, ok := value.(CommoditiesDuration)
		if !ok {
			return true
		}

		durationName := fmt.Sprintf("%d%s", v.DurationCount, GetDurationTypeMap()[v.DurationType])
		commoditiesRes = append(commoditiesRes, CommoditiesRes{
			ID:            v.ID,
			Name:          v.Name,
			Price:         v.Price,
			OriginalPrice: v.OriginalPrice,
			TotalCount:    v.TotalCount,
			Unit:          UnitMap[v.Type][v.Unit],
			Duration:      durationName,
			Currency:      v.Currency,
			Description:   strings.Split(v.Description, "\\n"),
			Weight:        v.Weight,
		})
		return true
	})
	return commoditiesRes
}

// GetCountryCommodities 获取静态住宅国家商品列表
func GetCountryCommodities() ([]TrafficCountryWithCommodities, error) {
	value, ok := commoditiesList.Load("static_commodities")
	if !ok {
		return nil, errors.New("static_commodities not found")
	}
	commodities, ok := value.([]TrafficCountryWithCommodities)
	if !ok {
		return nil, fmt.Errorf("类型转换失败: %v", value)
	}
	return commodities, nil
}

// GetCountryCommodities 获取数据中心国家商品列表
func GetDataCommodities() ([]TrafficCountryWithCommodities, error) {
	value, ok := commoditiesList.Load("data_commodities")
	if !ok {
		return nil, errors.New("data_commodities not found")
	}
	commodities, ok := value.([]TrafficCountryWithCommodities)
	if !ok {
		return nil, fmt.Errorf("类型转换失败: %v", value)
	}
	return commodities, nil
}

// GetCommodityByID 根据商品ID获取一条动态商品数据（缓存）
func GetCommodityByID(id int64) (*CommoditiesDuration, error) {
	commodity, ok := dynamicCommoditiesList.Load(id)
	if !ok {
		return nil, fmt.Errorf("commodity with id %d not found", id)
	}
	data, ok := commodity.(CommoditiesDuration)
	if !ok {
		return nil, errors.New("类型转换失败")
	}
	return &data, nil
}

// GetCountryCommodities 根据国家ID获取静态商品详情
func GetCountryCommoditieByCid(countryId int64, commoditie_type int8) (commoditie TrafficCountryWithCommodities, err error) {
	commoditie, err = models.TrafficCountryWithCommodities{
		CountryId: countryId,
	}.SelectByCountryIdByType(commoditie_type)
	if commoditie.CountryId == 0 || err != nil {
		return commoditie, fmt.Errorf("commodity with id %d not found", countryId)
	}
	return
}

// UpdateStaticCountry 定时任务：更新静态住宅国家库存缓存
func UpdateStaticCountry() {
	data, err := models.TrafficCountry{}.SelectByEnable()
	if err != nil {
		log.Error("updateCountryCommodities error:%v", err)
		return
	}
	staticCountryList = sync.Map{}
	for _, v := range data {
		staticCountryList.Store(v.Id, v)
	}
}

// CheckAndUpdateStaticCountryStock检查并更新静态住宅国家的库存
func CheckAndUpdateStaticCountryStock(countryQuantities map[int64]int64) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	// 首先检查所有国家的库存是否足够
	for countryId, quantity := range countryQuantities {
		value, ok := staticCountryList.Load(countryId)
		if !ok {
			return fmt.Errorf("static country not found for countryId: %d", countryId)
		}
		countries, ok := value.(models.TrafficCountry)
		if !ok {
			return fmt.Errorf("类型转换失败: %v", value)
		}
		// UsedQuantity：IP库存数量    StockQuantity：IP总库存数量
		if countries.UsedQuantity+quantity > countries.StockQuantity {
			return fmt.Errorf("insufficient inventory %s", countries.Name)
		}
	}

	// 如果所有国家的库存都足够，则更新库存
	for countryId, quantity := range countryQuantities {
		value, _ := staticCountryList.Load(countryId)
		countries := value.(models.TrafficCountry)
		countries.UsedQuantity += quantity
		// 更新库存
		staticCountryList.Store(countryId, countries)
		model := &models.TrafficCountry{
			Id:           countryId,
			UsedQuantity: quantity,
		}
		err := model.UpdateStockQuantity()
		if err != nil {
			log.Error("库存更新失败: %v", err)
			return fmt.Errorf("库存不足")
		}
	}

	return nil
}

func GetStaticCountry(countryId int64) (data *models.TrafficCountry, err error) {
	value, _ := staticCountryList.Load(countryId)
	countries := value.(models.TrafficCountry)
	if countries.Id == 0 {
		return nil, fmt.Errorf("static country not found for countryId: %d", countryId)
	}
	data = &countries
	return
}

type TrafficRegionRes struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetTrafficRegionList 获取静态住宅国家区域分类
func GetTrafficRegionList() ([]TrafficRegionRes, error) {
	var datas []TrafficRegionRes
	data, err := models.TrafficRegion{}.SelectByEnable()
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		datas = append(datas, TrafficRegionRes{
			ID:   v.Id,
			Name: v.Name,
		})
	}
	return datas, nil
}

func GetTrafficRegionInfo(code, cityCode string) (renewCommodities []models.RenewCommodities, err error) {
	countryCode := countrycode.GetCountryCode3(code)
	if countryCode == "" {
		return nil, fmt.Errorf("country code not found for code: %s", code)
	}
	data, err := models.TrafficCountry{
		Code:     countryCode,
		CityCode: cityCode,
	}.SelectByCountryCodeAndCityCode()
	if err != nil {
		log.Error("GetTrafficRegionInfo error:%v", err)
		return nil, fmt.Errorf("static country not found for code: %s", code)
	}
	if data.Id == 0 {
		return nil, fmt.Errorf("static country not found for code: %s", code)
	}
	renewCommodities, err = models.TrafficCountryCommodites{
		CountryId: data.Id,
	}.SelectRenewCommodities()
	if len(renewCommodities) == 0 {
		return nil, fmt.Errorf("renew commodities not found for code: %s", code)
	}
	return
}

// 根据国家id获取商品信息
func GetCommodity(countryId int64) (commodities Commodities, err error) {
	commoditesId, err := models.TrafficCountryCommodites{
		CountryId: countryId,
	}.SelectCommoditesIdByCountryID()
	if err != nil {
		log.Error("failed to get commodites id %v", err.Error())
		return
	}
	// 获取商品详细信息
	commodities, err = models.Commodities{
		ID: commoditesId,
	}.SelectByID()
	if err != nil {
		log.Error("failed to get commodites %v", err.Error())
		return
	}
	return commodities, nil
}
