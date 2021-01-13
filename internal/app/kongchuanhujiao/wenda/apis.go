package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/accounts"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

type (
	APIs struct{} // APIs 问答 APIs

	GetQuestionsReq struct { // GetQuestionsReq 问题请求
		Page uint32 // 页面
		ID   uint32 // 唯一识别码
	}

	PutQuestionStatusReq struct { // PutQuestionStatusReq 问题更新
		ID     uint32 // 唯一识别码
		Status uint8  // 状态
	}

	PostPraisePeq struct { // PostPraisePeq 表扬请求
		ID   uint32   // 唯一识别码
		List []uint64 // 名单
	}

	GetMarketsReq struct { // GetMarketsReq 市场请求
		Page    uint32 // 页面
		Subject uint8  // 学科
	}

	PostMarketsReq struct { // PostMarketsReq 市场复制
		ID     uint32   // 唯一识别码
		Target []uint64 // 目标集
	}

	PostPushcenterReq struct { // PostPushcenterReq 推送新建
		ID     uint32 // 唯一识别码
		Target string // 目标
	}

	DeleteQuestionsReq struct{ ID uint32 } // DeleteQuestionsReq 问题删除
)

// TODO 中间件安全校验

// GetQuestions 获取问题列表或问题。
// GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionsReq, c *context.Context) *kongchuanhujiao.Response {
	var (
		d   []*wenda.QuestionsTab
		err error
	)
	if v.ID != 0 {
		d, err = wenda.SelectQuestions(&wenda.QuestionsTab{ID: v.ID}, 0)
	} else {
		d, err = wenda.SelectQuestions(&wenda.QuestionsTab{ID: v.ID, Creator: c.GetCookie("account")}, v.Page)
	}
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok", Data: d}
}

// PutQuestionsStatus 更新问题状态。
// PUT /apis/wenda/questions/status
func (a *APIs) PutQuestionsStatus(v *PutQuestionStatusReq) *kongchuanhujiao.Response {
	err := wenda.UpdateQuestionStatus(v.ID, v.Status)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostQuestions 新建问题。
// POST /apis/wenda/questions
func (a *APIs) PostQuestions(v *wenda.QuestionsTab) *kongchuanhujiao.Response {
	err := wenda.InsertQuestion(v)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PutQuestions 更新问题。
// PUT /apis/wenda/questions
func (a *APIs) PutQuestions(v *wenda.QuestionsTab) *kongchuanhujiao.Response {
	err := wenda.UpdateQuestion(v)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostPraise 推送表扬列表。
// POST /apis/wenda/praise
func (a *APIs) PostPraise(v *PostPraisePeq) *kongchuanhujiao.Response {
	q, err := wenda.SelectQuestions(&wenda.QuestionsTab{ID: v.ID}, 0)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	msg := clientmsg.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range v.List {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&clientmsg.Target{Group: &clientmsg.Group{ID: q[0].Target}}))
	return &kongchuanhujiao.Response{Message: "ok"}
}

// GetMarkets 获取市场列表。
// GET /apis/wenda/markets
func (a *APIs) GetMarkets(v *GetMarketsReq) *kongchuanhujiao.Response {
	q, err := wenda.SelectQuestions(&wenda.QuestionsTab{Market: true, Subject: v.Subject}, v.Page)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok", Data: q}
}

// PostMarkets 复制市场问题。
// POST /apis/wenda/markets
func (a *APIs) PostMarkets(v *PostMarketsReq, c *context.Context) *kongchuanhujiao.Response {
	user := c.GetCookie("account")
	for _, t := range v.Target {
		err := wenda.CopyQuestions(v.ID, user, t)
		if err != nil {
			return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
		}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostPushcenter 推送数据到钉钉。
// POST /apis/wenda/pushcenter
func (a *APIs) PostPushcenter(v *PostPushcenterReq, c *context.Context) *kongchuanhujiao.Response {

	ac, err := accounts.SelectAccount(c.GetCookie("account"), 0)
	if err != nil || len(ac) == 0 {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	if v.Target == "dingtalk" {

		// FIXME 取消使用问题数据，而是学生作答数据，作答数据结果和作答数据是两张表
		// FIXME 有关作答数据计算结果的内容需要确定
		err := datahubpkg.PushMessage(ac[0].Token, ac[0].Push, ConvertToDTMessage(&wenda.QuestionsTab{}))

		if err != nil {
			logger.Error("发送钉钉消息失败", zap.Error(err))
			return &kongchuanhujiao.Response{Status: 1, Message: "发送失败"}
		}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// DeleteQuestions 删除问题。
// Delete /apis/wenda/questions
func (a *APIs) DeleteQuestions(v *DeleteQuestionsReq) *kongchuanhujiao.Response {
	err := wenda.DeleteQuestion(v.ID)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}
