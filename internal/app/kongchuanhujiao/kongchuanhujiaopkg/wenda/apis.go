package wenda

import (
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

	POSTPUSHCENTERReq struct {
		ID     uint32 `json:"id"`     // 唯一识别码
		Target string `json:"target"` // 目标
	}
)

func (w *APIs) BeforeActivation(_ mvc.BeforeActivation) {}

// GET /apis/wenda/
func (w *APIs) Get() map[string]string { return map[string]string{"status": "online"} }

// 获取问题列表或问题。
// GET /apis/wenda/questions
func (w *APIs) GetQuestions(v *GetQuestionReq) interface{} {
	if v.ID != 0 {
		return wenda.GetQuestions(0, v.ID, false, 0)[0]
	}
	return wenda.GetQuestions(v.Page, v.ID, false, 0)
}

// 获取市场列表。
// GET /apis/wenda/markets
func (w *APIs) GetMarkets(v *GetMarketsReq) []wenda.QuestionListTab {
	return wenda.GetQuestions(v.Page, 0, true, v.Subject)
}

// 推送数据到钉钉。
// POST /apis/wenda/pushcenter
func (w *APIs) PostPUSHCENTER(v *POSTPUSHCENTERReq) {

	/*
		TODO 通过ID 读取问题。
		 然后读取cookie，我也不知道MVC怎么读， 你看看读取user Cookie
		 然后数据库获取对应的教师的钉钉或qq工作群，(这我也还没写，你就假装他有)
		 弄个模板。然后生成再发送消息，两个模板函数最好拆分因为他们没有共同点
	*/

	if v.Target == "dingtalk" {
		// TODO 找 datahub / ~pkg/ 自己弄个发吧
		/*
		 TODO 预期调用
		  这里 -> datahub/pkg/wenda/ -> datahub/internal/dingtalk / 然后消息就发送出去了
		  预期是 internal/dingtalk/ 只能有 accessToken 和 密钥，因为在inernal里获取有可能引入包循环问题
		  所以进入 internal 前先把必要的数据准备好，如第一段所写

		*/

	}

}
