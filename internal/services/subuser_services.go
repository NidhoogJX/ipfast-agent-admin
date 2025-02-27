package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"time"
)

// 子账户列表查询
func GetSubuserList(userId int64, page, size int, subuserName string, status int8) (subuserList []models.SubUserInfo, total int64, err error) {
	subuserList, total, err = models.SubUser{}.SelectSubUserListByAgentId(userId, page, size, subuserName, status)
	if err != nil {
		err = fmt.Errorf("failed to get subuser list")
	}
	return
}

// 添加子账户
func CreateSubUser(userId int64, username, password string, maxCapacity float64) (err error) {
	// 验证子用户名、密码格式
	err = CheckNameAndPassword(username, password)
	if err != nil {
		return
	}
	// 判断当前用户的子账户名是否存在
	count, err := models.SubUser{}.IsUsernameExist(username, userId)
	if err != nil {
		return fmt.Errorf("failed to check username exist")
	}
	if count > 0 {
		return fmt.Errorf("username already exists")
	}
	now := time.Now().Unix()
	maxStatus := int8(0)
	if maxCapacity > 0 {
		maxStatus = 1
	}
	model := models.SubUser{
		ParentUserID: userId,
		Username:     username,
		Password:     password,
		MaxCapacity:  maxCapacity,
		Remarks:      "",
		Status:       1,
		CreatedTime:  now,
		UpdatedTime:  now,
		MaxStatus:    maxStatus,
	}
	_, err = model.Create()
	if err != nil {
		err = fmt.Errorf("failed to create subuser")
	}
	return
}

// 修改子账户
func EditSubuser(subuserId int64, password string, maxCapacity float64, status int8) (err error) {
	// 验证子用户名、密码格式,跳过用户名
	err = CheckNameAndPassword("username", password)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	maxStatus := int8(0)
	if maxCapacity > 0 {
		maxStatus = 1
	}
	model := models.SubUser{
		Id:          subuserId,
		Password:    password,
		MaxCapacity: maxCapacity,
		Remarks:     "",
		Status:      status,
		UpdatedTime: now,
		MaxStatus:   maxStatus,
	}
	err = model.Update()
	if err != nil {
		err = fmt.Errorf("failed to update subuser")
	}
	return
}

// 获取子用户的流量使用情况
func GetSubuserFlowStats(subuserId int64, startTime, endTime int64) (subuserFlowData []models.DateFlow, err error) {
	subuserFlowData, err = models.UserFlow{}.SelectDynamicFlowStats(subuserId, startTime, endTime)
	if err != nil {
		err = fmt.Errorf("failed to get subuser flow stats")
	}
	return
}

// 导出子账户列表
func ExportSubuserList(userId int64, status int8) (subuserList []models.SubUserInfo, err error) {
	subuserList, err = models.SubUser{}.SelectSubUserListByStatus(userId, status)
	if err != nil {
		err = fmt.Errorf("failed to export subuser list")
	}
	return
}
