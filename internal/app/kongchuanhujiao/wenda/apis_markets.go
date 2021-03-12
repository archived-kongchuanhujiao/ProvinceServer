package wenda

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"

	"github.com/kataras/iris/v12/context"
)

// GetMarketsReq 市场请求
type GetMarketsReq struct {
	Page    uint32 // 页面
	Subject uint8  // 学科
}

// GetMarkets 获取市场列表 APIs。
// 调用方法：GET /apis/wenda/markets
func (a *APIs) GetMarkets(v *GetMarketsReq) *kongchuanhujiao.Response {
	q, err := wenda.SelectQuestions(&public.QuestionsTab{Market: true, Subject: v.Subject}, v.Page)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}
	return &kongchuanhujiao.Response{Message: "ok", Data: q}
}

// ====================================================================================================================

// PostMarketsReq 复制市场问题 请求结构
type PostMarketsReq struct {
	ID     uint32   // 标识号
	Target []uint64 // 目标集
}

// PostMarkets 复制市场问题 APIs。
// 调用方法：POST /apis/wenda/markets
func (a *APIs) PostMarkets(v *PostMarketsReq, c *context.Context) *kongchuanhujiao.Response {
	for _, t := range v.Target {
		err := wenda.CopyQuestions(v.ID, c.Values().Get("account").(string), t)
		if err != nil {
			return kongchuanhujiao.DefaultErrResp
		}
	}
	return kongchuanhujiao.DefaultErrResp
}
