package qq

import (
	"coding.net/kongchuanhujiao/server/internal/app/clients/clientspublic"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"go.uber.org/zap"
)

// transformToChain 转化为消息链
func (q *QQ) transformToChain(ms *clientspublic.Message, m []message.IMessageElement) {
	for _, v := range m {
		switch e := v.(type) {
		case *message.TextElement:
			ms.Chain = append(ms.Chain, &clientspublic.Text{Content: e.Content})
		case *message.AtElement:
			ms.Chain = append(ms.Chain, &clientspublic.At{Target: uint64(e.Target)})
		case *message.ImageElement:
			ms.Chain = append(ms.Chain, &clientspublic.Image{URL: e.Url})
		}
	}
}

// receiveGroupMessage 接收群消息
func (q *QQ) receiveGroupMessage(_ *client.QQClient, m *message.GroupMessage) {

	ms := &clientspublic.Message{
		Client: q,
		Chain:  []clientspublic.Element{},
		Target: &clientspublic.Target{
			ID:    uint64(m.Sender.Uin),
			Name:  m.Sender.DisplayName(),
			Group: &clientspublic.Group{ID: uint64(m.GroupCode), Name: m.GroupName},
		},
	}

	q.transformToChain(ms, m.Elements)
	q.ReceiveMessage(ms)
}

// receiveFriendMessage 接收好友消息
func (q *QQ) receiveFriendMessage(_ *client.QQClient, m *message.PrivateMessage) {

	ms := &clientspublic.Message{
		Client: q,
		Chain:  []clientspublic.Element{},
		Target: &clientspublic.Target{ID: uint64(m.Sender.Uin), Name: m.Sender.DisplayName()},
	}

	q.transformToChain(ms, m.Elements)
	q.ReceiveMessage(ms)
}

// transformToMiraiGO 转化为 MiraiGO
func (q *QQ) transformToMiraiGO(ms *clientspublic.Message) (m *message.SendingMessage) {

	m = &message.SendingMessage{
		Elements: []message.IMessageElement{},
	}

	for _, v := range ms.Chain {
		switch e := v.(type) {
		case *clientspublic.Text:
			m.Elements = append(m.Elements, message.NewText(e.Content))
		case *clientspublic.At:
			mem := q.client.FindGroupByUin(int64(ms.Target.Group.ID)).FindMember(int64(ms.Target.ID))
			m.Elements = append(m.Elements, message.NewAt(int64(e.Target), mem.DisplayName()))
		case *clientspublic.Image:
			se, err := q.client.UploadGroupImage(int64(ms.Target.Group.ID), e.Data)
			if err != nil {
				loggerr.Error("上传图片错误", zap.Error(err))
			}
			m.Elements = append(m.Elements, se)
		}
	}
	return
}
