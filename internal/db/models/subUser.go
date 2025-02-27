package models

// 子用户关系表
type SubUser struct {
	Id              int64   `json:"id" gorm:"type:bigint; primary_key; AUTO_INCREMENT; not null; comment:子用户ID"`
	ParentUserID    int64   `json:"parent_user_id" gorm:"type:bigint; not null; comment:父用户ID"`
	Username        string  `json:"username" gorm:"type:varchar(150); not null; comment:子用户名"`
	Password        string  `json:"password" gorm:"type:varchar(150); not null; comment:子用户密码"`
	MaxCapacity     float64 `json:"max_capacity" gorm:"type:decimal(10,2); comment:使用流量上限(GB)"`
	MaxCapacityByte int64   `json:"max_capacity_byte" gorm:"type:bigint; comment:使用流量上限(B)"`
	UsedCapacity    int64   `json:"used_capacity" gorm:"type:bigint; ;default:0; comment:用户已使用流量(B)"`
	Remarks         string  `json:"remarks" gorm:"type:varchar(255); comment:子用户备注"`
	Status          int8    `json:"status" gorm:"type:tinyint(1); not null; default:1; comment:状态: 1-正常 0-禁用"`
	MaxStatus       int8    `json:"max_status" gorm:"type:tinyint(1); not null; default:0; comment:是否开启用户流量上限: 1-正常 0-禁用"`
	CreatedTime     int64   `json:"created_time" gorm:"type:bigint; not null; default:0; comment:创建时间"`
	UpdatedTime     int64   `json:"updated_time" gorm:"type:bigint; default:0; comment:更新时间"`
}

var SubUserUpdateFieid = []string{
	// "username",
	"password",
	"max_capacity",
	"max_capacity_byte",
	"used_capacity",
	"remarks",
	"status",
	"max_status",
	"updated_time",
}

type SubUserInfo struct {
	Id              int64   `json:"id"`
	ParentUserID    int64   `json:"parent_user_id"`
	ParentName      string  `json:"parent_name"`
	Username        string  `json:"username"`
	Password        string  `json:"password"`
	MaxCapacity     float64 `json:"max_capacity"`
	MaxCapacityByte int64   `json:"max_capacity_byte"`
	UsedCapacity    int64   `json:"used_capacity"`
	Remarks         string  `json:"remarks"`
	Status          int8    `json:"status"`
	MaxStatus       int8    `json:"max_status"`
	CreatedTime     int64   `json:"created_time"`
	UpdatedTime     int64   `json:"updated_time"`
}

// TableName 设置表名
func (SubUser) TableName() string {
	return "ip_sub_user"
}

// 新增子用户
func (model SubUser) Create() (subUser SubUser, err error) {
	model.MaxCapacityByte = (int64)(model.MaxCapacity * 1024 * 1024 * 1024) // 转换为B
	subUser, err = model, DB.Create(&model).Error
	return subUser, err
}

// 更新子用户信息
func (model *SubUser) Update() error {
	model.MaxCapacityByte = (int64)(model.MaxCapacity * 1024 * 1024 * 1024) // 转换为B
	return DB.Model(&model).Select(SubUserUpdateFieid).Updates(model).Error
}

// 删除子用户
func (model *SubUser) Delete() error {
	return DB.Model(&model).Where("id = ?", model.Id).Delete(model).Error
}

// 批量删除子用户
func (model *SubUser) DeleteByIDs(ids []int64, uid int64) error {
	return DB.Where("id in (?) and parent_user_id = ?", ids, uid).Delete(model).Error
}

// 分页查询子用户列表
func (model SubUser) GetSubUsers(page, size int, uid int64, subUserName string) (subUser []SubUser, total int64, err error) {
	tx := DB.Model(&model)
	if subUserName != "" {
		tx.Where("username like ?", "%"+subUserName+"%")
	}
	// 计算总记录数
	err = tx.Model(&model).Where("parent_user_id = ?", uid).Count(&total).Error
	if err != nil {
		return
	}
	// 分页查询
	err = tx.Limit(size).
		Offset((page - 1) * size).
		Order("created_time desc").
		Find(&subUser).Error
	return subUser, total, err
}

// 根据用户名模糊查询子用户信息
func (model SubUser) GetSubUserByUsername(username string, uid int64) (subUser SubUser, er error) {
	err := DB.Where("username like ? and parent_user_id = ?", username, uid).Find(&subUser).Error
	return subUser, err
}

// 根据状态查询子用户信息
func (model SubUser) GetSubUserByStatus() (subUser []SubUser, er error) {
	err := DB.Where("status = ? AND   parent_user_id = ?", model.Status, model.ParentUserID).Find(&subUser).Error
	return subUser, err
}

// 根据username查询子用户数量(判断子用户是否存在)
func (model SubUser) IsUsernameExist(username string, uid int64) (Count int64, err error) {
	err = DB.Model(&model).Where("username = ? and parent_user_id = ?", username, uid).Count(&Count).Error
	return Count, err
}

// 根据username查询子用户数量(排除当前子账号)
func (model SubUser) IsUsernameExistExcludeCurrent(username string, uid, subUserId int64) (Count int64, err error) {
	err = DB.Model(&model).Where("username = ? and parent_user_id = ? and id != ?", username, uid, subUserId).Count(&Count).Error
	return Count, err
}

// 根据id和parent_user_id查询子用户信息
func (model SubUser) GetSubUserById() (subUser SubUser, err error) {
	err = DB.Where("id = ? AND parent_user_id = ?", model.Id, model.ParentUserID).Find(&subUser).Error
	return subUser, err
}

// 根据子用户id查询子用户信息
func (model SubUser) GetSubUserBySubUserId() (subUser SubUser, err error) {
	err = DB.Where("id = ?", model.Id).Find(&subUser).Error
	return subUser, err
}

// 根据父用户id查询所有子用户信息
func (model SubUser) GetSubUsersByParentUserId() (subUsers []SubUser, err error) {
	err = DB.Where("parent_user_id = ?", model.ParentUserID).Find(&subUsers).Error
	return subUsers, err
}

// 查询代理商下用户的子账户
func (model SubUser) SelectSubUserListByAgentId(userId int64, page, size int, subuserName string, status int8) (subUsers []SubUserInfo, total int64, err error) {
	tx := DB.Table("ip_sub_user AS isu")
	if subuserName != "" {
		tx.Where("isu.username like ?", "%"+subuserName+"%")
	}
	if status != 2 {
		tx.Where("isu.status = ?", status)
	}
	err = tx.Select(`
			isu.*,
			iu.name AS parent_name
		`).
		Joins("LEFT JOIN ip_user AS iu ON isu.parent_user_id = iu.id").
		Where("parent_user_id = ?", userId).
		Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&subUsers).Error
	return
}

// 根据状态查询子账户列表
func (model SubUser) SelectSubUserListByStatus(userId int64, status int8) (subUsers []SubUserInfo, err error) {
	tx := DB.Table("ip_sub_user AS isu")
	if status != 2 {
		tx.Where("isu.status = ?", status)
	}
	err = tx.Select(`
			isu.*,
			iu.name AS parent_name
		`).
		Joins("LEFT JOIN ip_user AS iu ON isu.parent_user_id = iu.id").
		Where("isu.parent_user_id = ?", userId).
		Scan(&subUsers).Error
	return
}
