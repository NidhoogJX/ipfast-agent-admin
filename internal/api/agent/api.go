package agent

import (
	"fmt"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
)

// 代理商列表查询
func GetAgentList(resp server.Response) {
	param := struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Username string `json:"username"`
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
	agentList, total, err := services.GetAgentList(param.Page, param.PageSize, param.Username)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err.Error()))
		return
	}
	resp.Res["agent_list"] = agentList
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 添加代理商信息
func AddAgent(resp server.Response) {
	param := struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Description string `json:"description"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Username == "" || param.Password == "" {
		resp.Failed("param error")
		return
	}
	err = services.AddAgent(param.Username, param.Password, param.Description)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err.Error()))
		return
	}
	resp.Success("operate success")
}

// 编辑代理商信息
func EditAgent(resp server.Response) {
	param := struct {
		AgentId     int64  `json:"agent_id"`
		Password    string `json:"password"`
		Status      int8   `json:"status"`
		Description string `json:"description"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.AgentId <= 0 || param.Password == "" || (param.Status != 0 && param.Status != 1) {
		resp.Failed("param error")
		return
	}
	err = services.EditAgent(param.AgentId, param.Password, param.Description, param.Status)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err.Error()))
		return
	}
	resp.Success("operate success")
}

// 给代理商充值流量
func RechargeFlowToAgent(resp server.Response) {
	param := struct {
		AgentId     int64  `json:"agent_id"`
		Count       int64  `json:"count"`
		Sign        string `json:"sign"`
		Description string `json:"description"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.AgentId <= 0 || (param.Sign != "+" && param.Sign != "-") || param.Count <= 0 {
		resp.Failed("param error")
		return
	}
	err = services.RechargeFlowToAgent(param.AgentId, param.Count, param.Sign, param.Description)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err.Error()))
		return
	}
	resp.Success("operate success")
}
