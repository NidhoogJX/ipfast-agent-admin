package announcements

import (
	"ipfast_server/internal/handler/network/server"
	"ipfast_server/internal/services"
)

/*
获取公告
*/
func GetAnnouncementsList(resp server.Response) {
	param := struct {
		Page     int `json:"page" binding:"required,min=1"`
		PageSize int `json:"page_size"  binding:"required,min=1,max=100"`
	}{}
	if err := resp.Bind(&param); err != nil {
		resp.Failed("param error")
		return
	}
	announcements, total, err := services.GetAnnouncementsByPage(param.Page, param.PageSize)
	if err != nil {
		resp.Failed("get bulletin failed")
		return
	}
	resp.Res["announcements"] = announcements
	resp.Res["total"] = total
	resp.Success("operate success")
}

// 添加公告信息
func AddAnnouncement(resp server.Response) {
	param := struct {
		Enable  int8   `json:"enable"`
		TopUp   int8   `json:"top_up"`
		Sort    int64  `json:"sort"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	if param.Enable < 0 || param.TopUp < 0 || param.Sort < 0 {
		resp.Failed("param error")
		return
	}
	err = services.AddAnnouncement(param.Sort, param.Enable, param.TopUp, param.Title, param.Content)
	if err != nil {
		resp.Failed("failed to add announcement")
		return
	}
	resp.Success("operate success")
}

// 修改公告信息
func EditAnnouncement(resp server.Response) {
	param := struct {
		ID      int64  `json:"id"`
		Enable  int8   `json:"enable"`
		TopUp   int8   `json:"top_up"`
		Sort    int64  `json:"sort"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	if param.ID < 0 || param.Enable < 0 || param.TopUp < 0 || param.Sort < 0 {
		resp.Failed("param error")
		return
	}
	err = services.EditAnnouncement(param.ID, param.Sort, param.Enable, param.TopUp, param.Title, param.Content)
	if err != nil {
		resp.Failed("failed to edit announcement")
		return
	}
	resp.Success("operate success")
}

// 删除公告
func DeleteAnnouncement(resp server.Response) {
	param := struct {
		Id int64 `json:"id"`
	}{}
	err := resp.Json(&param)
	if err != nil {
		resp.Failed("param error")
		return
	}
	err = services.RemoveAnnouncement(param.Id)
	if err != nil {
		resp.Failed("failed to delete announcement")
		return
	}
	resp.Success("operate success")
}
