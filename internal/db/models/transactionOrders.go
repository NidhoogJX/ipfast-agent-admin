package models

import (
	"time"

	"gorm.io/gorm"
)

// 订单记录表
type TransactionOrders struct {
	UserID      int64       `json:"user_id" gorm:"type:bigint;comment:用户ID;not null;index"`
	Oid         string      `json:"oid" gorm:"primary_key;type:varchar(300);comment:订单号;not null;index"`
	Platform    string      `json:"platform" gorm:"type:varchar(50);comment:交易平台;not null;default:'';"`
	Tid         string      `json:"tid" gorm:"type:varchar(300);comment:交易平台订单号;not null;default:'';index"`
	Status      int8        `json:"status" gorm:"type:tinyint;comment:订单状态:0已创建,1未付款,2已付款,3交易完成,4交易超时,5交易失败;default:0;not null"`
	Desc        string      `json:"desc" gorm:"type:varchar(500);comment:订单描述;not null"`
	Amount      float64     `json:"amount" gorm:"type:decimal(10,2);comment:订单金额;default:0"`
	Currency    string      `json:"currency" gorm:"type:varchar(50);comment:订单货币;not null"`
	OrderType   int8        `json:"order_type" gorm:"type:tinyint;comment:订单类型;not null;default:1;comment:1动态住宅IP,2静态住宅IP,3数据中心IP,10静态订单续费,11数据中心订单续费"`
	PaymentLink string      `json:"payment_link" gorm:"type:text;comment:支付链接;"`
	CreatedTime int64       `json:"created_time" gorm:"type:bigint;comment:创建时间;not null;default:0"`
	UpdatedTime int64       `json:"updated_time" gorm:"type:bigint;comment:更新时间;not null;default:0"`
	Items       []OrderItem `gorm:"foreignkey:OrderID;references:Oid"`
}

type OrderItem struct {
	ID             int64   `gorm:"primary_key;AUTO_INCREMENT"`
	AreaId         int64   `gorm:"not null;type:bigint;default:0"`
	OrderID        string  `gorm:"not null;index;type:bigint"`
	CommodityID    int64   `gorm:"not null;type:bigint"`
	DurationTypeId int64   `gotm:"not null;type:bigint"` // 时长类型ID
	CommodityName  string  `gorm:"not null;type:varchar(200)"`
	Quantity       int64   `gorm:"not null;type:bigint;数量(流量)"`
	Unit           int8    `gorm:"not null;type:tinyint;default:0;comment:单位(1-GB;2-TB)"`
	Amount         float64 `gorm:"not null;type:decimal(10,2)"`
	Type           int8    `gorm:"not null;type:tinyint;comment:商品类型(1动态;2静态;3数据中心)"`
	Desc           string  `gorm:"not null;type:varchar(500)"`
	Ext1           int64   `gorm:"not null;type:bigint;default:0"`
}

type OrderInfo struct {
	OrderId     string  `json:"order_id"`
	Email       string  `json:"email"`
	Platform    string  `json:"platform"`
	Tid         string  `json:"tid"`
	OrderType   int8    `json:"order_type"`
	Countries   string  `json:"countries"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Quantity    int64   `json:"quantity"`
	Unit        int8    `json:"unit"`
	PaymentLink string  `json:"payment_link"`
	Status      int8    `json:"status"`
	Desc        string  `json:"desc"`
	CreatedTime int64   `json:"created_time"`
	UpdatedTime int64   `json:"updated_time"`
}

func (OrderItem) TableName() string {
	return "ip_order_items"
}

// 根据订单号查询订单所有订单项
func (model OrderItem) SelectOrderItemsByOid(oid string) (items []OrderItem, err error) {
	err = DB.Model(model).Where("order_id = ?", oid).Find(&items).Error
	return
}

func (TransactionOrders) TableName() string {
	return "ip_transaction_orders"
}

func (model TransactionOrders) Insert() error {
	return DB.Create(&model).Error
}

func (model TransactionOrders) GetByOid() (order TransactionOrders, err error) {
	err = DB.Preload("Items").Where("oid = ?", model.Oid).First(&order).Error
	return
}

func (model TransactionOrders) GetOrderByUidAndOid(page, pageSize int, startDate, endDate int64) (orders []TransactionOrders, total int64, err error) {
	tx := DB.Preload("Items")
	if model.OrderType != 0 {
		tx = tx.Where("order_type = ?", model.OrderType)
	}
	if model.Oid != "" {
		tx = tx.Where("oid = ?", model.Oid)
	}
	if startDate != 0 {
		tx = tx.Where("created_time >= ? ", startDate)
	}
	if endDate != 0 {
		tx = tx.Where("created_time <= ?", endDate)
	}

	// 计算总记录数
	err = tx.Model(&TransactionOrders{}).Where("user_id = ?", model.UserID).Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询
	err = tx.Limit(pageSize).
		Offset((page - 1) * pageSize).Order("created_time desc").
		Find(&orders).Error
	return
}
func (model TransactionOrders) Update() error {
	model.UpdatedTime = time.Now().Unix()
	return DB.Instance.Updates(&model).Error
}

func (model TransactionOrders) UpdateByTransaction(tx *gorm.DB) error {
	return tx.Where("oid = ?", model.Oid).Updates(&model).Error
}

// 查询未支付的订单
func (model TransactionOrders) SelectUnpaidStaticOrders() (orders []TransactionOrders, err error) {
	err = DB.Where("status = 1 and order_type = 2").Find(&orders).Error
	return
}

// 更新订单状态
func (model TransactionOrders) UpdateOrderStatus(order TransactionOrders) error {
	return DB.Table(model.TableName()).Where("oid = ?", order.Oid).Select("status", "updated_time").Updates(order).Error
}

// 根据ipid查询订单
func SelectOrderByIpId(ipId int64) (orders TransactionOrders, err error) {
	err = DB.Where("ip_id = ?", ipId).First(&orders).Error
	return
}

// 订单列表查询
func (model TransactionOrders) SelectOrderList(page, size int, email string) (orders []OrderInfo, total int64, err error) {
	tx := DB.Table("ip_transaction_orders as ito")
	if email != "" {
		tx = tx.Where("email like ?", "%"+email+"%")
	}
	err = tx.Select(`
			ito.oid,
			iu.email,
			ito.platform,
			ito.tid,
			ito.DESC,
			ito.amount,
			ito.currency,
			ito.order_type,
			ito.payment_link,
			ito.STATUS,
			ito.created_time,
			ito.updated_time,
			COALESCE ( ioi.quantity, 0 ) AS quantity,
			ioi.unit,
			GROUP_CONCAT( DISTINCT irc.country_name SEPARATOR ' | ' ) AS countries
		`).
		Joins("LEFT JOIN ip_user AS iu ON ito.user_id = iu.id").
		Joins("LEFT JOIN ip_order_items AS ioi ON ito.oid = ioi.order_id").
		Joins("LEFT JOIN ip_region_country AS irc ON ioi.area_id = irc.country_id AND ioi.area_id != 0").
		Group(`
			ioi.quantity,
			ioi.unit,
			ito.oid
		`).
		Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Scan(&orders).Error
	return
}
