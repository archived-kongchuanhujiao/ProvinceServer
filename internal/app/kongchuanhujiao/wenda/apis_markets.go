package wenda

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"

	"github.com/kataras/iris/v12/context"
)

type (
	// GetMarketsReq 市场请求
	GetMarketsReq struct {
		Page    uint32 // 页面
		Subject uint8  // 学科
	}

	// PostMarketsReq 市场复制
	PostMarketsReq struct {
		ID     uint32   // 唯一识别码
		Target []uint64 // 目标集
	}
)

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
