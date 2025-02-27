package models

/*
实时节点
*/
type Xfj struct {
	Id       int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:账号ID"`
	Ip       uint32 `json:"ip" gorm:"type:int unsigned;default:0;not null;comment:IP地址"`
	ServerIp uint32 `json:"server_ip" gorm:"type:int unsigned;default:0;not null;comment:服务器IP"`
	Uptime   int64  `json:"uptime" gorm:"type:bigint unsigned;default:0;not null;comment:账号创建时间"`
	Country  int32  `json:"country" gorm:"type:int unsigned;default:0;not null;comment:国家"`
	Province int32  `json:"province" gorm:"type:int unsigned;default:0;not null;comment:省份"`
	City     int32  `json:"city" gorm:"type:int unsigned;default:0;not null;comment:城市"`
}

/*
返回数据库表名

	struct:
		Device 客户端信息
	return:
		string: 表名
*/
func (Xfj) TableName() string {
	return "ip_xfj"
}

/* 获取可用的国家列表 */
func (model *Xfj) SelectEnableCountry() ([]int32, error) {
	var xfjIDs []int32
	err := DB.Table(model.TableName()).Select("DISTINCT country").Find(&xfjIDs).Error
	return xfjIDs, err
}

/* 获取可用的省份列表 */
func (model *Xfj) SelectEnableProvince() ([]int32, error) {
	var xfjIDs []int32
	err := DB.Table(model.TableName()).Select("DISTINCT province").Find(&xfjIDs).Error
	return xfjIDs, err
}

/* 获取可用的城市列表 */
func (model *Xfj) SelectEnableCity() ([]int32, error) {
	var xfjIDs []int32
	err := DB.Table(model.TableName()).Select("DISTINCT city").Find(&xfjIDs).Error
	return xfjIDs, err
}
