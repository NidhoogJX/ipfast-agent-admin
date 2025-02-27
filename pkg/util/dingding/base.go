package dingding

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// 钉钉群机器人的token
const token = "f7ce6c5198f0ce932bee9bef2d9f26bb5cfca977fa18d0866ab8a746c2e4a7b6"

type TextAt struct {
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type MarkdownAt struct {
	AtMobiles []string `json:"atMobiles"` // 根据手机号@
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextMessage struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At TextAt `json:"at"`
}

type MarkdownMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At MarkdownAt
}

func GetWebhookURL(token string) string {
	return "https://oapi.dingtalk.com/robot/send?access_token=" + token
}

// SendNotification 发送Text消息通知到钉钉
// isatall 是否@所有人
// message 是要发送的消息内容
// 返回 HTTP 响应的状态码和错误（如果有的话）
func SendText(message, webhookURL string, isatall bool) (int, error) {
	notification := TextMessage{
		Msgtype: "text",
		At: TextAt{
			IsAtAll: isatall,
		},
	}
	notification.Text.Content = message
	return SendJsonToDingDing(notification, webhookURL)
}

// SendMarkdown 发送Markdown消息通知到钉钉
// isatall 是否@所有人
// message 是要发送的消息内容 (markdown格式) 支持格式如下
// 标题
// # 一级标题
// ## 二级标题
// ### 三级标题
// #### 四级标题
// ##### 五级标题
// ###### 六级标题
// 引用
// > A man who stands for nothing will fall for anything.
// 文字加粗、斜体
// **bold**
// *italic*
// 链接
// [this is a link](http://name.com)
// 图片（建议不要超过20张）
// ![](http://name.com/pic.jpg)
// 无序列表
// - item1
// - item2
// 有序列表
// 1. item1
// 2. item2
// 返回 HTTP 响应的状态码和错误（如果有的话）

func SendMarkdown(title, message, webhookURL string, isatall bool) (int, error) {
	notification := MarkdownMessage{
		Msgtype: "markdown",
		At: MarkdownAt{
			IsAtAll: isatall,
		},
	}
	notification.Markdown.Title = title
	notification.Markdown.Text = message
	return SendJsonToDingDing(notification, webhookURL)
}

func SendJsonToDingDing(notification interface{}, webhookURL string) (int, error) {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(notificationBytes))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
