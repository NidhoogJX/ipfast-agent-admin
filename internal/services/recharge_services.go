package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
)

// 代理商充值记录查询
func GetRechargeList(agentId int64, rechargeId string, page, size int, status int8) (rechargeList []models.Recharge, total int64, err error) {
	rechargeList, total, err = models.Recharge{}.SelectRechargeList(agentId, rechargeId, page, size, status)
	if err != nil {
		err = fmt.Errorf("failed to get recharge list")
	}
	return
}
