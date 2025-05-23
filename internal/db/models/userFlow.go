package models

// 用户流量表
type UserFlow struct {
	UserID      int64 `json:"user_id" gorm:"type:bigint;comment:用户ID;not null;index"`
	SUserID     int64 `json:"s_user_id" gorm:"type:bigint;comment:子用户ID;not null;index"`
	Up          int64 `json:"up" gorm:"type:bigint;comment:上行流量;not null;default:0"`
	Down        int64 `json:"down" gorm:"type:bigint;comment:下行流量;not null;default:0"`
	CreatedTime int64 `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdatedTime int64 `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (UserFlow) TableName() string {
	return "ip_user_flow"
}

func (model UserFlow) Inserts(datas []UserFlow) error {
	return DB.Create(datas).Error
}

/*
根据用户ID查询用户流量
*/
func (model UserFlow) FindByUserID(userId int64) (userFlows []UserFlow, err error) {
	err = DB.Where("user_id = ?", userId).Find(userFlows).Error
	return
}

/*
根据开始-截止日期和用户ID查询用户流量
*/
func (model UserFlow) FindByDate(startTime, endTime, userId int64) (userFlows []struct {
	Bytes int64  `json:"bytes"`
	Date  string `json:"date"`
}, err error) {
	err = DB.Table(UserFlow{}.TableName()).
		Select("SUM(down) as bytes, FROM_UNIXTIME(created_time,'%Y%m%d') as date").
		Where("user_id = ?", userId).
		Where("created_time >= ? AND created_time <= ?", startTime, endTime).
		Group("date").
		Scan(&userFlows).Error
	return
}

// 子账户流量使用情况
type FlowStats struct {
	SUserID   int64      `json:"s_user_id"`
	SUsername string     `json:"s_username"`
	DateFlow  []DateFlow `json:"date_flow"`
}

type DateFlow struct {
	Date string `json:"date"`
	Flow int64  `json:"flow"`
}
