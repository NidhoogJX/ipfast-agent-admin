package services

import (
	"fmt"
	"ipfast_server/internal/config/countrycode"
	"ipfast_server/internal/db/models"
	"ipfast_server/internal/handler/ipipgo"
	stripehandler "ipfast_server/internal/handler/stripeHandler"
	"ipfast_server/pkg/util/dingding"
	"ipfast_server/pkg/util/log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderItem = models.OrderItem
type TransactionOrders = models.TransactionOrders

// 拉起策略
type PullUpStrategy struct {
	Threshold int             // 阈值
	Intervals []time.Duration // 间隔
}

// 拉起记录
type UserPullUpRecord struct {
	Count          int       // 次数
	LastPullUpTime time.Time // 最后拉起时间
}

// 统一商品拉起订单参数
type Cids struct {
	Cid              int64 `json:"cid" binding:"required"`              // 商品id
	Quantity         int64 `json:"quantity" binding:"required"`         // 数量
	TrafficCountryId int64 `json:"country_id" binding:"required"`       // 国家id
	DurationTypeId   int64 `json:"duration_type_id" binding:"required"` // 时长类型ID
	CmiId            int64 `json:"cmi_id"`                              // 静态IP套餐ID(续费时传递)
}

// 订单交易状态
const (
	OrderStatusCreated   int8 = 0 // 已创建
	OrderStatusUnpaid    int8 = 1 // 未付款
	OrderStatusPaid      int8 = 2 // 已付款
	OrderStatusCompleted int8 = 3 // 交易完成
	OrderStatusTimeout   int8 = 4 // 交易超时
	OrderStatusFailed    int8 = 5 // 交易失败
)

// 订单超时时间 单位分钟
const ORDER_TIME_OUT = 15

// 三方支付统一返回信息
type ThirdPartyPreOrderInfo struct {
	Oid string `json:"oid"`     // 订单号
	Url string `json:"pay_url"` // 支付链接
}

// 用户的拉起记录全局缓存 ToDo 使用Redis记录
var UserPullUpRecords = make(map[string]UserPullUpRecord)

// =========================检查=========================

// 检查用户拉起次数限制
func CheckPullCount(uid string, strategy PullUpStrategy) error {
	record, exists := UserPullUpRecords[uid]
	if !exists {
		record.Count = 0
		record.LastPullUpTime = time.Now()
		UserPullUpRecords[uid] = record
	}
	if record.Count >= strategy.Threshold {
		intervalIndex := record.Count - strategy.Threshold
		if intervalIndex >= len(strategy.Intervals) {
			intervalIndex = len(strategy.Intervals) - 1
		}
		waitTime := strategy.Intervals[intervalIndex]
		if time.Since(record.LastPullUpTime) < waitTime {
			return fmt.Errorf("拉起次数过多，请等待 %v 后再试", waitTime)
		}
	}
	record.LastPullUpTime = time.Now()
	record.Count += 1
	UserPullUpRecords[uid] = record
	return nil
}

// 记录用户的续费订单拉起记录全局缓存 ToDo 使用Redis记录
var UserPullUpRenewRecords = make(map[int64]UserPullUpRecord)

// 检查用户拉取续费订单次数
func RenewCheckPullCount(ipId int64, strategy PullUpStrategy) error {
	record, exists := UserPullUpRenewRecords[ipId]
	if !exists {
		record.Count = 0
		record.LastPullUpTime = time.Now()
		UserPullUpRenewRecords[ipId] = record
	}
	if record.Count >= strategy.Threshold {
		intervalIndex := record.Count - strategy.Threshold
		if intervalIndex >= len(strategy.Intervals) {
			intervalIndex = len(strategy.Intervals) - 1
		}
		waitTime := strategy.Intervals[intervalIndex]
		if time.Since(record.LastPullUpTime) < waitTime {
			return fmt.Errorf("续费订单拉起次数过多，请等待 %v 后再试", waitTime)
		}
	}
	record.LastPullUpTime = time.Now()
	record.Count += 1
	UserPullUpRenewRecords[ipId] = record
	return nil
}

// 支付前检查 （静态IP）
func staticCommoditesCheck(cids []Cids) (err error) {
	var countryIds = map[int64]int64{}
	for _, comdite := range cids {
		countryIds[comdite.TrafficCountryId] = comdite.Quantity
	}
	// err = CheckAndCreateIPIPGO(uid)
	// if err != nil {
	// 	return
	// }
	return CheckAndUpdateStaticCountryStock(countryIds) //检查并更新静态住宅国家的库存
}

// 检查支付平台是否存在
func checkPayPlatform(pid int64) (payIdentification string, err error) {
	payPlatform, err := models.PayPlatform{PayPlatformID: pid, Status: 1}.SelectPayPlatformByID()
	if err != nil || payPlatform.PayPlatformID <= 0 {
		err = fmt.Errorf("pay platform does not exist%v", err)
		return
	}
	payIdentification = payPlatform.PayIdentification
	return
}

// 拉起预订单
func PullUpPreOrders(uid int64, cids []Cids, pid int64, strategy PullUpStrategy, comditeType int8) (payInfo ThirdPartyPreOrderInfo, err error) {
	var (
		OrderItemes []OrderItem // 订单项
		currency    string      // 货币
		amount      float64     // 金额
		desc        string      // 订单描述
	)
	payIdentification, err := checkPayPlatform(pid)
	if err != nil || payIdentification == "" {
		log.Error("%v", err)
		err = fmt.Errorf("[ORDER]pay platform does not exist")
		return
	}
	// 商品类型
	switch comditeType {
	case StaticIPType, DataIPType:
		err = staticCommoditesCheck(cids)
		if err != nil {
			err = fmt.Errorf("commodities check failed")
			return
		}
		OrderItemes, currency, err = makeStaticAndDataCommodities(cids, comditeType)
		desc = "Purchase Static IP-"
		if comditeType == DataIPType {
			desc = "Purchase Data Center IP-"
		}
	case DynamicIPType:
		OrderItemes, currency, err = makeDynamicCommodities(cids)
		desc = "Purchase Dynamic IP-" + OrderItemes[0].CommodityName + OrderItemes[0].Desc
	case RenewStaticIPType:
		OrderItemes, currency, err = makeRenewStaticCommodities(cids[0], RenewStaticIPType)
		desc = "Renew Static IP-" + OrderItemes[0].CommodityName
	case RenewDataIPType:
		OrderItemes, currency, err = makeRenewDataCommodities(cids[0], RenewDataIPType)
		desc = "Renew Data Center IP-" + OrderItemes[0].CommodityName
	default:
		err = fmt.Errorf("[ORDER]unknown order item type")
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	oid, amount, err := makeOrder(uid, OrderItemes, currency, desc)
	if err != nil {
		log.Error("[ORDER]make order failed: %v", err)
		err = fmt.Errorf("make order failed")
		return
	}
	payInfo, err = initiatePayment(oid, amount, desc, currency, payIdentification)
	if err != nil {
		log.Error("[ORDER]pull up pre order failed: %v", err)
		// 拉取支付失败，修改订单状态为交易失败
		order, err := models.TransactionOrders{Oid: oid}.GetByOid()
		if err != nil || order.Oid == "" {
			log.Error("[ORDER]查询订单失败: %v", err)
		}
		order.Status = OrderStatusFailed
		order.UpdatedTime = time.Now().Unix()
		err = models.TransactionOrders{}.UpdateOrderStatus(order)
		if err != nil {
			log.Error("[ORDER]更新订单状态失败: %v", err)
		}
		err = fmt.Errorf("pull up pre order failed")
		return payInfo, err
	}
	// 插入订单的支付链接
	err = models.TransactionOrders{
		Oid:         oid,
		PaymentLink: payInfo.Url,
		UpdatedTime: time.Now().Unix(),
	}.Update()
	if err != nil {
		log.Error("[ORDER]MakeOrder failed: %v", err)
		return
	}
	CheckOrderTimeout(oid)
	// if comditeType == StaticIPType || comditeType == RenewStaticIPType {
	// 	CheckStaticOrderTimeout(oid)
	// } else if comditeType == DynamicIPType {
	// 	CheckDynamicOrderTimeout(oid)
	// } else if comditeType == RenewDataIPType || comditeType == RenewDataIPType {
	// 	CheckDataCenterOrderTimeout(oid)
	// } else {
	// 	log.Error("[ORDER]unknown order item type")
	// }
	return
}

// 生成动态IP订单关联资产信息
func makeDynamicCommodities(commoditieIds []Cids) (OrderItemes []OrderItem, currency string, err error) {
	cidobj, err1 := GetCommodityByID(commoditieIds[0].Cid)
	if err1 != nil {
		err = err1
		return
	}
	currency = cidobj.Currency
	desc := ""
	if cidobj.Unit == 1 {
		desc = "GB"
	} else if cidobj.Unit == 2 {
		desc = "TB"
	} else if cidobj.Unit == 3 {
		desc = "PB"
	} else {
		err = fmt.Errorf("unit not found")
		return
	}
	OrderItemes = append(OrderItemes, OrderItem{
		Type:           DynamicIPType,
		CommodityID:    commoditieIds[0].Cid,
		CommodityName:  cidobj.Name,
		Quantity:       cidobj.TotalCount,
		Unit:           cidobj.Unit,
		Amount:         cidobj.Price,
		Desc:           fmt.Sprintf("(%d%s)", cidobj.TotalCount, desc),
		DurationTypeId: cidobj.DurationTypeId,
	})
	return
}

// 生成静态IP订单关联资产信息
func makeStaticAndDataCommodities(commoditieIds []Cids, comditeType int8) (OrderItemes []OrderItem, currency string, err error) {
	for _, comdite := range commoditieIds {
		if !CheckTypeIsExist(comdite.DurationTypeId) {
			err = fmt.Errorf("duration not found")
			return
		}
		// 根据时长类型ID和商品ID获取静态住宅国家商品
		cidobj, err1 := GetCountryCommoditieByCid(comdite.TrafficCountryId, comditeType)
		if err1 != nil {
			err = err1
			return
		}
		if cidobj.Currency == "" {
			currency = "USD"
		} else {
			currency = cidobj.Currency
		}
		OrderItemes = append(OrderItemes, OrderItem{
			Type:           comditeType,
			CommodityID:    cidobj.CommodityId,
			CommodityName:  cidobj.CommodityName,
			Quantity:       comdite.Quantity,
			Amount:         cidobj.Price * float64(comdite.Quantity),
			Desc:           cidobj.CountryName,
			AreaId:         comdite.TrafficCountryId,
			DurationTypeId: comdite.DurationTypeId,
		})
	}
	return
}

// 生成静态IP续费订单关联资产信息
func makeRenewStaticCommodities(commoditie Cids, comditeType int8) (OrderItemes []OrderItem, currency string, err error) {
	// 根据商品ID获取静态住宅国家商品
	cidobj, err1 := GetCountryCommoditieByCid(commoditie.TrafficCountryId, StaticIPType)
	if err1 != nil {
		err = err1
		return
	}
	// 根据商品时长ID获取倍率
	durationType, err := models.DurationTypes{
		ID: commoditie.DurationTypeId,
	}.SelectByID()
	if err != nil {
		err = fmt.Errorf("duration type not found")
		return
	}
	if cidobj.Currency == "" {
		currency = "USD"
	} else {
		currency = cidobj.Currency
	}
	OrderItemes = append(OrderItemes, OrderItem{
		Type:          comditeType,
		CommodityID:   cidobj.CommodityId,
		CommodityName: cidobj.CommodityName,
		Quantity:      1,
		Amount:        cidobj.Price * durationType.MultiplyingPower,
		Desc:          cidobj.CountryName,
		AreaId:        commoditie.TrafficCountryId,
		Ext1:          commoditie.CmiId,
	})
	return
}

// 生成数据中心IP续费订单关联资产信息
func makeRenewDataCommodities(commoditie Cids, comditeType int8) (OrderItemes []OrderItem, currency string, err error) {
	// 根据商品ID获取静态住宅国家商品
	cidobj, err1 := GetCountryCommoditieByCid(commoditie.TrafficCountryId, DataIPType)
	if err1 != nil {
		err = err1
		return
	}
	// 根据商品时长ID获取倍率
	durationType, err := models.DurationTypes{
		ID: commoditie.DurationTypeId,
	}.SelectByID()
	if err != nil {
		err = fmt.Errorf("duration type not found")
		return
	}
	if cidobj.Currency == "" {
		currency = "USD"
	} else {
		currency = cidobj.Currency
	}
	OrderItemes = append(OrderItemes, OrderItem{
		Type:          comditeType,
		CommodityID:   cidobj.CommodityId,
		CommodityName: cidobj.CommodityName,
		Quantity:      1,
		Amount:        cidobj.Price * durationType.MultiplyingPower,
		Desc:          cidobj.CountryName,
		AreaId:        commoditie.TrafficCountryId,
		Ext1:          commoditie.CmiId,
	})
	return
}

// 创建新订单
func makeOrder(uid int64, items []OrderItem, currency, desc string) (string, float64, error) {
	// 使用年月日时分秒 + cid + uid + 商品长度 组装成 oid
	now := time.Now()
	oid := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), uid, len(items))
	// 计算总金额
	var totalAmount float64
	for i, item := range items {
		totalAmount += item.Amount
		items[i].OrderID = oid
	}
	// 判断总金额是否达到最小支付金额
	if currency == "USD" && totalAmount < 0.5 {
		log.Error("[ORDER]MakeOrder failed: total amount less than 0.5 USD")
		return "", 0, fmt.Errorf("total amount less than 0.5 USD")
	}
	if currency == "CNY" && totalAmount < 3.74 { // 人民币/美元汇率有波动
		log.Error("[ORDER]MakeOrder failed: total amount less than 3.74 CNY")
		return "", 0, fmt.Errorf("total amount less than 3.74 CNY")
	}
	if currency == "HKD" && totalAmount < 4 {
		log.Error("[ORDER]MakeOrder failed: total amount less than 4 HKD")
		return "", 0, fmt.Errorf("total amount less than 4 HKD")
	}
	// 判断是否超过最大金额
	if totalAmount > 999999.99 {
		log.Error("[ORDER]MakeOrder failed: total amount greater than 999999.99")
		return "", 0, fmt.Errorf("total amount greater than 999999.99")
	}
	// 生成订单
	err := models.TransactionOrders{
		UserID:      uid,
		Oid:         oid,
		Status:      OrderStatusCreated,
		Desc:        desc,
		Amount:      totalAmount,
		Currency:    currency,
		CreatedTime: now.Unix(),
		UpdatedTime: now.Unix(),
		Items:       items,
		OrderType:   items[0].Type,
	}.Insert()
	if err != nil {
		log.Error("[ORDER]MakeOrder failed: %v", err)
		return "", 0, fmt.Errorf("order creation failed")
	}
	err = updateOrderUnpaid(oid)
	if err != nil {
		log.Error("[ORDER]MakeOrder failed: %v", err)
		return "", 0, fmt.Errorf("order creation failed")
	}
	return oid, totalAmount, nil
}

