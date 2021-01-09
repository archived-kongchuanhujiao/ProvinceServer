package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/accounts"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"github.com/kataras/iris/v12/context"

	"go.uber.org/zap"
)

type (
	APIs struct{} // APIs 问答 APIs

	Response struct { // Response 响应
		Status  uint16      `json:"status"`  // 状态码
		Message string      `json:"message"` // 消息
		Data    interface{} `json:"data"`    // 数据
	}

	GetQuestionsReq struct { // GetQuestionsReq 问题请求
		Page uint32 // 页面
		ID   uint32 // 唯一识别码
	}

	PutQuestionReq struct { // PutQuestionReq 问题更新
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
)

// TODO 中间件安全校验

// GetQuestions 获取问题列表或问题。
// GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionsReq, c *context.Context) *Response {
	var (
		d   []*wenda.QuestionsTab
		err error
	)
	if v.ID != 0 {
		d, err = wenda.GetQuestions(&wenda.QuestionsTab{ID: v.ID}, 0)
	} else {
		d, err = wenda.GetQuestions(&wenda.QuestionsTab{ID: v.ID, Creator: c.GetCookie("account")}, v.Page)
	}
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}
	return &Response{0, "ok", d}
}

// PutQuestionsStatus 更新问题状态。
// PUT /apis/wenda/questions/status
func (a *APIs) PutQuestionsStatus(v *PutQuestionReq) *Response {
	err := wenda.UpdateQuestions(v.ID, v.Status)
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}
	return &Response{0, "ok", nil}
}

// PostQuestions 新建问题。
// POST /apis/wenda/questions
func (a *APIs) PostQuestions(v *wenda.QuestionsTab) *Response {
	err := wenda.InsertQuestions(v)
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}
	return &Response{0, "ok", nil}
}

// PostPraise 推送表扬列表。
// POST /apis/wenda/praise
func (a *APIs) PostPraise(v *PostPraisePeq) *Response {
	q, err := wenda.GetQuestions(&wenda.QuestionsTab{ID: v.ID}, 0)
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}
	msg := clientmsg.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range v.List {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&clientmsg.Target{Group: &clientmsg.Group{ID: q[0].Target}}))
	return &Response{0, "ok", nil}
}

// GetMarkets 获取市场列表。
// GET /apis/wenda/markets
func (a *APIs) GetMarkets(v *GetMarketsReq) *Response {
	q, err := wenda.GetQuestions(&wenda.QuestionsTab{Market: true, Subject: v.Subject}, v.Page)
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}
	return &Response{0, "ok", q}
}

// PostMarkets 复制市场问题。
// POST /apis/wenda/markets
func (a *APIs) PostMarkets(v *PostMarketsReq, c *context.Context) *Response {
	user := c.GetCookie("account")
	for _, t := range v.Target {
		err := wenda.CopyQuestions(v.ID, user, t)
		if err != nil {
			return &Response{1, "服务器错误", nil}
		}
	}
	return &Response{0, "ok", nil}
}

// PostPushcenter 推送数据到钉钉。
// POST /apis/wenda/pushcenter
func (a *APIs) PostPushcenter(v *PostPushcenterReq, c *context.Context) *Response {

	ac, err := accounts.GetAccount(c.GetCookie("account"), 0)
	if err != nil {
		return &Response{1, "服务器错误", nil}
	}

	/*
		TODO 弄个答题数据推送模板（独立一个函数）。生成后发送消息
		 两个模板函数拆分出来，因为一个是封装的消息结构体一个是第三方包的钉钉消息结构体
		 空接口（伪泛型）禁止
	*/

	if v.Target == "dingtalk" {
		/*
		 TODO 预期调用
		  这里 -> datahub/datahubpkg/wenda/ -> datahub/internal/dingtalk / 然后消息就发送出去了
		  预期是 internal/dingtalk/ 只能有 accessToken 和 密钥，因为在internal里获取有可能引入包循环问题
		  所以进入 internal 前先把必要的数据准备好（这里理论来说以后会有一张表专门存储答题数据计算结构，推送就是推这个结果的），如第一段所写
		*/

		// FIXME 取消使用问题数据，而是学生作答数据，作答数据结果和作答数据是两张表
		// FIXME 有关作答数据计算结果的内容需要确定
		// FIXME 这里单独一个文件（wenda/push.go）去生成，钉钉生成MarkDown,QQ就是封装的消息链

		err := datahubpkg.PushMessage(ac.Token, ac.Push, "", []string{}, false)
		if err != nil {
			logger.Error("发送钉钉消息失败", zap.Error(err))
			return &Response{Status: 1, Message: "发送失败"}
		}
	}
	return &Response{0, "", nil}
}
