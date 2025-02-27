package ipipgo

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"ipfast_server/internal/handler/network/request"
	"ipfast_server/pkg/util/log"
)

// ===================== 正式 接口 =====================
// // API 网关地址
// const ApiGateway = "https://manager.ap0.cn/proxy/api"
// const ApiKey = "13256669137"
// const ApiSign = "d6276b19abaa82e83c2d9cd0f2fb6cf7"

// // API 接口路径
// const (
//
//	AccountPath             = "/v2/account"             // 创建二级客户
//	StaticPath              = "/generate/static"        // 开套餐(静态IP)
//	StaticDetailPath        = "/generate/static/detail" // 开套餐(静态IP) 并返回详细信息
//	StaticIPRenewPath       = "/static/renew"           // 静态IP续费
//	StaticIPRenewDetailPath = "/static/renew/detail"    // 静态IP续费 并返回详细信息
//	GetStaticIPListPath     = "/ip/list"                // 获取静态IP列表
//
// )
// ===================== Mock 接口 =====================
// API 网关地址
const ApiGateway = "https://mock.apipost.net/mock/2ab6bf5e0c64000/proxy/api"
const ApiKey = "13256669137"
const ApiSign = "d6276b19abaa82e83c2d9cd0f2fb6cf7"

// API 接口路径
const (
	AccountPath             = "/v2/account?apipost_id=1b96975bf13040"             // 创建二级客户
	StaticPath              = "/generate/static?apipost_id=1b976043713064"        // 开套餐(静态IP)
	StaticDetailPath        = "/generate/static/detail?apipost_id=1b980652b13068" // 开套餐(静态IP) 并返回详细信息
	StaticIPRenewPath       = "/static/renew"                                     // 静态IP续费
	StaticIPRenewDetailPath = "/static/renew/detail?apipost_id=1b980652b13068"    // 静态IP续费 并返回详细信息
	GetStaticIPListPath     = "/ip/list?apipost_id=1b951c59b13030"                // 获取静态IP列表
)

// CommonResponse 公共响应字段
type CommonResponse struct {
	Code    int    `json:"code"`    //响应码 0:成功
	Message string `json:"message"` //响应消息
}

// 二级账号认证信息响应
type AuthInfo struct {
	Acc string `json:"acc"` //认账套餐账号
	Pwd string `json:"pwd"` //认证账号密码
}

// 二级客户账号数据
type AccountData struct {
	Account  string              `json:"account"`  //二级客户账户名
	Sign     string              `json:"sign"`     //二级客户账户加密字符串
	Key      *string             `json:"key"`      // 提取密钥
	AuthInfo map[string]AuthInfo `json:"authInfo"` //认证套餐账密信息
}

// 静态IP套餐账号数据
type StatciIPData struct {
	Cmild  int `json:"cmiId"` //客户套餐ID
	IpList []struct {
		IpPort          string `json:"ipPort"`          //IP地址端口
		AccountPassword string `json:"accountPassword"` //账号|密码
		AddressStr      string `json:"addressStr"`      //国家地区名称
		Status          int    `json:"status"`          //状态 0失效1过期2正常
		EndTime         string `json:"endTime"`         //到期时间
		MealId          int    `json:"mealId"`          //套餐ID
		CmiId           int    `json:"cmiId"`           //客户套餐ID
	} `json:"ipList"` //开通的IP
}

type CountryData struct {
	Country string `json:"country"`
	IpNum   int64  `json:"ipNum"`
	Region  string `json:"region"`
}
type CommonParam struct {
	Key      string `json:"key"`
	Sign     string `json:"sign"`
	NonceStr string `json:"nonceStr"`
}

// StaticIPParam 开套餐(静态IP)请求参数
type StaticIPParam struct {
	CommonParam
	MealId          int           `json:"mealId"`       //套餐ID
	CustomerName    string        `json:"customerName"` // ipipgo账号名
	CountryDataList []CountryData `json:"countryDataList"`
}

// StaticIPParam 续费套餐(静态IP)请求参数
type StaticIPRenewParam struct {
	CommonParam
	CmiId        int64  `json:"cmiId"`        // 客户套餐ID
	MealTime     int    `json:"mealId"`       // 套餐ID
	CustomerName string `json:"customerName"` // ipipgo账号名
}

// 二级账号接口响应
type AccountResponse struct {
	CommonResponse
	Data AccountData `json:"data"` //响应数据
}

