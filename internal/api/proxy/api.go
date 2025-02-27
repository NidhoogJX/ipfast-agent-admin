package proxy

import (
	"fmt"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
)

// 获取DFJ代理地址列表
func GetDFJProxyUrlList(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	data, err := services.GetEnableProxyServerList()
	if err != nil {
		resp.Failed("get proxy server list failed")
		return
	}
	var urlList []string
	for _, v := range data {
		urlList = append(urlList, fmt.Sprintf("%s", v.Address))
	}
	resp.Res["url_list"] = urlList
	resp.Success("operate success")
}

// API提取服务地址
func GetAPIProxyServiceURL(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	data, err := services.GetEnableProxyServerList()
	if err != nil {
		resp.Failed("get proxy server list failed")
		return
	}
	var urlList []string
	for _, v := range data {
		urlList = append(urlList, fmt.Sprintf("%s:%d", v.Address, v.Port))
	}
	resp.Res["url_list"] = urlList
	resp.Success("operate success")
}

// 账密提取节点
func GetAccountExtractionNodes(resp server.Response) {
	param := struct {
		Num           int    `json:"num" `           // 数量
		Country       int    `json:"country"`        // 国家
		Province      int    `json:"province" `      // 省份
		City          int    `json:"city" `          // 城市
		IPTime        int    `json:"ip_time" `       // IP时间
		SubUserID     int64  `json:"sub_user_id"`    // 子用户ID
		IPValidity    *int8  `json:"ip_validity"`    // IP有效性
		ProxyProtocol string `json:"proxy_protocol"` // 代理协议
		ProxyURL      string `json:"proxy_url"`      // 代理服务地址
		Format        string `json:"format"`         // 格式
	}{}
	if err := resp.Bind(&param); err != nil {
		resp.Failed(err.Error())
		return
	}
	if param.Country < 0 || param.Province < 0 || param.City < 0 || param.Num < 0 || param.SubUserID < 0 || param.IPTime < 0 {
		resp.Failed("param error")
		return
	}
	if param.ProxyURL == "" || param.Num > 100 || param.Format == "" {
		resp.Failed("param error")
		return
	}
	if param.ProxyProtocol == "" || (param.ProxyProtocol != "http" && param.ProxyProtocol != "socks5") {
		resp.Failed("param error")
		return
	}
	if param.IPValidity == nil {
		ipValidity := int8(1)
		param.IPValidity = &ipValidity
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	data, err := services.GetAccountExtractionNodes(
		services.AccountParams{
			Country:         param.Country,
			Province:        param.Province,
			City:            param.City,
			ProxyURL:        param.ProxyURL,
			ProxyProtocol:   param.ProxyProtocol,
			SessionValidity: *param.IPValidity,
			SessionTime:     param.IPTime,
			Format:          param.Format,
		}, param.Num, uid, param.SubUserID)
	if err != nil {
		resp.Failed("get account extraction nodes failed")
		return
	}
	resp.Res["proxy_url_list"] = data
	resp.Success("operate success")
}
