package models

import "gorm.io/gorm"

/*
流量账号管理
*/
type IpWhiteList struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:白名单ID"`
	UserId      int64  `json:"user_id" gorm:"type:bigint unsigned;default:0;not null;comment:用户ID"`
	Ip          string `json:"ip" gorm:"type:varchar(150);not null;comment:IP白名单地址"`
	Desc        string `json:"desc" gorm:"type:varchar(150);default:'';not null;comment:白名单描述"`
	CreatedTime int64  `json:"created_time" gorm:"type:bigint unsigned;default:0;not null;comment:创建时间"`
	UpdatedTime int64  `json:"updated_time" gorm:"type:bigint unsigned;default:0;not null;comment:更新时间"`
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (IpWhiteList) TableName() string {
	return "ip_white_list"
}

func (model IpWhiteList) CreateDatas(datas []IpWhiteList) error {
	return DB.Instance.Save(&datas).Error
}

// 查询当前ip是否存在
func (model IpWhiteList) Exist(ip string, uid int64) bool {
	var count int64
	DB.Model(&IpWhiteList{}).Where("ip = ?", ip).Where("user_id = ?", uid).Count(&count)
	return count > 0
}

// 批量更新记录
func (model IpWhiteList) Updates(datas []IpWhiteList) error {
	return DB.Instance.Transaction(func(tx *gorm.DB) error {
		for _, data := range datas {
			if err := tx.Model(&IpWhiteList{}).Where("id = ?", data.Id).Where("user_id = ?", model.UserId).Select("ip", "desc", "updated_time").Updates(data).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Paginate 分页查询 根据 IP 查询
func (m IpWhiteList) Paginate(page, pageSize int, ip string, uid int64) (results []IpWhiteList, total int64, err error) {
	query := DB.Model(&IpWhiteList{})
	if ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}
	err = query.Where("user_id = ?", uid).Count(&total).Error
	if err != nil {
		return
	}
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&results).Error
	return
}

func (model IpWhiteList) Delete(ids []int64, uid int64) error {
	return DB.Where("id in (?)", ids).Where("user_id = ?", uid).Delete(&model).Error
}
