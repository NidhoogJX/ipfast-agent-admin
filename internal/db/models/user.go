package models

import (
	"gorm.io/gorm"
)

/*
IPFast 用户管理
*/
type User struct {
	Id          int64   `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:账号ID"`
	Money       float64 `json:"money" gorm:"type:double(10,2);default:0;not null;comment:账号余额"`
	Name        string  `json:"name" gorm:"type:varchar(150);default '';comment:账号名称;"`
	Email       string  `json:"email" gorm:"type:varchar(150);default:'';not null;comment:邮箱;index"`
	Phone       string  `json:"phone" gorm:"type:varchar(150);default:'';not null;comment:手机号;index"`
	Password    string  `json:"password" gorm:"type:varchar(150);not null;comment:密码"`
	AppKey      string  `json:"app_key" gorm:"type:varchar(150);default '';comment:应用密钥"`
	Salt        string  `json:"salt" gorm:"type:varchar(150);not null;comment:密码盐"`
	AgentId     int64   `json:"agent_ip" gorm:"type:bigint;default 0;comment:代理商id"`
	LoginIp     string  `json:"login_ip" gorm:"type:varchar(50);default '';comment:上次登录IP"`
	LoginTime   int64   `json:"login_time" gorm:"type:bigint unsigned;default:0;not null;comment:上次登录时间"`
	Status      int8    `json:"status" gorm:"type:tinyint unsigned;default:0;not null;comment:账号状态"`
	Description string  `json:"description" gorm:"type:varchar(150);default:''; comment:描述信息"`
	CreateTime  int64   `json:"create_time" gorm:"type:bigint unsigned;default:0;not null;comment:注册时间"`
	UpdateTime  int64   `json:"update_time" gorm:"type:bigint unsigned;default:0;not null;comment:会员信息上次更新时间"`
}

var UserField = []string{
	"name",
	"email",
	"phone",
	"password",
	"salt",
	"status",
	"level",
	"description",
	"create_time",
	"update_time",
	"stop_time",
}

type UserInfo struct {
	Id            int64  `json:"user_id"`
	Name          string `json:"user_name"`
	AgentName     string `json:"agent_name"`
	Email         string `json:"user_email"`
	Phone         string `json:"user_phone"`
	TotalTraffic  int64  `json:"total_traffic"`
	UsedTraffic   int64  `json:"used_traffic"`
	EnableTraffic int64  `json:"enable_traffic"`
	Status        int8   `json:"user_status"`
	Description   string `json:"description"`
	CreateTime    int64  `json:"create_time"`
	UpdateTime    int64  `json:"update_time"`
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (User) TableName() string {
	return "ip_user"
}

/*
根据账户查询用户 // 通过邮箱、手机号查询用户
*/
func (model User) FindByAccount() (user *User, err error) {
	if model.Phone == "" {
		err = DB.Where("email = ?", model.Email).Find(&user).Limit(1).Error
	} else if model.Email == "" {
		err = DB.Where("phone = ?", model.Phone).Find(&user).Limit(1).Error
	} else if model.Email != "" && model.Phone != "" {
		err = DB.Where("email = ? or  phone = ?", model.Email, model.Phone).Find(&user).Limit(1).Error
	}
	return
}

/*
根据ID查询用户
*/
func (model User) FindById() (user *User, err error) {
	err = DB.Where("id = ? ", model.Id).Find(&user).Error
	return
}

// 根据id查询用户
func (model User) SelectById() (user User, err error) {
	err = DB.Where("id = ? ", model.Id).First(&user).Error
	return
}

/*
创建用户
*/
func (model User) Create() error {
	return DB.Create(&model).Error
}

/*
更新用户状态信息 -1 删除(注销账户) 0 禁用 1 正常
*/
func (model User) UpdateStatus() error {
	return DB.Model(&model).Update("status", model.Status).Error
}

/*
重置密码
*/
func (model User) UpdatePassword() error {
	return DB.Model(&model).Update("password", model.Password).Update("salt", model.Salt).Error
}

func (model User) UpdateFlows(updates []Account) error {
	// 开始事务
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 遍历更新 map
	for _, account := range updates {
		if err := tx.Model(&account).
			UpdateColumn("used_traffic", gorm.Expr("used_traffic + ?", account.UsedTraffic)).Error; err != nil {
			// 如果有错误，回滚事务
			tx.Rollback()
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error
}

// 更新用户登录信息
func (model User) UpdateLoginInfo() error {
	return DB.Model(&model).Select("login_ip", "login_time").Updates(model).Error
}

// 用户列表查询
func (model User) SelectUserList(page, size int, agentId int64, username string) (userList []UserInfo, total int64, err error) {
	tx := DB.Table("ip_user AS iu")
	if username != "" {
		tx.Where("iu.name LIKE ?", "%"+username+"%")
	}
	err = tx.Select(`
		iu.id,
		iu.name,
		ia.name AS agent_name,
		iu.status,
		SUM(ifr.purchased_flow) AS total_traffic,
		SUM(ifr.used_flow) AS used_traffic,
		SUM(ifr.purchased_flow - ifr.used_flow) AS enable_traffic,
		iu.description,
		iu.create_time,
		iu.update_time
		`).
		Joins("LEFT JOIN ip_agent AS ia ON iu.agent_id = ia.id").
		Joins("LEFT JOIN ip_flow_record AS ifr ON iu.id = ifr.user_id").
		Where("ia.id = ?", agentId).
		Group("iu.id").
		Order("iu.create_time DESC").
		Count(&total).
		Order("iu.create_time DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&userList).Error
	return
}

// 判断用户名是否存在
func (model User) CheckUsernameExist() (count int64, err error) {
	err = DB.Model(&User{}).Where("name = ?", model.Name).Count(&count).Limit(1).Error
	return
}

// 更新用户信息
func (model User) UpdateUserInfo() error {
	return DB.Model(&model).Select("password", "status", "salt", "description", "update_time").Updates(model).Error
}
