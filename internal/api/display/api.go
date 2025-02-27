package display

import (
	"ipfast_server/internal/config/i18n"
	"ipfast_server/internal/handler/network/server"
)

/*
获取商务经理列表
*/

func BusinessInfo(resp server.Response) {
	resp.Res["business_name"] = "Jason"
	resp.Res["business_email"] = "akkiki1008@gmail.com"
	resp.Res["business_phone"] = "+86 18955348589"
	resp.Res["business_telegram"] = "https://t.me/Psyyychooo"
	resp.Success("operate success")
}

type CountryIPInfoRes struct {
	CountryName string `json:"country_name"`
	IconUrl     string `json:"country_icon_url"`
	IPCount     int64  `json:"ip_count"`
}

/*
获取国家IP信息展示
*/
func CountryIPInfo(resp server.Response) {
	lang := resp.Context.GetHeader("Accept-Language")
	I18n := i18n.NewLocalizer(lang)
	resp.Res["country_list"] = []CountryIPInfoRes{
		{
			I18n.F("CN"),
			GetIcon("CN"),
			982939,
		},
		{
			I18n.F("US"),
			GetIcon("US"),
			829389,
		},
		{
			I18n.F("RU"),
			GetIcon("RU"),
			756378,
		},
		{
			I18n.F("ES"),
			GetIcon("ES"),
			654897,
		},
		{
			I18n.F("DE"),
			GetIcon("DE"),
			564893,
		},
		{
			I18n.F("FR"),
			GetIcon("FR"),
			467928,
		},
		{
			I18n.F("GB"),
			GetIcon("GB"),
			384789,
		},
		{
			I18n.F("IT"),
			GetIcon("IT"),
			283947,
		},
	}
	resp.Success("operate success")
}
