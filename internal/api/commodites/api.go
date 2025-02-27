package commodites

import (
	"ipfast_server/internal/api/display"
	"ipfast_server/internal/db/models"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/log"
)

/*
获取商品列表
*/
func CommoditiesList(resp server.Response) {
	resp.Res["commodities"] = services.GetCommoditiesList()
	resp.Success("operate success")
}

/*
获取静态住宅国家商品列表
*/
func GetCountryCommoditiesList(resp server.Response) {
	datas, err := services.GetCountryCommodities()
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	var data []interface{}
	for _, v := range datas {
		v.Flag = display.GetIcon(v.CountryCode)
		data = append(data, v)
	}
	resp.Res["country_commodities"] = data
	resp.Success("operate success")
}

/*
获取静态住宅国家列表
*/
func GetCountryList(resp server.Response) {
	datas, err := services.GetCountryCommodities()
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	type countrtinfo struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	var data []countrtinfo
	for _, v := range datas {
		data = append(data, countrtinfo{
			Id:   v.CountryId,
			Name: v.CountryName,
		})
	}
	resp.Res["country_commodities"] = data
	resp.Success("operate success")
}

/*
获取数据中心国家列表
*/
func GetDataCountryList(resp server.Response) {
	datas, err := services.GetDataCommodities()
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	type countrtinfo struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	var data []countrtinfo
	for _, v := range datas {
		data = append(data, countrtinfo{
			Id:   v.CountryId,
			Name: v.CountryName,
		})
	}
	resp.Res["country_commodities"] = data
	resp.Success("operate success")
}

/*
获取数据中心国家商品列表
*/
func GetCountryDataCenterCommoditiesList(resp server.Response) {
	datas, err := services.GetDataCommodities()
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	var data []interface{}
	for _, v := range datas {
		v.Flag = display.GetIcon(v.CountryCode)
		data = append(data, v)
	}
	resp.Res["data_commodities"] = data
	resp.Success("operate success")
}

/*
获取静态住宅国家和数据中心区域分类
*/
func GetTrafficRegionList(resp server.Response) {
	data, err := services.GetTrafficRegionList()
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	resp.Res["regiony_list"] = data
	resp.Success("operate success")
}

/* 获取商品时长类型列表 */
// func DurationTypesList(resp server.Response) {
// 	data, err := services.GetStaticDurationTypeList()
// 	if err != nil {
// 		resp.Failed("operate failed")
// 		return
// 	}
// 	resp.Res["duration_types"] = data
// 	resp.Success("operate success")
// }

// 根据类型id获取时长类型列表
func DurationTypesList(resp server.Response) {
	param := struct {
		TypeId int8 `json:"type_id"`
	}{}
	err := resp.Bind(&param)
	if err != nil || (param.TypeId != 1 && param.TypeId != 2 && param.TypeId != 3) {
		log.Error("EmailRegister json error:%v", err)
		resp.Failed("param error")
		return
	}
	data, err := services.GetDurationTypeList(param.TypeId)
	if err != nil {
		resp.Failed("operate failed")
		return
	}
	resp.Res["duration_types"] = data
	resp.Success("operate success")
}

/* 获取续费套餐时长类型 */
// func GetRenewCommodities(resp server.Response) {
// 	param := struct {
// 		AddressStr string `json:"addressstr"` //续费国家地址
// 	}{}
// 	err := resp.Bind(&param)
// 	if err != nil {
// 		log.Error("EmailRegister json error:%v", err)
// 		resp.Failed("param error")
// 		return
// 	}
// 	addressParts := strings.Split(param.AddressStr, "-")
// 	if len(addressParts) != 2 {
// 		resp.Failed("address format error")
// 		return
// 	}
// 	countryCode := addressParts[0]
// 	cityCode := addressParts[1]
// 	data, err := services.GetTrafficRegionInfo(countryCode, cityCode)
// 	if err != nil {
// 		log.Error("GetRenewCommodities error:%v", err)
// 		resp.Failed("get traffic region info failed")
// 		return
// 	}
// 	resp.Res["renew_commodities"] = data
// 	resp.Success("operate success")
// }

// 根据国家id，时长类型获取续费价格
func GetRenewCommodities(resp server.Response) {
	param := struct {
		CountryId      int64 `json:"country_id"`
		DurationTypeId int64 `json:"duration_type_id"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		log.Error("%v", err)
		resp.Failed("param error")
		return
	}
	// 获取续费商品信息
	commodities, err := services.GetCommodity(param.CountryId)
	if err != nil {
		log.Error("failed to get commodites info:%v", err)
		resp.Failed("failed to get commodites info")
		return
	}
	// 根据时长类型获取价格倍率
	durationType, err := services.GetDurationType(param.DurationTypeId)
	if err != nil {
		log.Error("failed to get duration type:%v", err)
		resp.Failed("failed to get duration type")
		return
	}
	renewCommodities := models.RenewCommodities{
		CommodityName:  commodities.Name,
		Price:          commodities.Price * durationType.MultiplyingPower,
		DurationTypeId: commodities.DurationTypeId,
		Currency:       commodities.Currency,
		CountryId:      param.CountryId,
	}
	resp.Res["renew_commodities"] = renewCommodities
	resp.Success("operate success")
}
