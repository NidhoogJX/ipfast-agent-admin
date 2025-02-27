package models

import (
	"time"

	"gorm.io/gorm"
)

// 流量记录表
type FlowRecord struct {
	Id            int64  `json:"id" gorm:"type:bigint; primary_key; AUTO_INCREMENT; not null; comment:流量记录ID"`
	UserID        int64  `json:"user_id" gorm:"type:bigint; not null; comment:用户ID"`
	Type          int8   `json:"type" gorm:"type:tinyint; not null; default:1;comment:流量记录类型(1:购买流量/管理员分配,2:代理商分配流量)"`
	AgentId       int64  `json:"agent_id" gorm:"type:bigint; not null; default:0; comment:代理商ID(非代理商分配为0)"`
	PurchasedFlow int64  `json:"purchased_flow" gorm:"type:bigint; not null; comment:购买流量(Byte)"`
	UsedFlow      int64  `json:"used_flow" gorm:"type:bigint; not null; default:0; comment:已使用流量(Byte)"`
	Deadline      int64  `json:"deadline" gorm:"type:bigint; not null; comment:流量有效期截止时间"`
	OrderId       string `json:"order_id" gorm:"type:varchar(100); not null; comment:订单ID(管理员添加的流量为0;代理商分配的为1)"`
	IdentityId    int8   `json:"identity_id" gorm:"type:tinyint; default:1; comment:用户身份ID(1:父用户,2:子用户)"`
	CreatedTime   int64  `json:"created_time" gorm:"type:bigint; not null; default:0; comment:创建时间"`
	UpdatedTime   int64  `json:"updated_time" gorm:"type:bigint; default:0; comment:更新时间"`
}

var FlowAccountUpdateFieid = []string{
	"id",
	"purchased_flow",
	"deadline",
	"order_id",
	"status",
	"updated_time",
}

// 流量日志
type DistributeFlowLog struct {
	UserId          int64  `json:"user_id"`
	Username        string `json:"username"`
	DistributeFlow  int64  `json:"distribute_flow"`
	BeforeFlow      int64  `json:"before_flow"`
	AfterFlow       int64  `json:"after_flow"`
	AgentEnableFlow int64  `json:"agent_enable_flow"`
	CreatedTime     int64  `json:"created_time"`
}

// TableName 设置表名
func (FlowRecord) TableName() string {
	return "ip_flow_record"
}

// 根据ID查询流量记录
func (model FlowRecord) GetById() (flowRecord FlowRecord, err error) {
	err = DB.Where("id = ?", model.Id).Find(&flowRecord).Error
	return flowRecord, err
}

// 批量添加流量记录
func (model FlowRecord) InsertFlowRecords(models []FlowRecord, tx *gorm.DB) error {
	return DB.Create(&models).Error
}

// 添加流量记录
func (model FlowRecord) Create() error {
	return DB.Create(&model).Error
}

// 更新流量记录
func (model FlowRecord) Update() error {
	err := DB.Model(&model).Where("id = ?", model.Id).Updates(model).Error
	return err
}

// 更新流量记录
func (model FlowRecord) UpdateByUsedFlow() error {
	err := DB.Model(&model).Where("id = ?", model.Id).UpdateColumn("used_flow", gorm.Expr("used_flow + ?", model.UsedFlow)).Error
	return err
}

// 删除流量记录
func (model *FlowRecord) Delete() error {
	return DB.Model(&model).Where("id = ?", model.Id).Delete(model).Error
}

// 根据用户ID查询所有未过期流量记录并根据deadline升序排列
func (model FlowRecord) GetByUserIds() (flowRecords []FlowRecord, err error) {
	err = DB.Where("user_id = ? and identity_id = ? and deadline > ? and used_flow < purchased_flow", model.UserID, model.IdentityId, time.Now().Unix()).Order("deadline ASC").Find(&flowRecords).Error
	return flowRecords, err
}

func (model FlowRecord) GetByUserSingel() (flowRecords FlowRecord, err error) {
	err = DB.Where("user_id = ? and identity_id = ? and deadline > ? and used_flow < purchased_flow", model.UserID, model.IdentityId, time.Now().Unix()).Order("deadline ASC").First(&flowRecords).Error
	return flowRecords, err
}

type FlowData struct {
	TotalFlow int64 `json:"total_flow"`
	UsedFlow  int64 `json:"used_flow"`
}

// 获取用户未过期流量
func (model FlowRecord) GetUserFlow() (flowData FlowData, err error) {
	err = DB.Instance.Raw(`SELECT SUM(purchased_flow) AS total_flow, SUM(used_flow) AS used_flow
		FROM ip_flow_recoord
		WHERE user_id = ? and identity_id = ? and deadline > ? 
		ORDER BY deadline ASC`,
		model.UserID, model.IdentityId, time.Now().Unix()).Scan(&flowData).Error
	return
}

type CommodityData struct {
	Desc     string `json:"desc"`
	Deadline int64  `json:"deadline"`
}

// 获取用户最快过期流量记录的时间和商品名
func (model FlowRecord) GetUserFastExpire() (commodityData CommodityData, err error) {
	err = DB.Instance.Raw(`SELECT
		ito.desc, ifr.deadline
	FROM
		ip_flow_recoord ifr
		LEFT JOIN ip_transaction_orders ito ON ifr.order_id = ito.oid 
	WHERE
		ifr.user_id = ?
		AND ifr.identity_id = ?
		AND ifr.deadline > ?
	ORDER BY
		ifr.deadline 
		LIMIT 1`,
		model.UserID, model.IdentityId, time.Now().Unix()).Scan(&commodityData).Error
	return
}

// 查询代理商的流量分配日志
func (model FlowRecord) SelectDistributeFlowLog(agentId int64, page, size int, username string) (distributeLogList []DistributeFlowLog, total int64, err error) {
	tx := DB.Table("ip_flow_record AS ifr")
	err = tx.Select(`
		ifr.user_id,
		iu.name AS username,
		ifr.purchased_flow AS distribute_flow,
		(SUM(ifr.purchased_flow) OVER (
			PARTITION BY ifr.user_id 
			ORDER BY ifr.created_time, ifr.id 
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) - ifr.purchased_flow) AS before_flow,
		SUM(ifr.purchased_flow) OVER (
			PARTITION BY ifr.user_id 
			ORDER BY ifr.created_time, ifr.id 
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) AS after_flow,
		ia.total_flow - SUM(ifr.purchased_flow) OVER (
			ORDER BY ifr.created_time, ifr.id 
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) AS agent_enable_flow,
		ifr.created_time
		`).
		Joins("LEFT JOIN ip_agent ia ON ifr.agent_id = ia.id").
		Joins("LEFT JOIN ip_user iu ON ia.id = iu.agent_id").
		Where("ifr.agent_id = ?", agentId).
		Order("ifr.created_time, ifr.id").
		Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&distributeLogList).Error
	return
}
