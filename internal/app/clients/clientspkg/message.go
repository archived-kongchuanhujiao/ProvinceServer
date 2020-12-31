package clientspkg

import "coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"

// NewMessage 新建消息
func NewMessage() *clientspublic.Message {
	return &clientspublic.Message{Chain: []clientspublic.Element{}}
}

func NewTextMessage(s string) *clientspublic.Message { return NewMessage().AddText(s) } // NewTextMessage 新建文本消息

func NewAtMessage(t uint64) *clientspublic.Message { return NewMessage().AddAt(t) } // NewAtMessage 新建@消息

func NewImageMessage(d []byte) *clientspublic.Message { return NewMessage().AddImage(d) } // NewImageMessage 新建图片消息
