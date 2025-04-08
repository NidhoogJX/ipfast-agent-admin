package models

// 用户流量汇总表
type UserFlowHour struct {
	UserID      int64 `json:"user_id" gorm:"type:bigint;comment:用户ID;not null;uniqueIndex:index_user_suser_created"`
	SUserID     int64 `json:"s_user_id" gorm:"type:bigint;comment:子用户ID;not null;uniqueIndex:index_user_suser_created"`
	Up          int64 `json:"up" gorm:"type:bigint;comment:上行流量;not null;default:0;"`
	Down        int64 `json:"down" gorm:"type:bigint;comment:下行流量;not null;default:0"`
	CreatedTime int64 `json:"created_time" gorm:"type:bigint;comment:创建时间();not null;default:0;uniqueIndex:index_user_suser_created"`
	UpdatedTime int64 `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
}

// 获取表名
func (UserFlowHour) TableName() string {
	return "ip_user_flow_hour"
}

// 根据子用户id查询流量使用情况
func (model UserFlowHour) SelectSubUserFlowByDate(subUserId, startTime, endTime int64) (userFlows []DateFlow, err error) {
	err = DB.Model(&model).
		Select(`
			FROM_UNIXTIME(created_time, '%Y-%m-%d') AS date, 
			sum(up) AS flow
		`).
		Where("s_user_id = ? AND created_time BETWEEN ? AND ?", subUserId, startTime, endTime).
		Group("date").
		Scan(&userFlows).Error
	return
}
