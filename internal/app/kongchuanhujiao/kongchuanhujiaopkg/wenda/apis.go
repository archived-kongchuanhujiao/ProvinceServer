package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"github.com/kataras/iris/v12/mvc"
)

type (
	WendaAPIs struct{} // WendaAPIs 问答 APIs

	GetQuestionReq struct { // GetQuestionReq 问题请求数据
		Page uint32 `json:"page"`
		ID   uint32 `json:"id"`
	}
)

func (w *WendaAPIs) BeforeActivation(_ mvc.BeforeActivation) {}

// GET /apis/wenda/
func (w *WendaAPIs) Get() map[string]string { return map[string]string{"status": "online"} }

// 获取问题列表或问题。
// GET /apis/wenda/questions?id=
func (w *WendaAPIs) GetQuestions(v *GetQuestionReq) interface{} {
	if v.ID != 0 {
		return wenda.GetQuestions(0, v.ID)[0]
	}
	return wenda.GetQuestions(v.Page, v.ID)
}
