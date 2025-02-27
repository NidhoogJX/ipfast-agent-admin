package models

/*
流量账号管理
*/
type IPIPGoAccountStaticIp struct {
	OrderID string `json:"orderId" gorm:"type:varchar(150);not null;comment:订单ID;index"`          //订单ID
	Uid     int64  `json:"uid" gorm:"type:bigint unsigned;not null;comment:用户ID;index"`           //用户ID
	IpPort  string `json:"ipPort" gorm:"varchar(150);default:'';not null;comment:IPIPGO二级客户对应账号"` //IP地址端口
	//AccountPassword string `json:"accountPassword" gorm:"varchar(150);default:'';not null;comment:IPIPGO静态IP账密"` //账号|密码
	Account     string `json:"account" gorm:"varchar(150);default:'';not null;comment:IPIPGO静态IP名称"`     //账号名称
	Password    string `json:"password" gorm:"varchar(150);default:'';not null;comment:IPIPGO静态IP密码"`    //密码
	AddressStr  string `json:"addressStr" gorm:"varchar(150);default:'';not null;comment:国家地区名称"`        //国家地区名称 //状态 0失效1过期2正常
	Status      int8   `json:"status" gorm:"tinyint(1);default:0;not null;comment:状态 0失效1过期2正常"`         //状态 0失效1过期2正常
	EndTime     string `json:"endTime" gorm:"varchar(150);default:'';not null;comment:到期时间"`             //到期时间
	MealId      int    `json:"mealId" gorm:"int(11);default:0;not null;comment:套餐ID"`                    //套餐ID
	CmiId       int64  `json:"cmiId" gorm:"bigint(20);default:0;not null;comment:客户套餐ID"`                //客户套餐ID
	CreatedTime int64  `json:"created_time" gorm:"type:bigint unsigned;not null;default:0;comment:创建时间"` //创建时间
	UpdatedTime int64  `json:"updated_time" gorm:"type:bigint unsigned;not null;default:0;comment:更新时间"` //更新时间
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (IPIPGoAccountStaticIp) TableName() string {
	return "ip_ipipgo_account_static_ip"
}

// GetAuthInfo 获取IPIPGO二级客户静态IP信息 支持分页
func (m IPIPGoAccountStaticIp) GetAccount() (data *IPIPGoAccount, err error) {
	data = &IPIPGoAccount{}
	err = DB.Where("uid = ?", m.Uid).Find(data).Error
	return
}

// Insert 插入IPIPGO二级客户静态IP信息
func (m IPIPGoAccountStaticIp) Insert() error {
	return DB.Create(m).Error
}

func (m IPIPGoAccountStaticIp) BatchInsert(data []IPIPGoAccountStaticIp) error {
	return DB.Create(&data).Error
}

// Paginate 分页查询，支持根据 IP 字段、状态和地区进行条件查询
func (m IPIPGoAccountStaticIp) Paginate(page, pageSize int, ip, addressStr string, status int8, uid int64) (results []IPIPGoAccountStaticIp, total int64, err error) {
	query := DB.Model(&IPIPGoAccountStaticIp{})

	if ip != "" {
		query = query.Where("ip_port LIKE ?", "%"+ip+"%")
	}
	if status > -1 && status < 3 {
		query = query.Where("status = ?", status)
	}
	if addressStr != "" {
		query = query.Where("address_str LIKE ?", "%"+addressStr+"%")
	}

	err = query.Where("uid = ?", uid).Count(&total).Error
	if err != nil {
		return
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&results).Error
	return
}
