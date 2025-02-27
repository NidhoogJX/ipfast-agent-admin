package pay

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/internal/handler/network/server"
	stripehandler "ipfast_server/internal/handler/stripeHandler"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/log"
	"time"
)

// 商品续费(静态/数据中心)
func IPRenew(resp server.Response) {
	param := struct {
		CountryId      int64 `json:"country_id" binding:"required"`       // 国家ID
		DurationTypeId int64 `json:"duration_type_id" binding:"required"` // 时长类型ID
	}{}
	if err := resp.Bind(&param); err != nil {
		log.Error("%v", err)
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	// 根据国家id获取商品信息
	commodities, err := services.GetCommodity(param.CountryId)
	if err != nil {
		log.Error("failed to get commodites info:%v", err)
		resp.Failed("failed to get commodites info")
		return
	}
	// 检查时长类型是否有效
	_, err = models.DurationTypes{
		ID: param.DurationTypeId,
	}.SelectByID()
	if err != nil {
		log.Error("failed to get duration type")
		resp.Failed("failed to get duration type")
		return
	}
	var ipType int8
	if commodities.Type == 2 {
		ipType = services.RenewStaticIPType
	} else if commodities.Type == 3 {
		ipType = services.RenewDataIPType
	} else {
		resp.Failed("param error")
		return
	}
	// 拉取预订单
	payInfo, err := services.PullUpPreOrders(
		uid,
		[]services.Cids{
			{
				Cid:              commodities.ID,
				Quantity:         1,
				TrafficCountryId: param.CountryId,
				DurationTypeId:   param.DurationTypeId,
			},
		}, // 商品列表
		1, // 支付平台ID,暂时只有一个，默认为1
		services.PullUpStrategy{
			Threshold: 3,
			Intervals: []time.Duration{
				time.Minute,
				5 * time.Minute,
				15 * time.Minute,
			},
		},
		ipType,
	)
	if err != nil {
		log.Error("failed to create renew order %v", err.Error())
		resp.Failed("failed to create renew order")
		return
	}
	resp.Res["pay_info"] = payInfo
	resp.Success("operate success")
}

/*
创建订单(静态)
*/
func StaticOrder(resp server.Response) {
	param := struct {
		Commodities []services.Cids `json:"commodities" binding:"required"` // 商品列表
	}{}

	if err := resp.Bind(&param); err != nil {
		log.Error("%v", err)
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || len(param.Commodities) == 0 {
		resp.Failed("param error")
		return
	}
	// 拉取预订单
	payInfo, err := services.PullUpPreOrders(
		uid,
		param.Commodities, // 商品列表
		1,                 // 支付平台ID
		services.PullUpStrategy{
			Threshold: 3,
			Intervals: []time.Duration{
				time.Minute,
				5 * time.Minute,
				15 * time.Minute,
			},
		},
		services.StaticIPType,
	)
	if err != nil {
		log.Error("failed to create static order %v", err.Error())
		resp.Failed("failed to create static order")
		return
	}
	resp.Res["pay_info"] = payInfo
	resp.Success("operate success")
}

/*
创建订单(数据中心)
*/
func DataOrder(resp server.Response) {
	param := struct {
		Commodities []services.Cids `json:"commodities" binding:"required"` // 商品列表
	}{}

	if err := resp.Bind(&param); err != nil {
		log.Error("%v", err)
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || len(param.Commodities) == 0 {
		resp.Failed("param error")
		return
	}
	// 拉取预订单
	payInfo, err := services.PullUpPreOrders(
		uid,
		param.Commodities, // 商品列表
		1,                 // 支付平台ID
		services.PullUpStrategy{
			Threshold: 3,
			Intervals: []time.Duration{
				time.Minute,
				5 * time.Minute,
				15 * time.Minute,
			},
		},
		services.DataIPType,
	)
	if err != nil {
		log.Error("failed to create data center order %v", err.Error())
		resp.Failed("failed to create data center order")
		return
	}
	resp.Res["pay_info"] = payInfo
	resp.Success("operate success")
}

/*
创建订单(动态)
*/
func DynamicOrder(resp server.Response) {
	param := struct {
		Cid int64 `json:"cid" binding:"required"`
	}{}

	if err := resp.Bind(&param); err != nil {
		log.Error("%v", err)
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	if param.Cid <= 0 || uid <= 0 {
		resp.Failed("param error")
		return
	}
	// 拉取预订单
	payInfo, err := services.PullUpPreOrders(
		uid,
		[]services.Cids{
			{
				Cid: param.Cid,
			},
		},
		1, // 支付平台ID
		services.PullUpStrategy{
			Threshold: 3,
			Intervals: []time.Duration{
				time.Minute,
				5 * time.Minute,
				15 * time.Minute,
			},
		},
		services.DynamicIPType,
	)
	if err != nil {
		log.Error("failed to create dynamic order %v", err.Error())
		resp.Failed("failed to create dynamic order")
		return
	}
	resp.Res["pay_info"] = payInfo
	resp.Success("operate success")
}

func Notify(resp server.Response) {
	var param services.ResponseData
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	log.Info("notify param:%v", param)
	if param.Type == stripehandler.CheckoutSessionCompleted || param.Type == stripehandler.CheckoutSessionAsyncPaymentSucceeded {
		go func() {
			// 校验签名
			signature := services.GenerateMD5(fmt.Sprintf("%s%s%d", param.Oid, param.Type, param.Time))
			if param.Signature != signature {
				resp.Failed("Signature verification failed")
				return
			}
			uid, err := services.Notify(param.Oid, param.PlatformOid, "stripe")
			if err != nil {
				log.Error("notify failed %v,%s", err, uid)
			}
		}()
	}

	// 支付成功后,用钉钉机器人发送通知
	services.SendDingTalkNotification(param.PlatformOid)

	resp.Response(200, 1, "success")
	resp.Context.Abort()
}

func RecordInfo(resp server.Response) {
	param := struct {
		Oid       string `json:"oid"`
		OrderType int8   `json:"order_type"`
		Page      int    `json:"page" binding:"required,min=1"`
		PageSize  int    `json:"page_size"  binding:"required,min=1,max=100"`
		StartDate int64  `json:"start_date"`
		EndDate   int64  `json:"end_date"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	if param.Page <= 0 || param.PageSize <= 0 {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	orderInfo, total, err := services.GetOrderByUid(uid, param.Oid, param.OrderType, param.Page, param.PageSize, param.StartDate, param.EndDate)
	if err != nil {
		resp.Failed("get order info failed")
		return
	}
	resp.Res["order_info"] = orderInfo
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 查询正在启用的支付平台id，名称
func GetActivePayPlatform(resp server.Response) {
	payPlatformList, err := services.GetActivePayPlatform()
	if err != nil {
		resp.Failed("get active pay platform failed")
		return
	}
	resp.Res["pay_platform_list"] = payPlatformList
	resp.Success("operate success")
}

// 获取支付链接
func TransferToPayLink(resp server.Response) {
	param := struct {
		Oid string `json:"oid"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	url, err := services.GetPayLink(param.Oid)
	resp.Res["pay_link"] = url
	if err != nil {
		resp.Failed("failed to obtain payment link")
		return
	}
	resp.Success("operate success")
}

// 检查订单状态
func CheckOrderStatus(resp server.Response) {
	param := struct {
		Oid string `json:"oid"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Oid == "" {
		resp.Failed("param error")
		return
	}
	status, err := services.CheckOrderStatus(param.Oid)
	if err != nil {
		resp.Failed("faile to obtain order status")
		return
	}
	resp.Res["status"] = status
	resp.Success("operate success")
}

// 订单列表查询
func GetOrderList(resp server.Response) {
	param := struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Email    string `json:"email"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Page < 0 || param.PageSize < 0 {
		resp.Failed("param error")
		return
	}
	orderList, total, err := services.GetOrderList(param.Page, param.PageSize, param.Email)
	if err != nil {
		resp.Failed("faile to obtain order list")
		return
	}
	resp.Res["total"] = total
	resp.Res["order_list"] = orderList
	resp.Success("operate success")
}

// 修改订单信息
func EditOrder(resp server.Response) {
	param := struct {
		OrderId string `json:"order_id"`
		Status  int8   `json:"status"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	err = services.UpdateOrderStatus(param.OrderId, param.Status)
	if err != nil {
		resp.Failed("faile to edit order info")
		return
	}
	resp.Success("operate success")
}
