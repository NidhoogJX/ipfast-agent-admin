package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Assets = models.Assets
type AssetsRes = models.AssetsRes

// GetAssetsByUserID 根据用户ID获取资产
func GetAssetsByUserID(userID int64) (asset []AssetsRes, err error) {
	asset, err = Assets{UserID: userID}.GetAssetsByUserID()
	return
}

// GetAssetsByUserID 根据用户ID获取当前使用的资产
func GetUsedAssetsByUserID(userID int64) (asset []AssetsRes, err error) {
	asset, err = Assets{UserID: userID}.GetUsedAssetsByUserID()
	return
}

// InsertAssets 新增资产
func InsertAssets(assets []Assets, tx *gorm.DB) error {
	return Assets{}.InsertBatch(assets, tx)
}

type UsedAssets struct {
	TotalCount float64 `json:"total_count"`
	UsedCount  float64 `json:"used_count"`
	Unit       string  `json:"unit"`
}

// GetAssetsByUserID 根据用户ID获取资产汇总
func GetAllAssetsByUserID(userID int64) (usedAssets map[int8]UsedAssets, err error) {
	asset, err := Assets{UserID: userID}.GetAssetsByUserID()
	if err != nil {
		return
	}
	staticTotalCount, err := models.IpRecord{
		Type:   StaticIPType,
		UserID: userID,
	}.GetCountByType()
	if err != nil {
		return
	}
	dataTotalCount, err := models.IpRecord{
		Type:   DataIPType,
		UserID: userID,
	}.GetCountByType()
	if err != nil {
		return
	}
	usedAssets = map[int8]UsedAssets{
		1: {
			TotalCount: 0,
			UsedCount:  0,
			Unit:       "GB",
		},
		2: {
			TotalCount: float64(staticTotalCount),
			UsedCount:  0,
			Unit:       "个",
		},
		3: {
			TotalCount: float64(dataTotalCount),
			UsedCount:  0,
			Unit:       "个",
		},
	}
	for _, v := range asset {
		if v.Type == DynamicIPType {
			if _, ok := usedAssets[v.Type]; !ok {
				continue
			}
			ua := usedAssets[v.Type]
			ua.TotalCount += v.TotalCount
			ua.UsedCount += v.UsedCount
			usedAssets[v.Type] = ua
		}
	}
	return
}

var ConvertUnits float64 = 1024 * 1024 * 1024

// 根据用户id获取所有资产信息
func GetTotalAssetsByUserIDAll(userId int64) (usedAssets map[int8]UsedAssets, err error) {
	// 获取用户未过期动态流量信息
	var dynamicTotalCount float64
	var dynamicUsedCount float64
	flowData, err := GetDynamicFlowRecords(userId)
	if err != nil {
		log.Printf("获取用户有效动态流量信息失败%v", err)
		return
	}
	dynamicTotalCount = float64(flowData.TotalFlow) / ConvertUnits // 单位转换为GB,保留两位小数
	dynamicTotalCount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", dynamicTotalCount), 64)
	dynamicUsedCount = float64(flowData.UsedFlow) / ConvertUnits
	dynamicUsedCount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", dynamicUsedCount), 64)
	// 获取用户未过期静态ip信息
	staticIpRecords, err := GetIPRecords(userId, 2)
	if err != nil {
		log.Printf("获取用户有效静态ip信息失败%v", err)
		return
	}
	staticTotalCount := float64(len(staticIpRecords))
	staticUsedCount := float64(0)
	now := time.Now().Unix()
	expire := now + 3600*24*7 // 7天内为临期
	for _, ipRecord := range staticIpRecords {
		if ipRecord.Deadline < expire {
			staticUsedCount += 1
		}
	}
	// 获取用户未过期数据中心ip信息
	centerIpRecords, err := GetIPRecords(userId, 3)
	if err != nil {
		log.Printf("获取用户有效数据中心ip信息失败%v", err)
		return
	}
	centerTotalCount := float64(len(centerIpRecords))
	centerUsedCount := float64(0)
	expire = now + 3600*24*7 // 7天内为临期
	for _, ipRecord := range centerIpRecords {
		if ipRecord.Deadline < expire {
			centerUsedCount += 1
		}
	}
	// 封装数据
	usedAssets = map[int8]UsedAssets{
		1: {
			TotalCount: dynamicTotalCount,
			UsedCount:  dynamicUsedCount,
			Unit:       "GB",
		},
		2: {
			TotalCount: staticTotalCount,
			UsedCount:  staticUsedCount,
			Unit:       "个",
		},
		3: {
			TotalCount: centerTotalCount,
			UsedCount:  centerUsedCount,
			Unit:       "个",
		},
	}
	return
}
