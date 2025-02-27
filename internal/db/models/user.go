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
	Email         string `json:"user_email"`
	Phone         string `json:"user_phone"`
	Password      string `json:"user_password"`
	TotalTraffic  int64  `json:"total_traffic"`
	UsedTraffic   int64  `json:"used_traffic"`
	EnableTraffic int64  `json:"enable_traffic"`
	LoginIp       string `json:"login_ip"`
	LoginTime     int64  `json:"login_time"`
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
func (model User) SelectUserList(agentId int64, username string, page, size int, status, totalSort, usedSort, enableSort int8) (userList []UserInfo, total int64, err error) {
	tx := DB.Table("ip_user as iu")
	if username != "" {
		tx = tx.Where("name like ?", "%"+username+"%")
	}
	if status != 2 {
		tx = tx.Where("status = ?", status)
	}
	err = tx.Select(`
		iu.id,
		iu.name,
		iu.status,
		iu.email,
		iu.phone,
		iu.password,
		iu.description,
		iu.login_ip,
		iu.login_time,
		iu.create_time,
		iu.update_time,
		IFNULL(SUM(ifr.purchased_flow),0) AS total_traffic,
		IFNULL(SUM(ifr.used_flow),0) AS used_traffic,
		IFNULL(SUM(ifr.purchased_flow-ifr.used_flow),0) AS enable_traffic
		`).
		Joins("LEFT JOIN ip_flow_record as ifr on iu.id = ifr.user_id").
		Where("iu.agent_id = ?", agentId).
		Count(&total).
		Group("iu.id").Error
	if totalSort != 0 {
		if totalSort == 1 {
			tx.Order("total_traffic DESC")
		}
		if totalSort == 2 {
			tx.Order("total_traffic ASC")
		}
	} else if usedSort != 0 {
		if usedSort == 1 {
			tx.Order("used_traffic DESC")
		}
		if usedSort == 2 {
			tx.Order("used_traffic ASC")
		}
	} else if enableSort != 0 {
		if enableSort == 1 {
			tx.Order("enable_traffic DESC")
		}
		if enableSort == 2 {
			tx.Order("enable_traffic ASC")
		}
	}
	tx.Offset((page - 1) * size).
		Limit(size).
		Scan(&userList)
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
