package ciyun

import (
	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/ciyun"
)

// HandleWordStat 处理词云
func HandleWordStat(m *message.Message) {

	ct := m.Chain[0].(*message.Text).Content

	// 不处理空消息
	if len(ct) == 0 {
		return
	}

	words := client.GetClient().ExtractWords(ct)

	ciyun.PushWord(m.Target.Group.ID, words)
}
