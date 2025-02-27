package models

type PayPlatform struct {
	PayPlatformID     int64  `json:"pay_platform_id" gorm:"type:bigint; primary_key; AUTO_INCREMENT; not null; comment:付费平台id"`
	PayPlatformName   string `json:"pay_platform_name" gorm:"type:varchar(100); not null; comment:付费平台名称"`
	PayIdentification string `json:"pay_identification" gorm:"type:varchar(30); not null; comment:付费平台标识(判断执行哪个支付SDK)"`
	PayMethod         string `json:"pay_method" gorm:"type:varchar(30); not null; comment:付费方式(信用卡,支付宝,微信)"`
	Currency          string `json:"currency" gorm:"type:varchar(30); not null; comment:货币"`
	AmountRange       string `json:"amount_range" gorm:"type:varchar(100); not null; comment:金额范围"`
	Status            int    `json:"status" gorm:"type:int; not null;  default:1; comment:支付平台状态(0:停用,1:启用)"`
	CreatedTime       int64  `json:"created_time" gorm:"type:bigint; not null; default:0; comment:创建时间"`
	UpdatedTime       int64  `json:"updated_time" gorm:"type:bigint; default:0; comment:更新时间"`
}

type PayPlatformList struct {
	PayPlatformID   int64  `json:"pay_platform_id"`
	PayPlatformName string `json:"pay_platform_name"`
}

// TableName 设置表名
func (PayPlatform) TableName() string {
	return "ip_pay_platform"
}

// 根据id查询支付平台信息
func (model PayPlatform) SelectPayPlatformByID() (payPlatform PayPlatform, err error) {
	err = DB.Table(model.TableName()).Where("pay_platform_id = ? AND status = ?", model.PayPlatformID, model.Status).Find(&payPlatform).Error
	if err != nil {
		return payPlatform, err
	}
	return payPlatform, err
}

// 查询正在启用的支付平台id，名称
func (model PayPlatform) SelectEnabledPayPlatform() (payPlatformList []PayPlatformList, err error) {
	err = DB.Table(model.TableName()).Select("pay_platform_id", "pay_platform_name").Where("status = ?", 1).First(&payPlatformList).Error
	if err != nil {
		return nil, err
	}
	return payPlatformList, err
}
