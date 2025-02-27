package models

import (
	"time"

	"gorm.io/gorm"
)

/*
购买IP记录
*/
type IpRecord struct {
	Id          int64  `json:"id" gorm:"type:bigint; primary_key; AUTO_INCREMENT; not null; comment:ip记录ID"`
	UserID      int64  `json:"user_id" gorm:"type:bigint; not null; comment:用户ID; index"`
	Ip          string `json:"ip" gorm:"type:varchar(100); not null; comment:IP地址; index"`
	Port        uint32 `json:"port" gorm:"type:int unsigned; not null; comment:端口"`
	OrderId     string `json:"order_id" gorm:"type:varchar(200); not null; comment:订单ID(管理员添加的ip,order_id为0)"`
	Username    string `json:"username" gorm:"type:varchar(150); not null; comment:用户名"`
	Password    string `json:"password" gorm:"type:varchar(200); not null; comment:密码"`
	Deadline    int64  `json:"deadline" gorm:"type:bigint; not null; comment:过期时间"`
	CountryId   int64  `json:"country_id" gorm:"type:bigint; comment:国家ID; index"`
	Status      int8   `json:"status" gorm:"type:tinyint; comment:状态(0失效,1过期,2正常,3分配中);"`
	Type        int8   `json:"type" gorm:"type:tinyint; comment:IP类型(2静态IP数据,3数据中心)"`
	CreatedTime int64  `json:"created_time" gorm:"type:bigint; not null; default:0; comment:创建时间"`
	UpdatedTime int64  `json:"updated_time" gorm:"type:bigint; default:0; comment:更新时间"`
}

