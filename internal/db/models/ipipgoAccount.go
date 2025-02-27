package models

/*
流量账号管理
*/
type IPIPGoAccount struct {
	Uid      int64  `json:"uid" gorm:"primary_key;type:bigint unsigned;not null;comment:用户ID"`
	Account  string `json:"id" gorm:"varchar(150);default:'';not null;comment:IPIPGO二级客户对应账号"`
	Sign     string `json:"sign" gorm:"type:varchar(150);default:'';not null;comment:IPIPGO二级客户对应签名"`
	AuthInfo string `json:"authInfo" gorm:"type:json;comment:认证套餐账密信息"`
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (IPIPGoAccount) TableName() string {
	return "ip_ipipgo_account"
}

// GetAuthInfo 获取IPIPGO二级客户认证套餐账密信息
func (m IPIPGoAccount) GetAccount() (data *IPIPGoAccount, err error) {
	data = &IPIPGoAccount{}
	err = DB.Where("uid = ?", m.Uid).Find(data).Error
	return
}

// Insert 插入IPIPGO二级客户认证套餐账密信息
func (m IPIPGoAccount) Insert() error {
	return DB.Create(m).Error
}
