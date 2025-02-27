package server

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
Response 响应结构体

	Context *gin.Context 请求上下文
	Res ResponseData 响应数据
	I18n i18n.MyLocalizer 国际化

	func Success() 用于处理成功响应，生成统一的响应格式
	func Failed() 用于处理失败响应，生成统一的响应格式
	func Response(httpCode int, code uint8) 用于处理响应，生成统一的响应格式
	func Json(dataStrcut interface{}) error 用于解析请求中的json数据
*/
type Response struct {
	Context *gin.Context
	Res     ResponseData
}
type ResponseData map[string]interface{}
type RequestData ResponseData

/*
写入响应时间

	param:
		time string 响应时间
*/
func (resp Response) WriteResponseTime(d time.Duration) {
	time := "0秒"
	us := d.Microseconds()
	if us < 1000 {
		time = fmt.Sprintf("%dus", us)
	} else {
		ms := d.Milliseconds()
		if ms < 1000 {
			time = fmt.Sprintf("%dms", ms)
		} else {
			s := d.Seconds()
			time = fmt.Sprintf("%.2fs", s)
		}
	}
	resp.Context.Writer.Header().Set("X-Elapsed-Time", time)
}
func (resp Response) Get(key, targetType string) string {
	valueStr, exits := resp.Context.Get(key)
	if !exits {
		return ""
	}
	value, err := dynamicTypeAssert(valueStr, targetType)
	if err != nil {
		return ""
	}
	return value.(string)
}
func (resp Response) GetUserID(key string) int64 {
	value, exists := resp.Context.Get(key)
	if !exists {
		return 0
	}

	valueStr, ok := value.(string)
	if !ok {
		return 0
	}

	valueInt, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0
	}

	return valueInt
}

// 动态类型转换函数
func dynamicTypeAssert(value interface{}, targetType string) (interface{}, error) {
	v := reflect.ValueOf(value)
	switch targetType {
	case "string":
		if v.Kind() == reflect.String {
			return v.String(), nil
		}
	case "int":
		if v.Kind() == reflect.Int {
			return int(v.Int()), nil
		}
	case "float64":
		if v.Kind() == reflect.Float64 {
			return v.Float(), nil
		}
	// 添加更多类型转换的情况
	default:
		return nil, fmt.Errorf("不支持的目标类型: %s", targetType)
	}
	return nil, fmt.Errorf("类型转换失败: %v 到 %s", reflect.TypeOf(value), targetType)
}

/*
返回成功响应 code 0 和 Response中的Res数据
*/
func (resp Response) Success(msg string) bool {
	resp.Context.JSON(200, gin.H{
		"code": 0,
		"data": resp.Res,
		"msg":  msg,
	})
	return true
}

/*
返回成功响应 code 0 和 Response中的Res数据
*/
func (resp Response) Code(code int, msg string) bool {
	resp.Context.JSON(code, gin.H{
		"code": 0,
		"data": resp.Res,
		"msg":  msg,
	})
	resp.Context.Abort()
	return true
}

/*
返回失败响应 code 1 和 Response中的Res数据
*/
func (resp Response) Failed(err string) bool {
	resp.Context.JSON(200, gin.H{
		"code": 1,
		"data": resp.Res,
		"msg":  err,
	})
	return false
}

/*
设置返回响应 httpCode 和 code 和 Response中的Res数据
*/
func (resp Response) Response(httpCode int, code uint8, msg string) {
	resp.Context.JSON(httpCode, gin.H{
		"code": code,
		"data": resp.Res,
		"msg":  msg,
	})
}

/*
Json 用于解析请求中的json数据

	param:
		dataStrcut interface{} 用于接收解析后的数据
*/
func (resp Response) Json(dataStrcut interface{}) error {
	return resp.Context.ShouldBindJSON(dataStrcut)
}

/*
Json 用于解析请求中的json数据

	param:
		dataStrcut interface{} 用于接收解析后的数据
*/
func (resp Response) Bind(dataStrcut interface{}) error {
	return resp.Context.Bind(dataStrcut)
}

/*
RawData 用于获取请求中的原始数据

	return:
		[]byte 请求中的原始数据
		error 错误信息
*/
func (resp Response) RawData() ([]byte, error) {
	return resp.Context.GetRawData()
}
