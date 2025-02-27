package ginHandler

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
路由函数结构体

	param:
		Path: 路由路径 (接口请求路径)
		Handler: 路由处理函数 (接口实际处理函数)
		RequestType: 请求类型 GET/POST
		JwtEnabled: 是否启用JWT验证
		RecaptchaEnabled: 是否启用Recaptcha验证
		PullEnabled: 是否启用拉取检查
		TranslateEnabled: 是否启用翻译
*/
type RouterFunc struct {
	Path             string
	RequestType      string
	Handler          func(c *gin.Context)
	JwtEnabled       bool
	RecaptchaEnabled bool
	PullEnabled      bool
	TranslateEnabled bool
}

/*
Gin Web 框架参数

	param:
		Mode: 运行模式
		RecordLog: 是否记录日志
		Recovery: 是否恢复
		IsCors: 是否允许跨域
		RouterFuncs: 路由函数
		Middlewares: 中间件
*/
type GinParams struct {
	Mode        string
	RecordLog   bool
	Recovery    bool
	AllowCors   bool
	RouterFuncs []RouterFunc
	Middlewares []gin.HandlerFunc
	Port        int
	ReadTime    int
	WriteTime   int
}

/*
返回成功响应 code 0 和 Response中的Res数据
*/
func Success(c *gin.Context, msg string) bool {
	c.JSON(200, gin.H{
		"code": 0,
		"data": nil,
		"msg":  msg,
	})
	c.Abort()
	return true
}

/*
返回成功响应 code 0 和 Response中的Res数据
*/
func CodeResponse(c *gin.Context, code int, msg string) bool {
	c.JSON(code, gin.H{
		"code": code,
		"data": nil,
		"msg":  msg,
	})
	c.Abort()
	return true
}

/*
返回失败响应 code 1 和 Response中的Res数据
*/
func FailedResponse(c *gin.Context, err string) bool {
	c.JSON(200, gin.H{
		"code": -1,
		"data": nil,
		"msg":  err,
	})
	c.Abort()
	return false
}

func FailedResponseCode(c *gin.Context, code int, err string) bool {
	c.JSON(200, gin.H{
		"code": code,
		"data": nil,
		"msg":  err,
	})
	c.Abort()
	return false
}

/*
自定义响应写入器
*/
type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (int, error) {
	return w.Body.Write(data)
}
func (w *CustomResponseWriter) WriteString(s string) (int, error) {
	return w.Body.WriteString(s)
}
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w *CustomResponseWriter) SetHeader(key, value string) {
	headers := w.Header()
	headers.Set(key, value)
}

/*
Gin Web 框架
*/
var engine *gin.Engine

/*
初始化 Gin Web 框架(路由和中间件和一些配置,需要在启动Web服务前调用)

	param:
		ginParams GinParams Gin Web 框架参数
*/
func SetGinEngine(ginParams GinParams) *http.Server {
	// 初始化设置
	gin.SetMode(ginParams.Mode) // 设置 Gin 的运行模式。有三个值：debug、release、test
	engine = gin.New()
	// engine.Static("/static", "./static") // 静态文件目录

	// 注册全局中间件
	if ginParams.RecordLog {
		engine.Use(gin.Logger()) //日志中间件 记录信息太少了 Todo 重写日志中间件
	}
	if ginParams.Recovery {
		engine.Use(gin.Recovery()) // 恢复中间件在出现 panic 错误时恢复 防止服务器崩溃。
	}
	if ginParams.AllowCors {
		engine.Use(CorsMiddleware) // 允许CORS跨域请求
	}
	if len(ginParams.Middlewares) > 0 {
		engine.Use(ginParams.Middlewares...) // 注册自定义中间件
	}

	// 注册路由 和 路由处理中间件
	for _, routerFunc := range ginParams.RouterFuncs {
		middlewares := []gin.HandlerFunc{}
		if !routerFunc.TranslateEnabled {
			middlewares = append(middlewares, TranslateMiddleware)
		}
		if routerFunc.RecaptchaEnabled {
			middlewares = append(middlewares, VerifyRecaptcha)
		}
		if routerFunc.JwtEnabled {
			middlewares = append(middlewares, ValidateJWT)
		}
		if routerFunc.PullEnabled {
			middlewares = append(middlewares, PullMiddleware)
		}
		switch routerFunc.RequestType {
		case "GET":
			engine.GET(routerFunc.Path, append(middlewares, routerFunc.Handler)...)
		case "POST":
			engine.POST(routerFunc.Path, append(middlewares, routerFunc.Handler)...)
		default:
			engine.Any(routerFunc.Path, append(middlewares, routerFunc.Handler)...)
		}
	}
	webServerEngine := &http.Server{
		Addr:           fmt.Sprintf(":%d", ginParams.Port),
		Handler:        engine,
		ReadTimeout:    time.Duration(ginParams.ReadTime) * time.Second,
		WriteTimeout:   time.Duration(ginParams.WriteTime) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return webServerEngine
}

/*
启动 Web 服务器

	param:
		port int 服务监听端口号
		readTime int 读取超时时间
		writeTime int 写入超时时间
	return:
		*http.Server Web 服务器实例
		error 错误信息
*/
func StartWebServer(webServerEngine *http.Server) (err error) {
	err = webServerEngine.ListenAndServe()
	return
}

/*
关闭 Web 服务器

	param:
		webServerEngine *http.Server Web 服务器实例
	return:
		error 错误信息
*/
func StopWebServer(webServerEngine *http.Server) error {
	webServerEngine.Close()
	return nil
}
