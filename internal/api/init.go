package api

import (
	"ipfast_server/internal/api/announcements"
	"ipfast_server/internal/api/recharge"
	"ipfast_server/internal/api/subuser"
	"ipfast_server/internal/api/user"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/pkg/util/log"

	"github.com/spf13/viper"
)

const (
	userPath                = "/user"          // 用户路径
	commoditiesPath         = "/commodities"   // 商品路径
	displayPath             = "/display"       // 展示数据路径
	proxyPath               = "/proxy"         // 代理路径
	announcementsRouterPath = "/announcements" // 系统公告路径
	payPath                 = "/pay"           // 支付路径
	rechargePath            = "/recharge"      // 充值路径
	subuserPath             = "/subuser"       // 子用户路径
)

// 用户路由
var userRouter = []server.Router{
	{
		RequestType: "POST",
		Path:        userPath + "/login",
		Handler:     user.Login,
		// RecaptchaEnabled: true,
	},
	{
		RequestType: "POST",
		Path:        userPath + "/logout",
		Handler:     user.LoginOut,
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/apply/verifcode",
		Handler:     user.ApplyVerificationCode,
		RequestType: "POST",
	},
	{
		Path:        userPath + "/info",
		Handler:     user.UserInfo,
		RequestType: "GET",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/reset/password",
		Handler:     user.ResetPassword,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/flow/total",
		Handler:     user.GetTotalFlowDetail,
		RequestType: "GET",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/flow/current",
		Handler:     user.GetCurrentFlowDetail,
		RequestType: "GET",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/flow/detail",
		Handler:     user.GetUserFlowDetail,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/list",
		Handler:     user.GetUserList,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/add",
		Handler:     user.AddUser,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/edit",
		Handler:     user.EditUser,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/distribute",
		Handler:     user.DistributeFlowToUser,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        userPath + "/distribute/log",
		Handler:     user.GetDistributeFlowLog,
		RequestType: "POST",
		JwtEnabled:  true,
	},
}

// 商品路由
var commoditiesRouter = []server.Router{}

var payRouter = []server.Router{}

// 系统公告路由
var announcementsRouter = []server.Router{
	{
		Path:        announcementsRouterPath + "/list",
		Handler:     announcements.GetAnnouncementsList,
		RequestType: "POST",
		JwtEnabled:  true,
	},
}

// 展示数据路由
var displayRouter = []server.Router{}

// 代理路由
var proxyRouter = []server.Router{}

// 充值路径
var rechargeRouter = []server.Router{
	{
		Path:        rechargePath + "/list",
		Handler:     recharge.GetRechargeList,
		RequestType: "POST",
		JwtEnabled:  true,
	},
}

// 子账户路径
var subuserRouter = []server.Router{
	{
		Path:        subuserPath + "/list",
		Handler:     subuser.GetSubuserList,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        subuserPath + "/add",
		Handler:     subuser.AddSubuser,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        subuserPath + "/edit",
		Handler:     subuser.EditSubuser,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        subuserPath + "/flowStats",
		Handler:     subuser.GetSubuserFlowStats,
		RequestType: "POST",
		JwtEnabled:  true,
	},
	{
		Path:        subuserPath + "/export",
		Handler:     subuser.ExportSubuserList,
		RequestType: "POST",
		JwtEnabled:  true,
	},
}

func mergeRouter(router ...[]server.Router) []server.Router {
	var routers []server.Router
	for _, r := range router {
		routers = append(routers, r...)
	}
	return routers
}

/*
初始化路由和Web服务监听
*/
func Setup() {
	port := viper.GetInt("web.port")
	server.Stop()
	go func() {
		server.InitGinEngine(
			viper.GetString("web.mode"),
			mergeRouter(
				userRouter,
				commoditiesRouter,
				displayRouter,
				announcementsRouter,
				payRouter,
				proxyRouter,
				subuserRouter,
				rechargeRouter,
			),
			viper.GetBool("web.recordLog"),
			viper.GetBool("web.recovery"),
			viper.GetBool("web.allowCors"),
			port,
			viper.GetInt("web.readTimeout"),
			viper.GetInt("web.weiteTimeout"),
		)
		err := server.Run()
		if err != nil && err.Error() != "http: Server closed" {
			log.Fatalln("接口服务启动失败: %v", err)
		}
	}()
	log.Info("接口服务已启动,端口号:[%d]", port)
}
