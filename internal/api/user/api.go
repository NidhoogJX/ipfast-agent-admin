package user

import (
	"fmt"
	"ipfast_server/internal/db/library"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/log"
	"time"
)

/*
账号密码登录
*/
func Login(resp server.Response) {
	param := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := resp.Json(&param)
	if err != nil || param.Username == "" || param.Password == "" {
		resp.Failed("param error")
		return
	}
	// 账号密码登录
	admin, err := library.Auth(param.Username, param.Password)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	token, err := server.GenerateToken(fmt.Sprintf("%d", admin.Id))
	if err != nil {
		log.Error("GenerateToken error:%v", err)
		resp.Failed("auth failed, please contact the administrator")
		return
	}
	// 记录登录时间和ip
	err = services.RecordAdminLoginIpAndTime(admin, resp.Context.ClientIP())
	if err != nil {
		resp.Failed("failed to record login information")
		return
	}
	resp.Res["auth_token"] = token
	resp.Res["user_name"] = admin.Name
	resp.Res["user_id"] = admin.Id
	resp.Res["login_time"] = admin.LoginTime
	resp.Res["create_time"] = formatTimestamp(admin.CreateTime)
	resp.Success("operate success")
}

/*
获取管理员信息
*/
func UserInfo(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	agent, err := services.GetAdminByUserId(uid)
	if err != nil {
		resp.Failed("failed to obtain current user information")
		return
	}
	resp.Res["user_name"] = agent.Name
	resp.Res["user_app_key"] = agent.AppKey
	resp.Res["user_id"] = agent.Id
	resp.Res["user_login_time"] = agent.LoginTime
	resp.Res["create_time"] = formatTimestamp(agent.CreateTime)
	resp.Success("operate success")
}

// formatTimestamp 格式化时间戳为日期字符串
func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

/*
账号退出登陆

	param resp: base.Response 响应对象
*/
func LoginOut(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	log.Info("userId:%d login out", uid)
	resp.Success("operate success")
}

// 查询代理商的用户列表
func GetUserList(resp server.Response) {
	param := struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		AgentId  int64  `json:"agent_id"`
		Username string `json:"username"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Page <= 0 || param.PageSize <= 0 || param.AgentId <= 0 {
		resp.Failed("param error")
		return
	}
	userList, total, err := services.GetUserList(param.Page, param.PageSize, param.AgentId, param.Username)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err.Error()))
		return
	}
	resp.Res["user_list"] = userList
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 子用户列表分页查询
func GetSubUserList(resp server.Response) {
	param := struct {
		Page         int    `json:"page"`
		PageSize     int    `json:"page_size"`
		SubUsrerName string `json:"user_name"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Page <= 0 || param.PageSize <= 0 {
		resp.Failed("param error")
		return
	}
	data, total, err := services.GetSubUserByPage(param.Page, param.PageSize, uid, param.SubUsrerName)
	if err != nil {
		resp.Failed("get subUser failed")
		return
	}
	resp.Res["sub_user_list"] = data
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 批量删除子用户
func DeleteSubUserBySubUserIds(resp server.Response) {
	param := struct {
		Ids []int64 `json:"ids" binding:"required"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	err = services.DeleteSubUserByIDs(param.Ids, uid)
	if err != nil {
		resp.Failed("failed to delete subUser")
		return
	}
	resp.Success("operate success")
}

// 获取认证账号列表
func GetSubUserByEnableStatus(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	data, err := services.SelectSubUserByEnableStatus(uid)
	if err != nil {
		resp.Failed("get subUser failed")
		return
	}
	type Res struct {
		Id       int64  `json:"id"`
		UserName string `json:"username"`
		PassWord string `json:"password"`
	}
	var sub_user_list = []Res{}
	for _, v := range data {
		sub_user_list = append(sub_user_list, Res{
			Id:       v.Id,
			UserName: v.Username,
			PassWord: v.Password,
		})
	}
	resp.Res["sub_user_list"] = sub_user_list
	resp.Success("operate success")
}

// 根据用户名查询子用户信息
func GetSubUserByUsername(resp server.Response) {
	param := struct {
		Username string `json:"username"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	data, err := services.SelectSubUserByUsername(param.Username, uid)
	if err != nil {
		resp.Failed("get subUser failed")
		return
	}
	resp.Res["sub_user_info"] = data
	resp.Success("operate success")
}

// 添加子账号
func AddSubUser(resp server.Response) {
	param := struct {
		Username    string  `binding:"required,min=3,max=20"`
		Password    string  `json:"password" binding:"required,min=6,max=20"`
		Status      int8    `json:"status" binding:"oneof=0 1"`
		MaxStatus   int8    `json:"max_status" binding:"oneof=0 1"`
		MaxCapacity float64 `json:"max_capacity" binding:"gte=0"`
		Remarks     string  `json:"remarks" binding:"max=255"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || (param.MaxStatus == 1 && param.MaxCapacity < 0) {
		resp.Failed("param error")
		return
	}
	err = services.AddSubUser(
		uid,
		param.Username,
		param.Password,
		param.MaxCapacity,
		param.Status,
		param.MaxStatus,
		param.Remarks,
	)
	if err != nil {
		resp.Failed("failed to add subUser")
		return
	}
	resp.Success("operate success")
}

// 更新子账号信息
func UpdateSubUser(resp server.Response) {
	param := struct {
		SubUserId   int64   `json:"id" binding:"required"`
		Username    string  `binding:"required,min=3,max=20"`
		Password    string  `json:"password" binding:"required,min=6"`
		Status      int8    `json:"status" binding:"oneof=0 1"`
		MaxStatus   int8    `json:"max_status" binding:"oneof=0 1"`
		MaxCapacity float64 `json:"max_capacity" binding:"gte=0"`
		Remarks     string  `json:"remarks" binding:"max=255"`
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	err = services.UpdateSubUser(
		uid,
		param.SubUserId,
		param.Username,
		param.Password,
		param.MaxCapacity,
		param.Remarks,
		param.Status,
		param.MaxStatus,
	)
	if err != nil {
		resp.Failed("failed to update subUser")
		return
	}
	resp.Success("operate success")
}
