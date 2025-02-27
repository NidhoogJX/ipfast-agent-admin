package test

import (
	"ipfast_server/internal/services"
	"testing"
)

type PullUpStrategy = services.PullUpStrategy

func TestPullUpPreOrders(t *testing.T) {
	// strategy := PullUpStrategy{
	// 	Threshold: 3,
	// 	Intervals: []time.Duration{
	// 		time.Minute,
	// 		5 * time.Minute,
	// 		15 * time.Minute,
	// 	},
	// }

	// uid := int64(1)
	// cid := []services.Cids{
	// 	{
	// 		Cid:              1001,
	// 		TrafficCountryId: 1,
	// 		Quantity:         1,
	// 	},
	// }

	// pid := int64(1001)

	// 第一次拉起
	// payInfo, err := services.PullUpPreOrders(uid, cid, pid, strategy, 1)
	// if err != nil {
	// 	t.Errorf("第一次拉起失败: %v", err)
	// 	return
	// }
	// t.Logf("第一次拉起成功%v", payInfo)
	// return
	// // 第二次拉起
	// err = services.PullUpPreOrders(uid, cid, pid, strategy)
	// if err != nil {
	// 	t.Errorf("第二次拉起失败: %v", err)
	// }
	// t.Logf("第二次拉起成功")

	// // 第三次拉起
	// err = services.PullUpPreOrders(uid, cid, pid, strategy)
	// if err != nil {
	// 	t.Errorf("第三次拉起失败: %v", err)
	// }
	// t.Logf("第三次拉起成功")

	// // 第四次拉起，应该失败并提示等待时间
	// err = services.PullUpPreOrders(uid, cid, pid, strategy)
	// if err == nil {
	// 	t.Errorf("第四次拉起应该失败，但没有失败")
	// } else {
	// 	t.Logf("第四次拉起失败，错误信息: %v", err)
	// }

	// // 等待 1 分钟后再次拉起，应该成功
	// time.Sleep(time.Minute)
	// err = services.PullUpPreOrders(uid, cid, pid, strategy)
	// if err != nil {
	// 	t.Errorf("等待 1 分钟后拉起失败: %v", err)
	// }
	// t.Logf("等待 1 分钟后拉起成功")

	// // 第六次拉起，应该失败并提示等待时间
	// err = services.PullUpPreOrders(uid, cid, pid, strategy)
	// if err == nil {
	// 	t.Errorf("第六次拉起应该失败，但没有失败")
	// } else {
	// 	t.Logf("第六次拉起失败，错误信息: %v", err)
	// }

}
