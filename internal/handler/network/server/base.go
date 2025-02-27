package server

import (
	"ipfast_server/internal/handler/ginHandler"
	"ipfast_server/pkg/util/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
Web 服务器实例
*/
var webServerEngine *http.Server
var serverStarted bool = false

/*
生成Jwt令牌
*/
func GenerateToken(userId string) (token string, err error) {
	token, err = ginHandler.GenerateToken(userId)
	if err != nil {
		// 处理错误，例如记录日志或返回错误响应
		log.Error("生成令牌时出错: %v", err)
	}
	return
}

/*
路由结构体

	param:
		Path: 路由路径 (接口请求路径)
		Handler: 路由处理函数 (接口实际处理函数)
		RequestType: 请求类型 GET/POST
		JwtEnabled: 是否启用JWT验证
		RecaptchaEnabled: 是否启用Recaptcha验证
*/
type Router struct {
	Path             string
	RequestType      string
	Handler          func(Response)
	JwtEnabled       bool
	RecaptchaEnabled bool
	PullEnabled      bool
	TranslateEnabled bool
}

/*
启动Web服务器

	param:
		port int 端口
		readTimeout int 读取超时
		writeTimeout int 写入超时
	return:
		err error 可能的错误
*/
func Run() (err error) {
	serverStarted = true
	return ginHandler.StartWebServer(webServerEngine)
}

/*
关闭 Web 服务器

	return:
		err error 可能的错误
*/
func Stop() error {
	if serverStarted {
		return ginHandler.StopWebServer(webServerEngine)
	}
	return nil
}

/*
初始化Gin引擎实例

	param:
		mode string 运行模式
		ginlog string Gin日志模式
		routers []Router 路由
	return:
		*gin.Engine Gin引擎实例
*/
func InitGinEngine(mode string, routers []Router, recordLog, recovery, allowCors bool, port, readTimeout, writeTimeout int) {
	webServerEngine = ginHandler.SetGinEngine(
		ginHandler.GinParams{
			Mode:        strings.ToLower(mode),
			RecordLog:   recordLog,
			Recovery:    recovery,
			AllowCors:   allowCors,
			Middlewares: nil,
			RouterFuncs: injection(routers),
			Port:        port,
			ReadTime:    readTimeout,
			WriteTime:   writeTimeout,
		},
	)
}

/*
注入路由修改调用参数

	param:
		handlerFunc func(resp *base.Response) 自定义路由处理函数
	return:
		func(c *gin.Context) Gin路由处理函数
*/
func injection(routers []Router) (routerFuncs []ginHandler.RouterFunc) {
	if len(routers) > 0 {
		for _, router := range routers {
			routerFuncs = append(
				routerFuncs,
				ginHandler.RouterFunc{
					Path: router.Path,
					Handler: func(c *gin.Context) {
						response := Response{Context: c, Res: make(map[string]interface{})}
						router.Handler(response)
					},
					RequestType:      router.RequestType,
					JwtEnabled:       router.JwtEnabled,
					RecaptchaEnabled: router.RecaptchaEnabled,
					PullEnabled:      router.PullEnabled,
					TranslateEnabled: router.TranslateEnabled,
				},
			)
		}
	}
	return
}
