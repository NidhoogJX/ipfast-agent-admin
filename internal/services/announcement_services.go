package services

import (
	"fmt"
	"time"

	"ipfast_server/internal/db/models"
	"ipfast_server/pkg/util/log"
)

// GetAnnouncementsByPage 获取分页公告列表
func GetAnnouncementsByPage(page int, pageSize int) (res []models.Announcement, total int64, err error) {
	res, total, err = models.Announcement{}.SelectByPage(page, pageSize)
	if err != nil {
		log.Error("GetAnnouncementsByPage error:%v", err)
		err = fmt.Errorf("failed to obtain system announcement")
		return
	}
	return
}

// 添加公告信息
func AddAnnouncement(sort int64, enable, topUp int8, title, content string) (err error) {
	now := time.Now().Unix()
	err = models.Announcement{
		Sort:        sort,
		Enable:      enable,
		TopUp:       topUp,
		Title:       title,
		Content:     content,
		CreatedTime: now,
		UpdatedTime: now,
	}.Create()
	if err != nil {
		err = fmt.Errorf("failed to add system announcement")
	}
	return
}

// 修改公告信息
func EditAnnouncement(id, sort int64, enable, topUp int8, title, content string) (err error) {
	now := time.Now().Unix()
	err = models.Announcement{
		ID:          id,
		Sort:        sort,
		Enable:      enable,
		TopUp:       topUp,
		Title:       title,
		Content:     content,
		UpdatedTime: now,
	}.Update()
	if err != nil {
		err = fmt.Errorf("failed to edit system announcement")
	}
	return
}

// 删除公告信息
func RemoveAnnouncement(id int64) (err error) {
	err = models.Announcement{
		ID: id,
	}.Delete()
	if err != nil {
		err = fmt.Errorf("failed to edit system announcement")
	}
	return
}
