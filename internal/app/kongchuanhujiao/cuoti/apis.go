package cuoti

import "github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"

type (
	APIs struct{} // APIs 账号 APIs

	GetWrongQuestionReq struct {
		UserQQ uint32
	}
)

// GetWrongQuestion 获取错题
func (a *APIs) GetWrongQuestion(v *GetWrongQuestionReq) *kongchuanhujiao.Response {
	_ = v
	// TODO: 数据库交互
	return &kongchuanhujiao.Response{Message: "ok"}
}
