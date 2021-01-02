package dingtalk

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"strconv"
)

type IDTMessage interface {
	Name() string
}

// DTMessage 钉钉消息结构
// 消息结构详见: https://ding-doc.dingtalk.com/doc#/serverapi2/ye8tup/1babf232
type DTMessage struct {
	MsgType string `json:"msgtype"` // MsgType 消息类型, 可选的有 [text, markdown, actionCard, feedCard, empty]
	At      At     `json:"at"`
}

type DTPlainText struct {
	DTMessage
	Text Text `json:"text"`
}

type DTMarkdown struct {
	DTMessage
	MarkDown MarkDown `json:"markdown"`
}

type DTImage struct {
	DTMessage
	Image Image `json:"image"`
}

// Text 纯文本消息
type Text struct {
	// Content 文本消息
	Content string `json:"content"`
}

// At @ 用户
type At struct {
	// AtUsers 被 @ 的人的手机号
	AtUsers []uint64 `json:"atMobiles"`
	// AtAll 是否 @ 所有人
	AtAll bool `json:"isAtAll"`
}

// MarkDown MD 样式消息
type MarkDown struct {
	// Title 消息标题
	Title string `json:"title"`
	// Text Markdown 格式文本
	Text string `json:"text"`
}

// Image 图片
type Image struct {
	// MediaID 媒体 ID, 用于获取图片
	MediaID string `json:"media_id"`
}

// ActionCard 动作卡片
type ActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
	Btns           []Btn  `json:"btns"`
	BtnOrientation string `json:"btnOrientation"`
	HideAvatar     string `json:"hideAvatar"`
}

// Btn
type Btn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

// ErrResponse 错误时的响应结构体
type ErrResponse struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int64  `json:"errcode"`
}

// 实现方法
func (d DTMessage) Name() string   { return "DTMessage" }
func (d DTPlainText) Name() string { return "DTPlainText" }
func (d DTMarkdown) Name() string  { return "DTMarkdown" }

func (d *DingTalk) transformToChain(ms *clientspublic.Message, m IDTMessage) {
	switch e := m.(type) {
	case DTPlainText:
		ms.Chain = append(ms.Chain, &clientspublic.Text{Content: e.Text.Content})
		for _, v := range e.At.AtUsers {
			ms.Chain = append(ms.Chain, &clientspublic.At{Target: v})
		}
	case DTMarkdown:
		ms.Chain = append(ms.Chain, &clientspublic.Text{Content: "标题: " + e.MarkDown.Title + "\n内容: " + e.MarkDown.Text})
		for _, v := range e.At.AtUsers {
			ms.Chain = append(ms.Chain, &clientspublic.At{Target: v})
		}
	}

	// @TODO Parse to clientspublic.Image
}

func (d *DingTalk) transformToDTMessage(ms *clientspublic.Message) IDTMessage {
	var result IDTMessage

	for _, v := range ms.Chain {
		switch e := v.(type) {
		case *clientspublic.Text:
			if result == nil {
				result = DTPlainText{
					DTMessage: DTMessage{MsgType: "text"},
					Text:      Text{Content: e.Content},
				}
			}
		case *clientspublic.At:
			if result != nil {
				switch r := result.(type) {
				case DTMessage:
					r.At.AtUsers = append(r.At.AtUsers, e.Target)
				case DTPlainText:
					r.At.AtUsers = append(r.At.AtUsers, e.Target)
					for _, user := range r.At.AtUsers {
						r.Text.Content += "@" + strconv.FormatUint(user, 10) + " "
					}
				case DTMarkdown:
					r.At.AtUsers = append(r.At.AtUsers, e.Target)
					for _, user := range r.At.AtUsers {
						r.MarkDown.Text += "@" + strconv.FormatUint(user, 10) + " "
					}
				}
			}
		}
	}

	return result
}
