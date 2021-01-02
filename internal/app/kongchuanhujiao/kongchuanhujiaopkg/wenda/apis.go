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