// 更新订单状态 支付平台 支付平台订单号
func updateOrderStatus(order *models.TransactionOrders, platform, tid string, status int8) error {
	order.Status = status
	order.Platform = platform
	order.UpdatedTime = time.Now().Unix()
	order.Tid = tid
	err := order.Update()
	if err != nil {
		log.Error("order update failed: %v, oid: %s, tid: %s", err, order.Oid, tid)
		return fmt.Errorf("order update failed")
	}
	return nil
}

// 更新用户资产
// func updateUserAssets(order *models.TransactionOrders, tx *gorm.DB) error {
// 	err := UpdateUserAssets(order.UserID, order.DurationTypeId, order.Items, tx)
// 	if err != nil {
// 		log.Error("update user assets failed: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
// 		return fmt.Errorf("update user assets failed")
// 	}
// 	return nil
// }

// 更新订单状态为未支付
func updateOrderUnpaid(oid string) (err error) {
	order, err := models.TransactionOrders{Oid: oid}.GetByOid()
	if err != nil {
		log.Error("order not found: %v, oid: %s", err, oid)
		err = fmt.Errorf("make order failed")
		return
	}
	order.Status = OrderStatusUnpaid
	err = order.Update()
	if err != nil {
		log.Error("order update failed: %v, oid: %s", err, oid)
		err = fmt.Errorf("make order failed")
		return
	}
	return
}

