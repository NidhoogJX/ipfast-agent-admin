package models

// 充值记录表
type Recharge struct {
	RechargeId  string  `json:"recharge_id" gorm:"primary_key;type:varchar(300);comment:充值订单号;not null;index"`
	AgentID     int64   `json:"agent_id" gorm:"type:bigint;comment:用户ID;not null;index"`
	Count       float64 `json:"count" gorm:"type:bigint;comment:充值流量;not null"`
	Unit        int8    `json:"unit" gorm:"type:tinyint;comment:单位:1GB,2TB,3PB;not null"`
	PayMethod   string  `json:"pay_method" gorm:"type:varchar(100);comment:支付方式;not null"`
	Description string  `json:"descriprion" gorm:"type:varchar(500);comment:订单描述;not null"`
	Status      int8    `json:"status" gorm:"type:tinyint;comment:订单状态:0充值失败,1充值成功;default:0;not null"`
	CreateTime  int64   `json:"create_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdateTime  int64   `json:"update_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
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
func (model Recharge) SelectRechargeList(agentId int64, rechargeId string, page, size int, status int8) (rechargeList []Recharge, total int64, err error) {
	tx := DB.Model(&model)
	if rechargeId != "" {
		tx.Where("recharge_id like ?", "%"+rechargeId+"%")
	}
	if status != 2 {
		tx.Where("status = ?", status)
	}
	err = tx.Select("*").
		Where("agent_id = ?", agentId).
		Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&rechargeList).Error
	return
}
