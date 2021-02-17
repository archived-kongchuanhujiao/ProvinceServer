package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/message"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/pkg/accounts"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"

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

	GetQuestionsRes struct { // GetQuestionsReq 问题响应
		Questions    []*public.QuestionsTab    `json:"questions"`    // 问题
		Groups       *public.Groups            `json:"groups"`       // 群
		GroupName    string                    `json:"group_name"`   // 群名称
		Members      *public.GroupMembers      `json:"members"`      // 群成员
		Calculations []*public.CalculationsTab `json:"calculations"` // 问题计算结果
	}

	PutQuestionStatusReq struct { // PutQuestionStatusReq 问题更新
		ID     uint32 // 唯一识别码
		Status uint8  // 状态
	}

	PostPraisePeq struct { // PostPraisePeq 表扬请求
		ID uint32 // 唯一识别码
	}

	GetMarketsReq struct { // GetMarketsReq 市场请求
		Page    uint32 // 页面
		Subject uint8  // 学科
	}

	PostMarketsReq struct { // PostMarketsReq 市场复制
		ID     uint32   // 唯一识别码
		Target []uint64 // 目标集
	}

	PostPushcenterReq struct { // PostPushcenterReq 推送激活
		ID     uint32 // 唯一识别码
		Target string // 目标
	}

	DeleteQuestionsReq struct{ ID uint32 } // DeleteQuestionsReq 问题删除
)

// TODO 中间件安全校验

// GetQuestions 获取问题列表或问题。
// GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionsReq, c *context.Context) *kongchuanhujiao.Response {

	// FIXME 需要拆分出更细的颗粒密度
	var (
		d   []*public.QuestionsTab
		g   *public.Groups
		n   string // 群名称
		m   *public.GroupMembers
		err error
	)

	if v.ID != 0 {
		d, err = wenda.SelectQuestions(&public.QuestionsTab{ID: v.ID}, 0)
		t := d[0].Target
		n = client.GetClient().GetGroupName(t)
		m = client.GetClient().GetGroupMembers(t)
	} else {
		d, err = wenda.SelectQuestions(&public.QuestionsTab{Creator: c.GetCookie("account")}, v.Page)
		g = client.GetClient().GetGroups()
	}
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	var calcs []*public.CalculationsTab

	for _, e := range d {
		calc, err := wenda.SelectCalculations(e.ID)

		if err != nil {
			return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
		}

		if calc != nil {
			calcs = append(calcs, calc...)
		}
	}

	return &kongchuanhujiao.Response{
		Message: "ok",
		Data:    &GetQuestionsRes{d, g, n, m, calcs},
	}
}

// PutQuestionsStatus 更新问题状态。
// PUT /apis/wenda/questions/status
func (a *APIs) PutQuestionsStatus(v *PutQuestionStatusReq) *kongchuanhujiao.Response {

	var (
		q   = &public.QuestionsTab{ID: v.ID}
		err error
	)

	var qs []*public.QuestionsTab
	qs, err = wenda.SelectQuestions(&public.QuestionsTab{ID: v.ID}, 0)
	q = qs[0]

	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	if wenda.UpdateQuestionStatus(q, v.Status) != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostQuestions 新建问题。
// POST /apis/wenda/questions
func (a *APIs) PostQuestions(v *public.QuestionsTab) *kongchuanhujiao.Response {
	err := wenda.InsertQuestion(v)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PutQuestions 更新问题。
// PUT /apis/wenda/questions
func (a *APIs) PutQuestions(v *public.QuestionsTab) *kongchuanhujiao.Response {
	err := wenda.UpdateQuestion(v)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostPraise 推送表扬列表。
// POST /apis/wenda/praise
func (a *APIs) PostPraise(v *PostPraisePeq) *kongchuanhujiao.Response {
	q, err := wenda.SelectQuestions(&public.QuestionsTab{ID: v.ID}, 0)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	cache, err := wenda.SelectCalculations(q[0].ID)

	if err == nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	details := cache[0]

	msg := message.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range details.Right {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&message.Target{Group: &message.Group{ID: q[0].Target}}))
	return &kongchuanhujiao.Response{Message: "ok"}
}

// GetMarkets 获取市场列表。
// GET /apis/wenda/markets
func (a *APIs) GetMarkets(v *GetMarketsReq) *kongchuanhujiao.Response {
	q, err := wenda.SelectQuestions(&public.QuestionsTab{Market: true, Subject: v.Subject}, v.Page)
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
		err := PushDigestToDingtalk(ac[0].Token, ac[0].Push, ConvertToDTMessage(&public.QuestionsTab{}))

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
