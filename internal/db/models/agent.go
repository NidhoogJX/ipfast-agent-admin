package models

// 代理商管理
type Agent struct {
	Id             int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:账号ID"`
	Name           string `json:"name" gorm:"type:varchar(150);default '';comment:账号名称;"`
	Email          string `json:"email" gorm:"type:varchar(150);default:'';not null;comment:邮箱;index"`
	Phone          string `json:"phone" gorm:"type:varchar(150);default:'';not null;comment:手机号;index"`
	Password       string `json:"password" gorm:"type:varchar(150);not null;comment:密码"`
	AppKey         string `json:"app_key" gorm:"type:varchar(150);default '';comment:应用密钥"`
	Salt           string `json:"salt" gorm:"type:varchar(150);not null;comment:密码盐"`
	TotalFlow      int64  `json:"total_flow" gorm:"type:bigint unsigned;default:0;not null;comment:总流量"`
	DistributeFlow int64  `json:"distribute_flow" gorm:"type:bigint unsigned;default:0;not null;comment:已分配流量"`
	LoginIp        string `json:"login_ip" gorm:"type:varchar(50);default '';comment:上次登录IP"`
	LoginTime      int64  `json:"login_time" gorm:"type:bigint unsigned;default:0;not null;comment:上次登录时间"`
	Status         int8   `json:"status" gorm:"type:tinyint unsigned;default:0;not null;comment:账号状态"`
	Description    string `json:"description" gorm:"type:varchar(500);comment:备注"`
	CreateTime     int64  `json:"create_time" gorm:"type:bigint unsigned;default:0;not null;comment:注册时间"`
	UpdateTime     int64  `json:"update_time" gorm:"type:bigint unsigned;default:0;not null;comment:会员信息上次更新时间"`
}

// 整体流量明细
type AgentFlowInfo struct {
	TotalFlow      int64 `json:"total_flow"`      // 总流量
	DistributeFlow int64 `json:"distribute_flow"` // 已分配流量
	EnableFlow     int64 `json:"enable_flow"`     // 剩余流量
	UsedFlow       int64 `json:"used_flow"`       // 已消耗流量
}

// 当天流量明细
type CurrentAgentFlowInfo struct {
	UserCount      int64 `json:"user_count"`      // 用户数量
	DistributeFlow int64 `json:"distribute_flow"` // 已分配流量
	UsedFlow       int64 `json:"used_flow"`       // 已消耗流量
}

// 获取表名
func (Agent) TableName() string {
	return "ip_agent"
}

// 根据邮箱查询代理商信息
func (model Agent) FindByEmail() (agent Agent, err error) {
	err = DB.Model(&model).Select("*").Where("email = ?", model.Email).First(&agent).Error
	return
}

// 根据用户名查询代理商信息
func (model Agent) FindByName() (agent Agent, err error) {
	err = DB.Model(&model).Select("*").Where("name = ?", model.Name).First(&agent).Error
	return
}

// 根据id查询代理商信息
func (model Agent) FindById() (agent Agent, err error) {
	err = DB.Model(&model).Select("*").Where("id = ?", model.Id).First(&agent).Error
	return
}

// 更新代理商流量信息
func (model Agent) UpdateFlowInfo() error {
	return DB.Model(&model).Select("distribute_flow", "update_time").Updates(model).Error
}

// 更新代理商登录信息
func (model Agent) UpdateLoginInfo() error {
	return DB.Model(&model).Select("login_ip", "login_time").Updates(model).Error
}

// 查询代理商的流量统计
func (model Agent) SelectTotalFlowInfo(id int64) (agentFlowInfo AgentFlowInfo, err error) {
	tx := DB.Table("ip_agent AS ia")
	err = tx.Select(`
		ia.total_flow,
		ia.distribute_flow,
		ia.total_flow - ia.distribute_flow AS enable_flow,
		SUM(ifr.used_flow) AS used_flow
		`).
		Joins("LEFT JOIN ip_flow_record AS ifr ON ia.id = ifr.agent_id AND ifr.type = 2").
		Where("ia.id = ?", id).
		Group("ia.id").
		Scan(&agentFlowInfo).Error
	return
}

// 查询代理商当天的流量统计
func (model Agent) SelectCurrentFlowInfo(id int64, startTime, endTime int64) (currentAgentFlowInfo CurrentAgentFlowInfo, err error) {
	tx := DB.Table("ip_flow_record")
	err = tx.Select(`
		(SELECT
			COUNT(*) 
		FROM 
			ip_user 
		WHERE 
			agent_id = 1
		) AS user_count,
		SUM(purchased_flow) AS distribute_flow,
		SUM(used_flow)
		`).
		Where("agent_id = ? AND created_time BETWEEN ? AND ?", id, startTime, endTime).
		Scan(&currentAgentFlowInfo).Error
	return
}

// 查询代理商下用户的流量明细
func (model Agent) SelectUserFlowInfo(id int64, startTime, endTime int64) (flowDate []DateFlow, err error) {
	err = DB.Instance.Raw(`
		SELECT 
			DATE_FORMAT(a.date, "%Y-%m-%d") as date, 
			IFNULL(SUM(b.up), 0) as flow 
		FROM 
			(SELECT CURDATE() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY as date 
			 FROM (SELECT 0 as a union all SELECT 1 union all SELECT 2 union all SELECT 3 union all SELECT 4 union all SELECT 5 union all SELECT 6 union all SELECT 7 union all SELECT 8 union all SELECT 9) as a 
			 CROSS JOIN (SELECT 0 as a union all SELECT 1 union all SELECT 2 union all SELECT 3 union all SELECT 4 union all SELECT 5 union all SELECT 6 union all SELECT 7 union all SELECT 8 union all SELECT 9) as b 
			 CROSS JOIN (SELECT 0 as a union all SELECT 1 union all SELECT 2 union all SELECT 3 union all SELECT 4 union all SELECT 5 union all SELECT 6 union all SELECT 7 union all SELECT 8 union all SELECT 9) as c) a 
		LEFT JOIN ip_user_flow b 
			ON FROM_UNIXTIME(b.created_time, "%Y-%m-%d") = a.date 
			AND b.user_id IN (SELECT id FROM ip_user WHERE agent_id = ?)
		WHERE a.date BETWEEN FROM_UNIXTIME(?, "%Y-%m-%d") AND FROM_UNIXTIME(?, "%Y-%m-%d") 
		GROUP BY a.date
		ORDER BY a.date 
	`, id, startTime, endTime).Scan(&flowDate).Error
	return
}
