package recharge

import (
	"fmt"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
)

// 获取代理商的充值记录列表
func GetRechargeList(resp server.Response) {
	param := struct {
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
		RechargeId string `json:"recharge_id"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Page <= 0 || param.PageSize <= 0 {
		resp.Failed("param error")
		return
	}
	rechargeList, total, err := services.GetRechargeList(param.RechargeId, param.Page, param.PageSize)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Res["recharge_list"] = rechargeList
	resp.Res["total"] = total
	resp.Success("operate success")
}
