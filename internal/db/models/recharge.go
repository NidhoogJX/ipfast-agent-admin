package models

// 充值记录表
type Recharge struct {
	Id          string `json:"id" gorm:"primary_key;type:varchar(300);comment:充值订单号;not null;index"`
	AgentID     int64  `json:"agent_id" gorm:"type:bigint;comment:用户ID;not null;index"`
	Count       int64  `json:"count" gorm:"type:bigint;comment:充值流量;not null"`
	Unit        int8   `json:"unit" gorm:"type:tinyint;comment:单位:1GB,2TB,3PB;not null"`
	PayMethod   string `json:"pay_method" gorm:"type:varchar(100);comment:支付方式;not null"`
	Description string `json:"descriprion" gorm:"type:varchar(500);comment:订单描述;not null"`
	Status      int8   `json:"status" gorm:"type:tinyint;comment:订单状态:0充值失败,1充值成功;default:0;not null"`
	CreateTime  int64  `json:"create_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdateTime  int64  `json:"update_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
}

type RechargeInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Count       int64  `json:"count"`
	BeforeFlow  int64  `json:"before_flow"`
	AfterFlow   int64  `json:"after_flow"`
	Unit        int8   `json:"unit"`
	Description string `json:"descriprion"`
	CreateTime  int64  `json:"create_time"`
}

// 获取表名
func (Recharge) TableName() string {
	return "ip_recharge"
}

// 创建充值记录
func (model Recharge) Create() (err error) {
	return DB.Create(&model).Error
}

// 获取充值记录列表
func (model Recharge) SelectRechargeList(rechargeId string, page, size int) (rechargeList []RechargeInfo, total int64, err error) {
	tx := DB.Table("ip_recharge AS ir")
	if rechargeId != "" {
		tx.Where("ir.id like ?", "%"+rechargeId+"%")
	}
	err = tx.Select(`
		ir.id,
		ia.name,
		ir.count,
		ir.unit,
		(SUM(ir.count) OVER (
			PARTITION BY ir.agent_id 
			ORDER BY ir.create_time, ir.id 
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) - ir.count) AS before_flow,
		SUM(ir.count) OVER (
			PARTITION BY ir.agent_id 
			ORDER BY ir.create_time, ir.id 
			ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		) AS after_flow,
		ir.description,
		ir.create_time
		`).
		Joins("LEFT JOIN ip_agent AS ia ON ir.agent_id = ia.id").
		Count(&total).
		Order("ir.create_time DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&rechargeList).Error
	return
}
