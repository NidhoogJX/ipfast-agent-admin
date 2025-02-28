package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
)

// 子账户列表查询
func GetSubuserList(userId int64, page, size int, subuserName string) (subuserList []models.SubUserInfo, total int64, err error) {
	subuserList, total, err = models.SubUser{}.SelectSubUserListByUserId(userId, page, size, subuserName)
	if err != nil {
		err = fmt.Errorf("failed to get subuser list")
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
