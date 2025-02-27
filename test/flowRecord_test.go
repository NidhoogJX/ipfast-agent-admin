package test

import (
	// _ "ipfast_server/internal/config/api"
	"ipfast_server/internal/db/models"
	"ipfast_server/internal/services"
	"log"
	"testing"

	"gorm.io/gorm"
)

func TestFlowRecord(t *testing.T) {

	order := &models.TransactionOrders{
		UserID:      1,
		Oid:         "oid-00001",
		Platform:    "platform-test",
		Tid:         "tid-00001",
		Status:      1,
		Desc:        "desc-test",
		Amount:      100,
		Currency:    "CNY",
		OrderType:   1,
		CreatedTime: 1732521354,
		UpdatedTime: 1732521354,
		Items: []models.OrderItem{
			{
				ID:             00001,
				AreaId:         1,
				OrderID:        "oid-00001",
				CommodityID:    1,
				CommodityName:  "commodity-test",
				DurationTypeId: 1,
				Quantity:       1,
				Amount:         100,
				Type:           1,
				Desc:           "item1-desc-test",
				Ext1:           1,
			},
			{
				ID:             00002,
				AreaId:         1,
				OrderID:        "oid-00001",
				CommodityID:    1,
				CommodityName:  "commodity-test",
				DurationTypeId: 1,
				Quantity:       1,
				Amount:         100,
				Type:           1,
				Desc:           "item2-desc-test",
				Ext1:           1,
			},
		},
	}

	err := models.DB.Instance.Transaction(func(tx *gorm.DB) error {
		services.UpdateFlowRecord(order.UserID, order.Items[0].DurationTypeId, order.Items, tx)
		return nil
	})
	if err != nil {
		log.Printf("失败:%v", err)
	}

}
