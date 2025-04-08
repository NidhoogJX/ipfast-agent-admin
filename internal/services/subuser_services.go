package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"time"
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
	subuserFlowData, err = models.UserFlowHour{}.SelectSubUserFlowByDate(subuserId, startTime, endTime)
	if err != nil {
		err = fmt.Errorf("failed to get subuser flow stats")
	}
	subuserFlowData = FillMissingDates(startTime, endTime, subuserFlowData)
	return
}

// 填充缺失日期的流量数据
func FillMissingDates(startDate, endDate int64, flow []models.DateFlow) []models.DateFlow {
	start := time.Unix(startDate+8*3600, 0).Truncate(24 * time.Hour) // +8*3600:设置为东八区
	end := time.Unix(endDate+8*3600, 0).Truncate(24 * time.Hour)
	dateMap := make(map[string]int64)
	for _, f := range flow {
		dateMap[f.Date] = f.Flow
	}
	var result []models.DateFlow
	for current := start; !current.After(end); current = current.Add(24 * time.Hour) {
		dateStr := current.Format("2006-01-02")
		if flowValue, exists := dateMap[dateStr]; exists {
			result = append(result, models.DateFlow{Date: dateStr, Flow: flowValue})
		} else {
			result = append(result, models.DateFlow{Date: dateStr, Flow: 0})
		}
	}
	return result
}
