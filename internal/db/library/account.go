package library

import "ipfast_server/internal/db/models"

type Account = models.Account

/*
获取所有正常状态的账号列表作为白名单
*/
func GetAccountsByStatus(status int8) ([]models.Account, error) {
	model := &models.Account{}
	model.Status = status
	return model.SelectByAccountIdAndName()
}

/*
更新账号流量
*/
func UpdateAccountFlows(accounts []Account) error {
	model := &models.Account{}
	return model.UpdateFlows(accounts)
}
