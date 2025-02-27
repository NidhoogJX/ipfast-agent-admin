package services

import (
	"ipfast_server/internal/db/core/kafka"
	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/api/accountFlow"
	"ipfast_server/pkg/util/log"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"google.golang.org/protobuf/proto"
)

type FlowRecord = models.FlowRecord

var UserCurrentFlow = make(map[int64]FlowRecord)
var SubUserAndParentUser = make(map[int64]int64)

const batchSize = 5000 // 批量处理的大小
const insterSize = 500 // 批量插入的大小

// 创建一个 channel 来接收消息
var msgChan = make(chan *accountFlow.AccountFlow, batchSize)

/*
开始接收埋点数据,并入库
*/
func StartReceiveFlowData() {
	log.Info("[FLOW]Start Receive Flow Data")
	ants.Submit(func() {
		for {
			message, err := kafka.Reader.ReadMessage(kafka.Topic)
			if err != nil {
				log.Error("[FLOW]Read Kafka Message Error: %v", err)
				continue
			}
			flowData := &accountFlow.AccountFlows{}
			err = proto.Unmarshal(message.Value, flowData)
			if err != nil {
				log.Error("[FLOW]Unmarshal FlowData Error: %v", err)
				continue
			}
			for _, v := range flowData.AccountFlowDatas {
				if v.UploadFlow > 0 {
					msgChan <- v
				}
			}

		}
	})
	ants.Submit(func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			batch := make([]*accountFlow.AccountFlow, 0, insterSize)
			i := 0
		outerLoop:
			for i < insterSize {
				select {
				case msg, ok := <-msgChan:
					if !ok {
						log.Error("[FLOW]msgChan 通道已关闭")
						break
					}
					i = i + 1
					batch = append(batch, msg)
				case <-ticker.C:
					break outerLoop
				}
			}
			if len(batch) > 0 {
				log.Info("读取数据%d条", len(batch))
				ants.Submit(func() {
					UpdateUserFlowData(batch)
				})
				ants.Submit(func() {
					InsertFlowInfo(batch)
				})
			}
		}
	})
}

func InsertFlowInfo(data []*accountFlow.AccountFlow) {
	var InsertBuriedData = []models.UserFlow{}
	now := time.Now().Unix()
	for _, v := range data {
		InsertBuriedData = append(InsertBuriedData, models.UserFlow{
			SUserID:     v.AccountId,
			Up:          v.UploadFlow,
			CreatedTime: now,
			UpdatedTime: now,
		})
	}
	log.Info("开始写入流量数据: %d/条", len(InsertBuriedData))
	err := models.UserFlow{}.Inserts(InsertBuriedData)
	if err != nil {
		log.Error("[FLOW]Insert Flow Datas Error: %v", err)
	}
}

func UpdateUserFlowData(data []*accountFlow.AccountFlow) {
	for _, v := range data {
		UpdateUserFlowRecord(v.AccountId, v.UploadFlow)
	}
}

func UpdateUserFlowRecord(userId int64, flow int64) {
	var data FlowRecord
	var ok bool
	var err error
	data, ok = UserCurrentFlow[userId]
	if !ok {
		data, err = GetFlowRecordByUserId(userId)
		if err != nil {
			log.Error("[FLOW]Get Flow Record Error: %v", err)
		}
		data.UsedFlow = 0
	}
	if data.UserID == 0 {
		return
	}
	data.UsedFlow += flow
	UserCurrentFlow[userId] = data
}

func GetFlowRecordByUserId(userID int64) (flowRecords FlowRecord, err error) {
	flowRecords, err = models.FlowRecord{
		UserID:     userID,
		IdentityId: 1,
	}.GetByUserSingel()
	if err != nil || flowRecords.Id == 0 {
		var parentUserId int64
		parentUserId, ok := SubUserAndParentUser[userID]
		if !ok {
			subUser, err := models.SubUser{}.GetSubUserBySubUserId()
			if err != nil {
				return flowRecords, err
			}
			SubUserAndParentUser[userID] = subUser.ParentUserID
			parentUserId = subUser.ParentUserID
		}
		flowRecords, err := models.FlowRecord{
			UserID:     parentUserId,
			IdentityId: 1,
		}.GetByUserSingel()
		if err != nil {
			return flowRecords, err
		}
	}
	return flowRecords, err
}

// 定义锁
var mu sync.Mutex

// 从缓存同步用户流量数据至数据库
func GetUserCurrentFlowData() {
	log.Info("Start Update User Flow Data")
	mu.Lock()
	// 根据流量记录id更新子用户流量数据
	for _, v := range UserCurrentFlow {
		// 更新子用户流量记录
		err := FlowRecord{
			Id:       v.Id,
			UsedFlow: v.UsedFlow,
		}.UpdateByUsedFlow()
		if err != nil {
			log.Error("update Flow Record Error: %v, flowRecordId: %v", err, v.Id)
		}
	}
	UserCurrentFlow = make(map[int64]FlowRecord)
	mu.Unlock()
}
