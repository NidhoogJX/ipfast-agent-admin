package services

import (
	"fmt"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
	"sync"
)

// 套餐时长类型表 秒 分 时 日  月  年
const (
	DurationTypeSecond = 1
	DurationTypeMinute = 2
	DurationTypeHour   = 3
	DurationTypeDay    = 4
	DurationTypeMonth  = 5
	DurationTypeYear   = 6
)

type DurationTypes = models.DurationTypes

type DurationTypeList struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Weight           int8    `json:"weight"`
	Time             string  `json:"time"`
	MultiplyingPower float64 `json:"multiplying_power"`
}

var durationTypeList = sync.Map{}

// var

// 定时更新 时长类型列表
func UpdateDurationTypeList() {
	data, err := models.DurationTypes{}.SelectAll()
	if err != nil {
		return
	}
	durationTypeList = sync.Map{}
	for _, v := range data {
		durationTypeList.Store(v.ID, v)
	}
}

// 检查传入的时长类型是否存在
func CheckTypeIsExist(id int64) bool {
	_, ok := durationTypeList.Load(id)
	return ok
}

func GetDurationTypeMap() map[int8]string {
	return map[int8]string{
		DurationTypeSecond: "秒钟",
		DurationTypeMinute: "分钟",
		DurationTypeHour:   "小时",
		DurationTypeDay:    "天",
		DurationTypeMonth:  "个月",
		DurationTypeYear:   "年",
	}
}

var durationtypeMap = map[int8]int64{
	DurationTypeSecond: 1,
	DurationTypeMinute: 60,
	DurationTypeHour:   60 * 60,
	DurationTypeDay:    60 * 60 * 24,
	DurationTypeMonth:  60 * 60 * 24 * 30,
	DurationTypeYear:   60 * 60 * 24 * 365,
}

// 根据类型ID获取时长类型
func GetDurationTypeList(typeId int8) (durationTypeList []DurationTypeList, err error) {
	if typeId == 1 {
		durationTypeList, err = GetDynamicDurationTypeList()
		if err != nil {
			log.Error("获取动态IP时长类型列表失败: %v", err)
			return
		}
	} else if typeId == 2 {
		durationTypeList, err = GetStaticDurationTypeList()
		if err != nil {
			log.Error("获取静态IP时长类型列表失败: %v", err)
			return
		}
	} else if typeId == 3 {
		durationTypeList, err = GetDataDurationTypeList()
		if err != nil {
			log.Error("获取数据中心IP时长类型列表失败: %v", err)
			return
		}
	} else {
		return nil, fmt.Errorf("类型ID错误: %d", typeId)
	}
	return
}

// GetDynamicDurationTypeList 获取动态IP时长类型列表
func GetDynamicDurationTypeList() ([]DurationTypeList, error) {
	var durationTypes []DurationTypeList
	durationTypeList.Range(func(key, value interface{}) bool {
		v := value.(DurationTypes)
		if v.CType == 1 {
			durationTypes = append(durationTypes, DurationTypeList{
				ID:               v.ID,
				Name:             v.Name,
				Weight:           v.Weight,
				MultiplyingPower: v.MultiplyingPower,
				Time:             fmt.Sprintf("%d%s", v.Count, GetDurationTypeMap()[v.Type]),
			})
		}
		return true
	})
	return durationTypes, nil
}

// GetStaticDurationTypeList 获取静态IP时长类型列表
func GetStaticDurationTypeList() ([]DurationTypeList, error) {
	var durationTypes []DurationTypeList
	durationTypeList.Range(func(key, value interface{}) bool {
		v := value.(DurationTypes)
		if v.CType == 2 {
			durationTypes = append(durationTypes, DurationTypeList{
				ID:               v.ID,
				Name:             v.Name,
				Weight:           v.Weight,
				MultiplyingPower: v.MultiplyingPower,
				Time:             fmt.Sprintf("%d%s", v.Count, GetDurationTypeMap()[v.Type]),
			})
		}
		return true
	})
	return durationTypes, nil
}

// 获取数据中心IP时长类型列表
func GetDataDurationTypeList() ([]DurationTypeList, error) {
	var durationTypes []DurationTypeList
	durationTypeList.Range(func(key, value interface{}) bool {
		v := value.(DurationTypes)
		if v.CType == 3 {
			durationTypes = append(durationTypes, DurationTypeList{
				ID:               v.ID,
				Name:             v.Name,
				Weight:           v.Weight,
				MultiplyingPower: v.MultiplyingPower,
				Time:             fmt.Sprintf("%d%s", v.Count, GetDurationTypeMap()[v.Type]),
			})
		}
		return true
	})
	return durationTypes, nil
}

func CalculateExpireTime(id int64) (timestamp int64, err error) {
	timestamp = 0
	value, ok := durationTypeList.Load(id)
	if !ok {
		err = fmt.Errorf("durationTypeId not found")
		return
	}
	// 类型断言并检测是否成功
	durationtype, ok := value.(DurationTypes)
	if !ok {
		err = fmt.Errorf("类型转换失败: %v", value)
	}
	timestamp = durationtypeMap[durationtype.Type] * durationtype.Count
	return
}

// 根据时长类型ID获取时长类型信息
func GetDurationType(id int64) (durationType DurationTypes, err error) {
	durationType, err = models.DurationTypes{
		ID: id,
	}.SelectByID()
	if err != nil {
		log.Error("获取时长类型信息失败: %v", err)
		return
	}
	return
}
