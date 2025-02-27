package models

// ProxyServer 代理服务模型表结构体
type ProxyServer struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:varchar(100);not null" json:"name"`           // 代理服务名称
	Address       string `gorm:"type:varchar(255);not null" json:"address"`        // 代理服务地址
	Port          int    `gorm:"type:int;not null" json:"port"`                    // 代理服务端口
	ProxyProtocol string `gorm:"type:VARCHAR(100);not null" json:"proxy_protocol"` // 代理服务协议
	Status        int8   `gorm:"type:tinyint;default:1" json:"status"`             // 状态 (1: 启用, 0: 禁用)
	CreatedTime   int64  `gorm:"type:bigint" json:"created_time"`                  // 创建时间
	UpdatedTime   int64  `gorm:"autoUpdateTime" json:"updated_time"`               // 更新时间
	Description   string `gorm:"type:varchar(255)" json:"description"`             // 描述
}

func (ProxyServer) TableName() string {
	return "ip_proxy_server"
}

// 查询地址并去重
func (model ProxyServer) GetProxyServerList() (proxyServerList []ProxyServer, err error) {
	err = DB.Model(&model).Distinct("address").Find(&proxyServerList).Error
	return
}

// 根据代理服务协议查询对应的端口号
func (model ProxyServer) GetProxyServerPortByProtocol(proxyProtocol string) (port int, err error) {
	err = DB.Model(&model).Where("proxy_protocol = ?", proxyProtocol).Select("port").First(&port).Error
	return
}