// StaticResponse 开套餐(静态IP)响应
type StaticResponse struct {
	CommonResponse
	Data StatciIPData `json:"data"` //响应数据
}

// StaticResponse 续费套餐(静态IP)响应
type StaticRenewResponse struct {
	CommonResponse
	Data []struct {
		IpPort          string `json:"ipPort"`          //IP地址端口
		AccountPassword string `json:"accountPassword"` //账号|密码
		AddressStr      string `json:"addressStr"`      //国家地区名称
		Status          int    `json:"status"`          //状态 0失效1过期2正常
		EndTime         string `json:"endTime"`         //到期时间
		MealId          int    `json:"mealId"`          //套餐ID
	} `json:"data"` //开通的IP `json:"data"` //响应数据
}

func getNonceStr() string {
	hash := md5.New()
	hash.Write([]byte("test"))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateAccount 创建二级客户
func CreateAccount() (*AccountResponse, error) {

	data := &struct {
		Key      string `json:"key"`
		Sign     string `json:"sign"`
		NonceStr string `json:"nonceStr"`
	}{
		Key:      ApiKey,
		Sign:     ApiSign,
		NonceStr: getNonceStr(),
	}
	body, _, err := request.Post(ApiGateway+AccountPath, data)
	if err != nil {
		return nil, err
	}
	log.Info("CreateAccount response: %s", string(body))
	// 解析响应
	resp := &AccountResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf(" api:[%s],code: %d, message: %s", ApiGateway+AccountPath, resp.Code, resp.Message)
	}
	return resp, nil
}

// StaticIP 开套餐(静态IP)
func StaticIP(data *StaticIPParam) error {
	body, _, err := request.Post(ApiGateway+StaticPath, data)
	if err != nil {
		return err
	}
	log.Info("StaticIP response: %s", string(body))
	resp := &CommonResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf(" api:[%s],code: %d, message: %s", ApiGateway+StaticPath, resp.Code, resp.Message)
	}
	return nil
}

// StaticIPDetail 开套餐(静态IP) 并返回详细信息
func StaticIPDetail(data *StaticIPParam) (*StaticResponse, error) {
	data.Key = ApiKey
	data.Sign = ApiSign
	data.NonceStr = getNonceStr()
	jsonData, _ := json.Marshal(data)
	log.Info("StaticIPDetail request: %s", jsonData)
	body, _, err := request.Post(ApiGateway+StaticDetailPath, data)
	if err != nil {
		return nil, err
	}
	log.Info("StaticIPDetail response: %s", string(body))
	resp := &StaticResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf(" api:[%s],code: %d, message: %s", ApiGateway+StaticDetailPath, resp.Code, resp.Message)
	}
	return resp, nil
}

// StaticIPRenew 静态IP续费
func StaticIPRenew(data *StaticIPParam) (*CommonResponse, error) {
	body, _, err := request.Post(ApiGateway+StaticIPRenewPath, data)
	if err != nil {
		return nil, err
	}
	log.Info("StaticIPRenew response: %s", string(body))
	resp := &CommonResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf(" api:[%s],code: %d, message: %s", ApiGateway+StaticIPRenewPath, resp.Code, resp.Message)
	}
	return resp, nil
}

// StaticIPRenewDetail 静态IP续费 并返回详细信息
func StaticIPRenewDetail(data *StaticIPRenewParam) (*StaticRenewResponse, error) {
	body, _, err := request.Post(ApiGateway+StaticIPRenewDetailPath, data)
	if err != nil {
		return nil, err
	}
	log.Info("StaticIPRenewDetail response: %s", string(body))
	resp := &StaticRenewResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("api:[%s],code: %d, message: %s", ApiGateway+StaticIPRenewDetailPath, resp.Code, resp.Message)
	}
	return resp, nil
}

type StaticIPListParam struct {
	CommonParam
	Type  int   `json:"type"`  //类型 1:静态IP
	CmiId int64 `json:"cmiId"` //客户套餐ID
}

func GetStaticIPList(data StaticIPListParam) (*StaticResponse, error) {
	body, _, err := request.Post(ApiGateway+GetStaticIPListPath, data)
	if err != nil {
		return nil, err
	}
	log.Info("GetStaticIPList response: %s", string(body))
	resp := &StaticResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("api:[%s],code: %d, message: %s", ApiGateway+StaticIPRenewDetailPath, resp.Code, resp.Message)
	}
	return resp, nil
}
