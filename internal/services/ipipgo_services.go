package services

import (
	"encoding/json"
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/internal/handler/ipipgo"
	"ipfast_server/pkg/util/log"
	"strings"
	"time"
)

type IPIPGoAccountStaticIp = models.IPIPGoAccountStaticIp

// 分页查询 用户的静态IP
// func StaticIPPaginate(page, pageSize int, ip, addressStr string, status int8, uid int64) (results []models.IPIPGoAccountStaticIp, total int64, err error) {
// 	return models.IPIPGoAccountStaticIp{}.Paginate(page, pageSize, ip, addressStr, status, uid)
// }

// 批量插入静态IP
func StaticIPInsert(datas []IPIPGoAccountStaticIp) error {
	mdoel := IPIPGoAccountStaticIp{}
	return mdoel.BatchInsert(datas)
}

// 检查对应的IPIPGO二级账号是否存在，不存在则创建
func CheckAndCreateIPIPGO(uid int64) (err error) {
	ipipgoAccount, err := models.IPIPGoAccount{
		Uid: uid,
	}.GetAccount()
	if err != nil {
		return
	}
	if ipipgoAccount == nil || ipipgoAccount.Account == "" {
		ipipgodata, err := ipipgo.CreateAccount()
		if err != nil {
			return err
		}
		authInfo, err := json.Marshal(&ipipgodata.Data.AuthInfo)
		if err != nil {
			return err
		}
		err = models.IPIPGoAccount{
			Uid:      uid,
			Account:  ipipgodata.Data.Account,
			Sign:     ipipgodata.Data.Sign,
			AuthInfo: string(authInfo),
		}.Insert()
		if err != nil {
			return err
		}
	}
	return
}

// 购买静态IP（IPIPGO）
func purchaseStaticIP(order *models.TransactionOrders) error {
	ipipgoAccount, err := models.IPIPGoAccount{Uid: order.UserID}.GetAccount()
	if err != nil {
		log.Error("获取IPIPGO账号失败: %v,uid:%d, oid: %s, tid: %s", err, order.UserID, order.Oid, order.Tid)
		return fmt.Errorf("get auth account failed")
	}
	if ipipgoAccount == nil || ipipgoAccount.Account == "" {
		log.Error("IPIPGO账号不存在: %v, uid:%d, oid: %s, tid: %s", err, order.UserID, order.Oid, order.Tid)
		return fmt.Errorf("auth account not exists")
	}
	countryDataList, err := getCountryDataList(order)
	if err != nil {
		log.Error("获取静态国家数据失败: %+v, uid:%d, oid: %s, tid: %s", err, order.UserID, order.Oid, order.Tid)
		return err
	}
	log.Info("开通静态IP国家: %+v, uid:%d, oid: %s, tid: %s", countryDataList, order.UserID, order.Oid, order.Tid)

	var staticIPlist = []ipipgo.StatciIPData{}
	for _, countryData := range countryDataList {
		res, err := ipipgo.StaticIPDetail(&ipipgo.StaticIPParam{ // 调用IPIPGO接口开通静态IP  开套餐并返回详细信息
			// MealId:          int(order.DurationTypeId), // 套餐ID
			MealId:          int(0), // 套餐ID
			CustomerName:    ipipgoAccount.Account,
			CountryDataList: []ipipgo.CountryData{countryData},
		})
		if err != nil && res != nil && res.Code != 0 {
			resData, _ := json.Marshal(res)
			log.Error("开通静态IP失败: %v,uid:%d,oid: %s, tid: %s", string(resData), order.UserID, order.Oid, order.Tid)
			return fmt.Errorf("open static ip failed")
		}
		staticIPlist = append(staticIPlist, res.Data)
		log.Info("开通静态IP成功: %v,uid:%d, oid: %s, tid: %s, IP: %+v", res, order.UserID, order.Oid, order.Tid, res.Data.IpList)
	}
	err = saveStaticIP(staticIPlist, order.Oid, order.Tid, order.UserID) // 保存开通后静态IP信息
	if err != nil {
		return err
	}
	log.Debug("保存静态IP信息成功: %v,uid:%d, oid: %s, tid: %s")
	return nil
}

func RenewStaticIP(order *models.TransactionOrders) error {
	ipipgoAccount, err := models.IPIPGoAccount{Uid: order.UserID}.GetAccount()
	if err != nil {
		log.Error("获取IPIPGO账号失败: %v,uid:%d, oid: %s, tid: %s", err, order.UserID, order.Oid, order.Tid)
		return fmt.Errorf("get auth account failed")
	}
	if ipipgoAccount == nil || ipipgoAccount.Account == "" {
		log.Error("IPIPGO账号不存在: %v, uid:%d, oid: %s, tid: %s", err, order.UserID, order.Oid, order.Tid)
		return fmt.Errorf("auth account not exists")
	}
	ipipgo.StaticIPRenewDetail(
		&ipipgo.StaticIPRenewParam{
			CustomerName: ipipgoAccount.Account,
			CmiId:        int64(order.Items[0].Ext1),
			// MealTime:     int64(order.DurationTypeId),
		},
	)
	return nil
}

// 保存静态IP信息
func saveStaticIP(ipdata []ipipgo.StatciIPData, oid, tid string, uid int64) error {
	var iplist []models.IPIPGoAccountStaticIp
	now := time.Now().Unix()
	for _, ipres := range ipdata {
		for _, item := range ipres.IpList {
			accountPassword := strings.Split(item.AccountPassword, "|")
			iplist = append(iplist, models.IPIPGoAccountStaticIp{
				OrderID: oid,
				Uid:     uid,
				IpPort:  item.IpPort,
				//AccountPassword: item.AccountPassword,
				Account:     accountPassword[0],
				Password:    accountPassword[1],
				AddressStr:  item.AddressStr,
				Status:      int8(item.Status),
				EndTime:     item.EndTime,
				MealId:      item.MealId,
				CmiId:       int64(item.CmiId),
				CreatedTime: now,
				UpdatedTime: now,
			})
		}
	}
	err := StaticIPInsert(iplist)
	if err != nil {
		log.Error("开通静态IP成功记录IP数据失败: %v, oid: %s, tid: %s", err, oid, tid)
		return fmt.Errorf("open static ip failed")
	}
	return nil
}