type IpRecordByCountry struct {
	Id          int64  `json:"id" gorm:"type:bigint; primary_key; AUTO_INCREMENT; not null; comment:ip记录ID"`
	UserID      int64  `json:"user_id" gorm:"type:bigint; not null; comment:用户ID; index"`
	Ip          string `json:"ip" gorm:"type:varchar(100); not null; comment:IP地址; index"`
	Port        uint32 `json:"port" gorm:"type:int unsigned; not null; comment:端口"`
	OrderId     string `json:"order_id" gorm:"type:varchar(200); not null; comment:订单号"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	Status      int8   `json:"status" gorm:"type:tinyint; comment:状态(0失效,1过期,2正常,3分配中);"`
	Username    string `json:"username" gorm:"type:varchar(150); not null; comment:用户名"`
	Password    string `json:"password" gorm:"type:varchar(200); not null; comment:密码"`
	Deadline    int64  `json:"deadline" gorm:"type:bigint; not null; comment:过期时间"`
	CountryId   int64  `json:"country_id" gorm:"type:bigint; comment:国家ID; index"`
	CreatedTime int64  `json:"created_time" gorm:"type:bigint; not null; default:0; comment:创建时间"`
	UpdatedTime int64  `json:"updated_time" gorm:"type:bigint; default:0; comment:更新时间"`
}

var fieldList = []string{
	"ip_ip_record.id",
	"ip_ip_record.user_id",
	"ip_ip_record.ip",
	"ip_ip_record.port",
	"ip_ip_record.order_id",
	"ip_ip_record.username",
	"ip_ip_record.password",
	"ip_ip_record.status",
	"ip_ip_record.country_id",
	"ip_ip_record.deadline",
	"ip_ip_record.created_time",
	"ip_ip_record.updated_time",
}

var IpRecordField = []string{
	"deadline",
	"status",
	"ip",
	"port",
	"username",
	"password",
	"updated_time",
}

type IpRecordInfo struct {
	Id          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Ip          string `json:"ip"`
	Port        uint32 `json:"port"`
	OrderId     string `json:"order_id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Deadline    int64  `json:"deadline"`
	CountryId   int64  `json:"country_id"`
	RegionName  string `json:"region_name"`
	CountryName string `json:"country_name"`
	Status      int8   `json:"status"`
	Type        int8   `json:"type"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}

// 返回表名
func (IpRecord) TableName() string {
	return "ip_ip_record"
}

// Paginate 分页查询，支持根据 IP 字段、状态和地区进行条件查询
func (model IpRecord) Paginate(page, pageSize int, ip, port string, status int8, country_id, uid int64, iptype int8) (results []IpRecordByCountry, total int64, err error) {
	query := DB.Model(&IpRecord{})
	if ip != "" {
		query = query.Where("port LIKE ?", "%"+ip+"%")
	}
	if port != "" {
		query = query.Where("port LIKE ?", "%"+port+"%")
	}
	if status > -1 && status < 3 {
		query = query.Where("status = ?", status)
	}
	if country_id > 0 {
		query = query.Where("country_id = ?", country_id)
	}
	err = query.Where("user_id = ?", uid).Where("type = ?", iptype).Count(&total).Error
	if err != nil {
		return
	}
	offset := (page - 1) * pageSize
	err = query.
		Select(fieldList, "ip_traffic_country.Name AS country_name", "ip_traffic_region.name AS region_name").
		Joins("JOIN ip_traffic_country ON ip_traffic_country.id = ip_ip_record.country_id").
		Joins("JOIN ip_traffic_region ON ip_traffic_region.id = ip_traffic_country.region_id").
		Offset(offset).
		Limit(pageSize).
		Order("ip_ip_record.created_time DESC").
		Find(&results).Error
	return
}

// 插入数据
func (model IpRecord) Insert() error {
	return DB.Create(&model).Error
}

// 批量插入数据
func (model IpRecord) InsertBatch(models []IpRecord, tx *gorm.DB) error {
	return tx.Create(&models).Error
}

// 更新数据
func (model IpRecord) Update() error {
	return DB.Model(&model).Where("id = ?", model.Id).Select(IpRecordField).Updates(model).Error
}

// 删除数据
func (model IpRecord) Delete() error {
	return DB.Model(&model).Where("id = ?", model.Id).Delete(model).Error
}

func (model IpRecord) GetCountByType() (TotalCount int64, err error) {
	err = DB.Table(model.TableName()).Select("count(id)").Where("ip_ip_record.type = ? AND ip_ip_record.user_id = ? ", model.Type, model.UserID).Count(&TotalCount).Error
	return
}

// 根据用户id查询未过期的静态IP/数据中心IP
func (model IpRecord) GetByUserId() (ipRecords []IpRecord, err error) {
	err = DB.Model(&IpRecord{}).Where("user_id = ? and type = ? and deadline > ? ", model.UserID, model.Type, time.Now().Unix()).
		Find(&ipRecords).Error
	return ipRecords, err
}

type WaitIpCount struct {
	Type      int8  `json:"type"`       // 类型	2静态IP数据,3数据中心
	WaitCount int64 `json:"wait_count"` // 等待分配的IP数量
}

// 查询待分配ip的用户数
func (model IpRecord) SelectWaitIpUsers() (waitIpCounts []WaitIpCount, err error) {
	err = DB.Table("ip_ip_record").Select("type, COUNT(DISTINCT user_id) AS wait_count").
		Where("status = 3").
		Group("type").
		Scan(&waitIpCounts).Error
	return
}

// 根据用户id查询ip列表
func (model IpRecord) SelectIpList(page, size int, userId int64, sortStatus, ipType int8, ip string) (staticIpList []IpRecordInfo, total int64, err error) {
	tx := DB.Table("ip_ip_record AS iir")
	if ip != "" {
		tx.Where("iir.ip LIKE ?", "%"+ip+"%")
	}
	if sortStatus == 1 {
		tx.Order("iir.status DESC")
	} else if sortStatus == 2 {
		tx.Order("iir.status ASC")
	}
	err = tx.Select(`
		iir.id,
		iir.user_id,
		iir.ip,
		iir.deadline,
		itr.name AS region_name,
		irc.country_name,
		iir.port,
		iir.order_id,
		iir.username,
		iir.password,
		iir.country_id,
		iir.type,
		iir.status,
		iir.created_time,
		iir.updated_time
	`).
		Joins("LEFT JOIN ip_region_country AS irc ON iir.country_id = irc.country_id").
		Joins("LEFT JOIN ip_traffic_country AS itc ON iir.country_id = itc.id ").
		Joins("LEFT JOIN ip_traffic_region AS itr ON itc.region_id = itr.id").
		Where("iir.user_id = ? and iir.type = ?", userId, ipType).
		Count(&total).
		Scan(&staticIpList).
		Offset(size).
		Limit(page).Error
	return
}
