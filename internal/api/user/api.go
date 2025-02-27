package user

import (
	"fmt"
	"ipfast_server/internal/db/library"
	emailhandler "ipfast_server/internal/handler/emailHandler"
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
	"ipfast_server/pkg/util/log"
	"math/rand"
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
	agent, err := library.Auth(param.Username, param.Password)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	token, err := server.GenerateToken(fmt.Sprintf("%d", agent.Id))
	if err != nil {
		log.Error("GenerateToken error:%v", err)
		resp.Failed("auth failed, please contact the administrator")
		return
	}
	// 记录登录时间和ip
	err = services.RecordAgentLoginIpAndTime(agent, resp.Context.ClientIP())
	if err != nil {
		resp.Failed("failed to record login information")
		return
	}
	resp.Res["auth_token"] = token
	resp.Res["user_id"] = agent.Id
	resp.Res["app_key"] = agent.AppKey
	resp.Res["login_time"] = agent.LoginTime
	resp.Res["create_time"] = formatTimestamp(agent.CreateTime)
	resp.Success("operate success")
}

/*
获取代理商信息
*/
func UserInfo(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	agent, err := services.GetAgentInfo(uid)
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

// 获取整体流量统计
func GetTotalFlowDetail(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	flowDetail, err := services.GetTotalFlowDetail(uid)
	if err != nil {
		resp.Failed("failed to obtain user flow detail")
		return
	}
	resp.Res["flow_detail"] = flowDetail
	resp.Success("operate success")
}

// 获取当天流量统计
func GetCurrentFlowDetail(resp server.Response) {
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	flowDetail, err := services.GetCurrentFlowDetail(uid)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Res["current_flow_detail"] = flowDetail
	resp.Success("operate success")
}

// 获取用户流量明细
func GetUserFlowDetail(resp server.Response) {
	param := struct {
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.StartTime < 0 || param.EndTime < 0 || param.EndTime < param.StartTime {
		resp.Failed("param error")
		return
	}
	flowDate, err := services.GetUserFlowDetail(uid, param.StartTime, param.EndTime)
	if err != nil {
		resp.Failed("failed to obtain user flow detail")
		return
	}
	resp.Res["flow_date"] = flowDate
	resp.Success("operate success")
}

// 获取代理商的用户列表
func GetUserList(resp server.Response) {
	param := struct {
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
		Username   string `json:"username"`
		Status     int8   `json:"status"`
		TotalSort  int8   `json:"total_sort"`
		UsedSort   int8   `json:"used_sort"`
		EnableSort int8   `json:"enable_sort"`
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
	userList, total, err := services.GetUserListByPage(uid, param.Username, param.Page, param.PageSize, param.Status, param.TotalSort, param.UsedSort, param.EnableSort)
	if err != nil {
		resp.Failed("failed to obtain user list")
		return
	}
	resp.Res["user_list"] = userList
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 代理商添加用户
func AddUser(resp server.Response) {
	param := struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Description string `json:"description"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.Username == "" || param.Password <= "" {
		resp.Failed("param error")
		return
	}
	err = services.AddUser(uid, param.Username, param.Password, param.Description)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Success("operate success")
}

// 修改用户信息
func EditUser(resp server.Response) {
	param := struct {
		UserId      int64  `json:"user_id"`
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
	if uid <= 0 || param.UserId <= 0 || param.Status < 0 {
		resp.Failed("param error")
		return
	}
	err = services.EditUser(param.UserId, param.Password, param.Description, param.Status)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Success("operate success")
}

// 给用户分配流量
func DistributeFlowToUser(resp server.Response) {
	param := struct {
		UserId int64   `json:"user_id"`
		Count  float64 `json:"count"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 || param.UserId <= 0 || param.Count <= 0 {
		resp.Failed("param error")
		return
	}
	err = services.DistributeFlowToUser(uid, param.UserId, param.Count)
	if err != nil {
		resp.Failed(fmt.Sprintf("%v", err))
		return
	}
	resp.Success("operate success")
}

// 获取流量分配日志
func GetDistributeFlowLog(resp server.Response) {
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
	if uid <= 0 || param.Page <= 0 || param.PageSize <= 0 {
		resp.Failed("param error")
		return
	}
	distributeLogList, total, err := services.GetDistributeFlowLog(uid, param.Page, param.PageSize, param.Username)
	resp.Res["distribute_log_list"] = distributeLogList
	resp.Res["total"] = total
	resp.Success("operate success")
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

// 申请验证码
func ApplyVerificationCode(resp server.Response) {
	param := &struct {
		Email string `json:"email"`
	}{}
	err := resp.Json(param)
	if err != nil || param.Email == "" {
		resp.Failed("param error")
		return
	}
	if !library.IsValidEmail(param.Email) {
		resp.Failed("email format error")
		return
	}
	code, err := library.MakeVerificationCode(param.Email, resp.Context.ClientIP())
	if err != nil {
		resp.Failed("failed to generate verification code")
		return
	}
	emailhandler.AsyncSendEmail(emailhandler.EmailConfig{
		To:      []string{param.Email},
		Subject: "Verification code",
		Body:    fmt.Sprintf("Verification code:%s", code),
	})
	resp.Success("operate success")
}

// 根据开始-截止日期获取用户流量
func GetUserTraffic(resp server.Response) {
	param := struct {
		StartDate int64 `json:"start_time" binding:"required"`
		EndDate   int64 `json:"end_time" binding:"required"`
	}{}
	// 时间参数检查
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}

	if param.StartDate > param.EndDate {
		resp.Failed("param error")
		return
	}
	if param.StartDate > time.Now().Unix() || param.EndDate > time.Now().Unix() {
		resp.Failed("param error")
		return
	}

	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	trafficData1 := generateFakeTrafficData()
	trafficData3 := generateFakeTrafficData()
	resp.Res["static_traffic_data"] = trafficData1
	resp.Res["data_traffic_data"] = trafficData3

	trafficData2, err := services.GetUserFlowByDate(param.StartDate, param.EndDate, uid)
	if err != nil {
		resp.Failed("failed to obtain user traffic information")
		return
	}
	resp.Res["dynamic_traffic_data"] = trafficData2

	resp.Success("operate success")
}

// TrafficData 表示单日的流量数据
type TrafficData struct {
	Date  string `json:"date"`
	Bytes int64  `json:"bytes"`
}

// generateFakeTrafficData 生成最近 7 天的假数据
func generateFakeTrafficData() []TrafficData {
	var data []TrafficData
	now := time.Now()

	for i := 7; i > 0; i-- {
		date := now.AddDate(0, 0, -i).Format("20060102")
		bytes := rand.Int63n(1000000) // 随机生成流量数据
		data = append(data, TrafficData{
			Date:  date,
			Bytes: bytes,
		})
	}

	return data
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

/*
重置密码
*/
func ResetPassword(resp server.Response) {
	param := struct {
		UserId   int64  `json:"user_id"`  //用户ID
		Password string `json:"password"` //用户密码
	}{}
	err := resp.Bind(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	if param.Password == "" {
		resp.Failed("param error")
		return
	}
	uid := resp.GetUserID("user_id")
	if uid <= 0 {
		resp.Failed("param error")
		return
	}
	err = library.ResetPassword(param.UserId, param.Password)
	if err != nil {
		resp.Failed("password reset failed")
		return
	}
	resp.Success("operate success")

}
