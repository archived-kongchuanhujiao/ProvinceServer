package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/client"
	"coding.net/kongchuanhujiao/server/internal/app/client/clientmsg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"github.com/kataras/iris/v12/mvc"
)

type (
	APIs struct{} // APIs 问答 APIs

	GetQuestionReq struct { // GetQuestionReq 问题请求数据
		Page uint32 `json:"page"` // 页面
		ID   uint32 `json:"id"`   // 唯一识别码
	}

	GetMarketsReq struct { // GetMarketsReq 市场请求数据
		Page    uint32 `json:"page"` // 页面
		Subject uint8  `json:"sub"`  // 学科
	}

	POSTPraisePeq struct { // POSTPraisePeq 表扬请求数据
		ID   uint32   `json:"id"`   // 唯一识别码
		List []uint64 `json:"list"` // 名单
	}

	POSTPUSHCENTERReq struct {
		ID     uint32 `json:"id"`     // 唯一识别码
		Target string `json:"target"` // 目标
	}
)

func (a *APIs) BeforeActivation(_ mvc.BeforeActivation) {}

// 获取状态。
// GET /apis/wenda/
func (a *APIs) Get() map[string]string { return map[string]string{"status": "online"} }

// 获取问题列表或问题。
// GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionReq) interface{} {
	if v.ID != 0 {
		return wenda.GetQuestions(0, v.ID, false, 0)[0]
	}
	return wenda.GetQuestions(v.Page, v.ID, false, 0)
}

// 推送表扬列表。
// POST /apis/wenda/praise
func (a *APIs) PostPraise(v *POSTPraisePeq) {
	q := wenda.GetQuestions(0, v.ID, false, 0)
	msg := clientmsg.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range v.List {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&clientmsg.Target{Group: &clientmsg.Group{ID: q[0].Target}}))
}

// 获取市场列表。
// GET /apis/wenda/markets
func (a *APIs) GetMarkets(v *GetMarketsReq) []wenda.QuestionListTab {
	return wenda.GetQuestions(v.Page, 0, true, v.Subject)
}

// 推送数据到钉钉。
// POST /apis/wenda/pushcenter
func (a *APIs) PostPUSHCENTER(v *POSTPUSHCENTERReq) {

	/*
		TODO 通过ID 读取问题。
		 然后读取cookie，我也不知道MVC怎么读， 你看看读取user Cookie
		 然后数据库获取对应的教师的钉钉或qq工作群，(这我也还没写，你就假装他有)
		 弄个模板。然后生成再发送消息，两个模板函数最好拆分因为他们没有共同点
	*/

	if v.Target == "dingtalk" {
		/*
		 TODO 预期调用
		  这里 -> datahub/pkg/wenda/ -> datahub/internal/dingtalk / 然后消息就发送出去了
		  预期是 internal/dingtalk/ 只能有 accessToken 和 密钥，因为在internal里获取有可能引入包循环问题
		  所以进入 internal 前先把必要的数据准备好，如第一段所写
		*/

		datahubpkg.PushMessage("fakeToken", "fakeSecret", "", []string{}, false)
	}

}
