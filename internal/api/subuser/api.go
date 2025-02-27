package subuser

import (
	"fmt"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
)

// 获取子账户列表
func GetSubuserList(resp server.Response) {
	param := struct {
		Page        int    `json:"page"`
		PageSize    int    `json:"page_size"`
		UserId      int64  `json:"user_id"`
		SubuserName string `json:"subuser_name"`
		Status      int8   `json:"status"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Page <= 0 || param.PageSize <= 0 || (param.Status != 0 && param.Status != 1 && param.Status != 2) {
		resp.Failed("param error")
		return
	}
	subuserList, total, err := services.GetSubuserList(param.UserId, param.Page, param.PageSize, param.SubuserName, param.Status)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Res["subuser_list"] = subuserList
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 创建子账户
func AddSubuser(resp server.Response) {
	param := struct {
		UserId      int64   `json:"user_id"`
		Username    string  `json:"subuser_name"`
		Password    string  `json:"subuser_password"`
		MaxCapacity float64 `json:"subuser_max_capacity"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.UserId <= 0 || param.Username == "" || param.Password == "" || param.MaxCapacity < 0 {
		resp.Failed("param error")
		return
	}
	err = services.CreateSubUser(param.UserId, param.Username, param.Password, param.MaxCapacity)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Success("operate success")
}

// 编辑子账户
func EditSubuser(resp server.Response) {
	param := struct {
		SubuserId   int64   `json:"subuser_id"`
		Password    string  `json:"subuser_password"`
		MaxCapacity float64 `json:"subuser_max_capacity"`
		Status      int8    `json:"subuser_status"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.SubuserId <= 0 || param.Password == "" || param.MaxCapacity < 0 || (param.Status != 0 && param.Status != 1) {
		resp.Failed("param error")
		return
	}
	err = services.EditSubuser(param.SubuserId, param.Password, param.MaxCapacity, param.Status)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Success("operate success")
}

// 获取子账户流量统计
func GetSubuserFlowStats(resp server.Response) {
	param := struct {
		SubuserId int64 `json:"subuser_id"`
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.SubuserId <= 0 || param.StartTime <= 0 || param.EndTime <= 0 {
		resp.Failed("param error")
		return
	}
	flowStats, err := services.GetSubuserFlowStats(param.SubuserId, param.StartTime, param.EndTime)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Res["flow_stats"] = flowStats
	resp.Success("operate success")
}

// 导出子账户列表
func ExportSubuserList(resp server.Response) {
	param := struct {
		UserId int64 `json:"user_id"`
		Status int8  `json:"status"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.UserId <= 0 || (param.Status != 0 && param.Status != 1) {
		resp.Failed("param error")
		return
	}
	subuserList, err := services.ExportSubuserList(param.UserId, param.Status)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Res["subuser_list"] = subuserList
	resp.Success("operate success")
}