// 更新用户购买的静态IP和动态IP信息
func updateUserIpRecord(order *models.TransactionOrders, iptype int8, tx *gorm.DB) error {
	var datas = []models.IpRecord{}
	for _, item := range order.Items {
		now := time.Now().Unix()
		datas = append(datas, models.IpRecord{
			UserID:      order.UserID,
			Ip:          "分配中",
			Port:        0,
			OrderId:     order.Oid,
			Username:    "分配中",
			Password:    "分配中",
			Deadline:    now,
			Status:      3,
			CountryId:   item.AreaId,
			Type:        item.Type,
			CreatedTime: now,
			UpdatedTime: now,
		})
	}
	model := &models.IpRecord{}
	err := model.InsertBatch(datas, tx)
	if err != nil {
		log.Error("update user assets failed: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
		return fmt.Errorf("update user assets failed")
	}
	return nil
}

// var datas = []models.IpRecord{}

// func updateUserIpRecord(order *models.TransactionOrders, tx *gorm.DB) error {
// 	for _, item := range order.Items {
// 		now := time.Now().Unix()
// 		datas = append(datas, models.IpRecord{
// 			UserID:      order.UserID,
// 			Ip:          "分配中",
// 			Port:        0,
// 			OrderId:     order.Oid,
// 			Username:    "分配中",
// 			Password:    "分配中",
// 			Deadline:    now,
// 			Status:      3,
// 			CountryId:   item.AreaId,
// 			Type:        item.Type,
// 			CreatedTime: now,
// 			UpdatedTime: now,
// 		})
// 	}
// 	model := &models.IpRecord{}
// 	err := model.InsertBatch(datas, tx)
// 	if err != nil {
// 		log.Error("update user flow record failed: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
// 		return fmt.Errorf("update user flow record failed")
// 	}
// 	return nil
// }

// 根据订单记录的资产信息 获取静态国家数据 用于购买静态IP
func getCountryDataList(order *models.TransactionOrders) ([]ipipgo.CountryData, error) {
	var countryDataList []ipipgo.CountryData
	for _, item := range order.Items {
		country, err := GetStaticCountry(item.AreaId)
		if err != nil {
			log.Error("获取静态国家失败: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
			return nil, fmt.Errorf("get static country failed")
		}
		countryCode := countrycode.GetCountryCode(country.Code)
		if countryCode == "" {
			log.Error("获取静态国家代码失败: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
			return nil, fmt.Errorf("get static country code failed")
		}

		for range item.Quantity {
			countryDataList = append(countryDataList, ipipgo.CountryData{
				Country: countryCode,
				IpNum:   1,
				Region:  country.CityCode,
			})
		}
	}
	return countryDataList, nil
}

// 更新用户资产
// func UpdateUserAssets(uid int64, items []OrderItem, tx *gorm.DB) error {
// 	var assets []models.Assets
// 	for _, item := range items {
// 		var Unit int8
// 		// 查询商品信息
// 		if item.Type == DynamicIPType {
// 			_, err := GetCommodityByID(item.CommodityID)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			_, err := GetCountryCommoditieByCid(item.AreaId)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		// 计算过期时间
// 		expireTime, err := CalculateExpireTime(item.DurationTypeId)
// 		if err != nil {
// 			return err
// 		}
// 		now := time.Now().Unix()
// 		assets = append(assets, models.Assets{
// 			UserID:        uid,
// 			CommodityName: item.CommodityName,
// 			ExpireTime:    now + expireTime,
// 			Type:          item.Type,
// 			Status:        0,
// 			TotalCount:    item.Quantity,
// 			UsedCount:     0,
// 			Unit:          Unit,
// 			AreaId:        item.AreaId,
// 			CreatedTime:   now,
// 			UpdatedTime:   now,
// 		})
// 	}
// 	return InsertAssets(assets, tx)
// }

// 更新用户流量记录
func UpdateFlowRecord(uid, durationTypeId int64, items []OrderItem, tx *gorm.DB) error {
	var flowRecords []models.FlowRecord
	for _, item := range items {
		var Unit int8 // 单位类型:1(GB), 2(TB), 3(PB)
		// 查询商品信息,计算购买流量
		commodity, err := GetCommodityByID(item.CommodityID)
		if err != nil {
			log.Error("failed to get commodity: %v", err)
			return err
		}
		Unit = commodity.Unit
		purchasedFlow := commodity.TotalCount
		// 转换为字节单位
		purchasedFlow = purchasedFlow * 1024 * 1024 * 1024 // 转换为byte单位
		// 若Unit为2(TB) 或 3(PB)进行相应处理
		for i := 1; i < int(Unit); i++ {
			purchasedFlow *= 1024
		}
		// 计算过期时间
		expireTime, err := CalculateExpireTime(durationTypeId)
		if err != nil {
			return err
		}
		now := time.Now().Unix()
		flowRecords = append(flowRecords, models.FlowRecord{
			UserID:        uid,
			PurchasedFlow: purchasedFlow,
			UsedFlow:      0,
			Deadline:      now + expireTime,
			OrderId:       item.OrderID,
			CreatedTime:   now,
			UpdatedTime:   now,
		})
	}
	model := &models.FlowRecord{}
	err := model.InsertFlowRecords(flowRecords, tx)
	if err != nil {
		log.Error("update user flow record failed: %v", err)
		return fmt.Errorf("update user flow record failed")
	}
	return nil
}

// 更新用户IP记录
var ipRecords []models.IpRecord

func UpdateIPRecord(uid, durationTypeId int64, items []OrderItem, tx *gorm.DB) error {
	for _, item := range items {
		// 计算过期时间
		expireTime, err := CalculateExpireTime(durationTypeId)
		if err != nil {
			return err
		}
		now := time.Now().Unix()
		ipRecords = append(ipRecords, models.IpRecord{
			UserID:      uid,
			Ip:          "",
			Port:        0, // 0:分配中
			OrderId:     item.OrderID,
			Username:    "",
			Password:    "",
			Deadline:    now + expireTime,
			CountryId:   item.AreaId,
			Type:        item.Type,
			Status:      3, // 3:分配中
			CreatedTime: now,
			UpdatedTime: now,
		})
	}
	model := &models.IpRecord{}
	err := model.InsertBatch(ipRecords, tx)
	if err != nil {
		log.Error("update user ip record failed: %v", err)
		return fmt.Errorf("update user ip record failed")
	}
	return nil
}

// 查询正在启用的支付平台id，名称
func GetActivePayPlatform() (payPlatformList []models.PayPlatformList, err error) {
	payPlatformList, err = models.PayPlatform{}.SelectEnabledPayPlatform()
	if err != nil {
		log.Error("查询支付平台失败: %v", err)
		return nil, err
	}
	return payPlatformList, err
}

// 检测订单超时未支付
func CheckOrderTimeout(oid string) {
	ticker := time.NewTicker(time.Minute * ORDER_TIME_OUT)
	// 开启协程,避免阻塞
	go func() {
		defer ticker.Stop() // 确保定时器在协程结束时停止
		<-ticker.C
		retryCount := 3 // 重试次数
		for i := 0; i < retryCount; i++ {
			order, err := models.TransactionOrders{Oid: oid}.GetByOid()
			if err != nil || order.Oid == "" {
				log.Error("[ORDER]查询订单失败: %v", err)
				continue
			}
			if order.Status == OrderStatusPaid || order.Status == OrderStatusCompleted {
				log.Info("[ORDER]" + order.Oid + "号订单已支付，不需要处理")
				return
			}
			// 修改订单状态为交易超时
			order.Status = OrderStatusTimeout
			order.UpdatedTime = time.Now().Unix()
			err = models.TransactionOrders{}.UpdateOrderStatus(order)
			if err != nil {
				log.Error("[ORDER]更新订单状态失败: %v", err)
				continue
			}
			log.Info("[ORDER]" + order.Oid + "号订单超时未支付，已修改订单状态并返回库存")
			return
		}
		log.Error("[ORDER]订单处理失败，已达到最大重试次数")
	}()
}

// 检测订单超时未支付
func CheckStaticOrderTimeout(oid string) {
	ticker := time.NewTicker(time.Minute * ORDER_TIME_OUT)
	go func() {
		defer ticker.Stop() // 确保定时器在协程结束时停止
		<-ticker.C
		retryCount := 3 // 重试次数
		for i := 0; i < retryCount; i++ {
			order, err := models.TransactionOrders{Oid: oid}.GetByOid()
			if err != nil || order.Oid == "" {
				log.Error("[ORDER]查询订单失败: %v", err)
				continue
			}
			if order.Status == OrderStatusPaid || order.Status == OrderStatusCompleted {
				log.Info("[ORDER]" + order.Oid + "号订单已支付，不需要处理")
				return
			}
			// 查询该订单下的所有订单项
			orderItems, err := models.OrderItem{}.SelectOrderItemsByOid(oid)
			if err != nil {
				log.Error("[ORDER]查询订单项失败: %v", err)
				continue
			}
			for _, item := range orderItems {
				trafficCountry, err := models.TrafficCountry{}.GetById(item.AreaId)
				if err != nil {
					log.Error("[TrafficCountry]查询地区库存失败: %v", err)
					continue
				}
				trafficCountry.UsedQuantity -= item.Quantity
				if trafficCountry.UsedQuantity < 0 {
					log.Error("[TrafficCountry]库存数据有误")
					continue
				}
				// 更新库存
				err = models.TrafficCountry{}.UpdateUsedQuantity(trafficCountry)
				if err != nil {
					log.Error("[TrafficCountry]更新库存失败: %v", err)
					continue
				}
			}
			// 修改订单状态为交易超时
			order.Status = OrderStatusTimeout
			order.UpdatedTime = time.Now().Unix()
			err = models.TransactionOrders{}.UpdateOrderStatus(order)
			if err != nil {
				log.Error("[ORDER]更新订单状态失败: %v", err)
				continue
			}
			log.Info("[ORDER]" + order.Oid + "号订单超时未支付，已修改订单状态并返回库存")
			return
		}
		log.Error("[ORDER]订单处理失败，已达到最大重试次数")
	}()
}

// 支付成功后，用钉钉机器人发送通知
func SendDingTalkNotification(oid string) {
	// 定义标题
	title := "新的订单通知"
	// 查询订单信息
	orderModel := models.TransactionOrders{Oid: oid}
	order, err := orderModel.GetByOid()
	if err != nil {
		log.Error("查询订单信息失败: %v", err)
		return
	}
	if order.Oid == "" {
		log.Error("订单信息不存在: %v", err)
		return
	}
	// 查询用户信息
	uid := order.UserID
	userModel := models.User{Id: uid}
	user, err := userModel.FindById()
	if err != nil {
		log.Error("查询用户信息失败: %v", err)
		return
	}
	if user == nil {
		log.Error("用户信息不存在: %v", err)
		return
	}

	// 处理金额格式
	orderAmount := fmt.Sprintf("%.2f", order.Amount)
	// 转换时间格式
	timeObj := time.Unix(order.CreatedTime, 0)           // 将时间戳转换为时间对象
	createdTime := timeObj.Format("2006-01-02 15:04:05") // 格式化时间
	// 查询所有订单项信息
	orderItemModel := models.OrderItem{OrderID: oid}
	orderItems, err := orderItemModel.SelectOrderItemsByOid(oid)
	if err != nil {
		log.Error("查询订单项信息失败: %v", err)
		return
	}
	if len(orderItems) == 0 {
		log.Error("订单项信息不存在: %v", err)
		return
	}
	// 拼接订单项信息
	var orderItemsStr string
	for index, item := range orderItems {
		itemAmount := fmt.Sprintf("%.2f", item.Amount)
		orderItemsStr += fmt.Sprintf("\n\n"+"%d.商品名称:%s"+"\n\n"+"\x20数量:%d"+"\n\n"+"\x20单价:%s", index+1, item.CommodityName, item.Quantity, itemAmount)
	}
	// 构建通知内容
	message := fmt.Sprintf(`## 用户:%s,有新订单:`+"\n\n"+
		`- 用户ID: %d`+"\n\n"+
		`- 邮箱: %s `+"\n\n"+
		`- 订单号: %s`+"\n\n"+
		`- 交易平台订单号: %s`+"\n\n"+
		`- 订单描述: %s`+"\n\n"+
		`- 金额: %s `+"\n\n"+
		`- 货币: %s `+"\n\n"+
		`- 发起时间: %s`+"\n\n"+
		`- 订单详情: `+"\n\n"+
		`%s`,
		user.Name, uid, user.Email, oid, order.Tid, order.Desc, orderAmount, order.Currency, createdTime, orderItemsStr)
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=f7ce6c5198f0ce932bee9bef2d9f26bb5cfca977fa18d0866ab8a746c2e4a7b6"
	// 发送通知
	dingding.SendMarkdown(title, message, webhookURL, false)
}

// =========================查询=========================
// 分页信息获取订单
func GetOrderByUid(uid int64, oid string, order_type int8, page, pageSize int, startDate, endDate int64) (orders []models.TransactionOrders, total int64, err error) {
	return models.TransactionOrders{
		UserID:    uid,
		Oid:       oid,
		OrderType: order_type,
	}.GetOrderByUidAndOid(page, pageSize, startDate, endDate)
}

// 根据订单号查询订单信息
func getOrder(oid string) (*TransactionOrders, error) {
	order, err := TransactionOrders{Oid: oid}.GetByOid()
	if err != nil {
		log.Error("[ORDER]order not found: %v, oid: %s", err, oid)
		return nil, fmt.Errorf("order not found")
	}
	if order.Status != OrderStatusUnpaid {
		log.Error("Order Status Is %s", oid)
		return nil, fmt.Errorf("order status is not unpaid")
	}
	return &order, nil
}

// =========================支付=========================

// 拉起第三方支付 暂时只支持Stripe
func initiatePayment(oid string, amount float64, commodityName, currency, payIdentification string) (res ThirdPartyPreOrderInfo, err error) {
	res.Oid = oid
	switch payIdentification {
	case "stripe":
		res.Url, err = stripePay(oid, amount, commodityName, currency)
		if err != nil {
			err = fmt.Errorf("stripe pay failed")
		}
		return
	default:
		err = fmt.Errorf("unknown pay platform")
		return
	}
}

// 三方支付[Stripe]
func stripePay(oid string, amount float64, commodityName, currency string) (string, error) {
	return stripehandler.CreateCheckoutSession(
		&stripehandler.PayParam{
			Amount:        amount,
			CommodityName: commodityName,
			Currency:      currency,
			Quantity:      1,
			Oid:           oid,
		},
	)
}

type ResponseData struct {
	Oid         string `json:"oid"`
	PlatformOid string `json:"platform_oid"`
	Type        string `json:"type"`
	ServerId    int    `json:"server_id"`
	Time        int64  `json:"time"`
	Signature   string `json:"signature"`
	TryCount    int    `json:"try_count"`
	MaxCount    int    `json:"max_count"`
}

// 异步通知 回调处理
func Notify(oid, tid, platform string) (uid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("交易失败 Notify: %v,oid:%s,tid:%s", err, oid, tid)
		}
	}()
	order, err := getOrder(oid)
	if err != nil {
		return
	}
	err = updateOrderStatus(order, platform, tid, OrderStatusPaid)
	if err != nil {
		return
	}
	// 开启事务处理
	err = models.DB.Instance.Transaction(func(tx *gorm.DB) error {
		err = updateOrderStatus(order, platform, tid, OrderStatusCompleted) // 4. 更新订单状态为交易完成
		if err != nil {
			log.Error("更新订单状态失败,开始回滚数据: %v, uid:%d,oid: %s, tid: %s", err, order.UserID, oid, tid)
			return err
		}
		delete(UserPullUpRecords, fmt.Sprintf("%d", order.UserID)) // 5. 删除用户的拉起次数记录
		if order.Items[0].Type == StaticIPType || order.Items[0].Type == DataIPType {
			// 静态ip和数据中心，需更新ip记录
			err := UpdateIPRecord(order.UserID, order.Items[0].DurationTypeId, order.Items, tx)
			if err != nil {
				log.Error("update user ip record failed: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
				return fmt.Errorf("update user ip record failed")
			}
		} else if order.Items[0].Type == DynamicIPType {
			// 动态ip，需更新流量记录
			err := UpdateFlowRecord(order.UserID, order.Items[0].DurationTypeId, order.Items, tx)
			if err != nil {
				log.Error("update user flow record failed: %v, oid: %s, tid: %s", err, order.Oid, order.Tid)
				return fmt.Errorf("update user flow record failed")
			}
		}
		return nil
	})
	log.Info("交易成功: uid:%s,oid: %s, tid: %s", uid, oid, tid)
	return
}

// 获取用户未过期动态流量信息
func GetDynamicFlowRecords(userId int64) (flowData models.FlowData, err error) {
	flowData, err = models.FlowRecord{
		UserID:     userId,
		IdentityId: 1,
	}.GetUserFlow()
	if err != nil {
		log.Error("get dynamic flow record failed: %v", err)
		return models.FlowData{}, err
	}
	return flowData, nil
}

// 获取用户未过期静态IP/数据中心IP信息
func GetIPRecords(userId int64, ipType int8) (ipRecords []models.IpRecord, err error) {
	ipRecordModel := models.IpRecord{
		UserID: userId,
		Type:   ipType, // ip类型:2(静态IP), 3(数据中心IP)
	}
	ipRecords, err = ipRecordModel.GetByUserId()
	if err != nil {
		log.Error("get static ip or data center ip record failed: %v", err)
		return nil, err
	}
	return ipRecords, nil
}

// 根据订单获取商品名
func GetCommodityNameAndDeadline(uid int64) (commodityName string, deadline int64) {
	commodityData, err := models.FlowRecord{
		UserID:     uid,
		IdentityId: 1,
	}.GetUserFastExpire()
	if err != nil {
		log.Error("get dynamic flow record failed: %v", err)
		return "未知商品", 0
	}
	if commodityData.Desc != "" {
		commodityName = strings.Split(commodityData.Desc, "-")[1]
	}
	return commodityName, commodityData.Deadline
}

// 获取支付链接
func GetPayLink(oid string) (url string, err error) {
	orderModel := models.TransactionOrders{Oid: oid}
	order, err := orderModel.GetByOid()
	if err != nil {
		log.Error("get order failed: %v, oid: %s", err, oid)
		return "", err
	}
	// 判断订单是否超过支付时间(15min)
	now := time.Now().Unix()
	if now-order.CreatedTime > 900 {
		log.Error("the order has expired and payment cannot continue: %v, oid: %s", err, oid)
		return "", fmt.Errorf("订单已过期,无法继续支付")
	}
	if order.Status != OrderStatusUnpaid {
		log.Error("the order is not in an unpaid status and payment cannot continue: %v, oid: %s", err, oid)
		return "", fmt.Errorf("订单不是未付款状态,无法继续支付")
	}
	// 获取支付链接
	return order.PaymentLink, nil
}

// 检查订单状态
func CheckOrderStatus(oid string) (status int8, err error) {
	order, err := models.TransactionOrders{
		Oid: oid,
	}.GetByOid()
	if err != nil {
		err = fmt.Errorf("failed to obtain order information")
		return
	}
	// 判断是否过期并且未支付(15分钟)
	now := time.Now().Unix()
	if order.CreatedTime+ORDER_TIME_OUT*60 < now &&
		(order.Status != OrderStatusCompleted && order.Status != OrderStatusPaid && order.Status != OrderStatusFailed) {
		order.Status = OrderStatusTimeout
		order.UpdatedTime = now
		err = models.TransactionOrders{}.UpdateOrderStatus(order)
		if err != nil {
			log.Error("failed to update order status: %v, oid: %s", err, oid)
			return
		}
	}
	status = order.Status
	return
}

// 订单列表查询
func GetOrderList(page, pageSize int, email string) (orderList []models.OrderInfo, total int64, err error) {
	orderList, total, err = models.TransactionOrders{}.SelectOrderList(page, pageSize, email)
	if err != nil {
		err = fmt.Errorf("failed to get order list")
	}
	return
}

// 更新订单状态
func UpdateOrderStatus(oid string, status int8) (err error) {
	order := models.TransactionOrders{
		Oid:         oid,
		Status:      status,
		UpdatedTime: time.Now().Unix(),
	}
	err = models.TransactionOrders{}.UpdateOrderStatus(order)
	if err != nil {
		log.Error("failed to update order status: %v, oid: %s", err, oid)
	}
	return
}
