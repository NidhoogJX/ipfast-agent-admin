package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
)

// 代理商充值记录查询
func GetRechargeList(rechargeId string, page, size int) (rechargeList []models.RechargeInfo, total int64, err error) {
	rechargeList, total, err = models.Recharge{}.SelectRechargeList(rechargeId, page, size)
	if err != nil {
		err = fmt.Errorf("failed to get recharge list")
	}
	return
}
