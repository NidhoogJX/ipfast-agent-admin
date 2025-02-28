package models

/*
IPFast 管理员信息
*/
type Admin struct {
	Id         int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:账号ID"`
	Name       string `json:"name" gorm:"type:varchar(150);default '';comment:账号名称;"`
	Password   string `json:"password" gorm:"type:varchar(150);not null;comment:密码"`
	AppKey     string `json:"app_key" gorm:"type:varchar(150);default '';comment:应用密钥"`
	Salt       string `json:"salt" gorm:"type:varchar(150);not null;comment:密码盐"`
	LoginIp    string `json:"login_ip" gorm:"type:varchar(50);default '';comment:上次登录IP"`
	LoginTime  int64  `json:"login_time" gorm:"type:bigint unsigned;default:0;not null;comment:上次登录时间"`
	Status     int8   `json:"status" gorm:"type:tinyint unsigned;default:0;not null;comment:账号状态"`
	CreateTime int64  `json:"create_time" gorm:"type:bigint unsigned;default:0;not null;comment:注册时间"`
	UpdateTime int64  `json:"update_time" gorm:"type:bigint unsigned;default:0;not null;comment:会员信息上次更新时间"`
}

// TableName 表名
func (Admin) TableName() string {
	return "ip_admin"
}

// 根据用户名查询管理员信息
func (model Admin) FindByName() (admin Admin, err error) {
	err = DB.Where("name = ?", model.Name).Find(&admin).Error
	return
}

// 更新管理员登录信息
func (model Admin) UpdateLoginInfo() error {
	return DB.Model(&model).Select("login_ip", "login_time").Updates(model).Error
}

// 根据ID查询管理员信息
func (model Admin) FindById() (admin Admin, err error) {
	err = DB.Where("id = ?", model.Id).First(&admin).Error
	return
}
