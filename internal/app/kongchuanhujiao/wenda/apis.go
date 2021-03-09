package wenda

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/accounts"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/kataras/iris/v12"
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
		Questions []*public.QuestionsTab `json:"questions"`  // 问题
		Groups    *public.Groups         `json:"groups"`     // 群
		GroupName string                 `json:"group_name"` // 群名称
		Members   *public.GroupMembers   `json:"members"`    // 群成员
		Result    *public.Result         `json:"result"`     // 结果
	}

	PutQuestionStatusReq struct { // PutQuestionStatusReq 问题更新
		ID     uint32 // 唯一识别码
		Status uint8  // 状态
	}

	PostPraiseReq struct { // PostPraiseReq 表扬请求
		ID uint32 // 唯一识别码
	}

	PostPushcenterReq struct { // PostPushcenterReq 推送激活
		ID     uint32 // 唯一识别码
		Target string // 目标
	}

	DeleteQuestionsReq struct{ ID uint32 } // DeleteQuestionsReq 问题删除

	GetAnswersReq struct{ ID uint32 } // GetAnswersReq 获取对应问题答案
)

// GetQuestions 获取问题列表或问题。
// GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionsReq, c *context.Context) *kongchuanhujiao.Response {

	// FIXME 需要拆分出更细的颗粒密度
	var (
		d       []*public.QuestionsTab
		g       *public.Groups
		n       string // 群名称
		m       *public.GroupMembers
		calc    *public.Result
		err     error
		account = c.Values().Get("account").(string)
	)

	if v.ID != 0 {

		d, err = wenda.SelectQuestions(&public.QuestionsTab{
			ID:      v.ID,
			Creator: account,
		}, 0)
		if err != nil {
			return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
		}
		t := d[0].Topic.Target
		n = client.GetClient().GetGroupName(t)
		m = client.GetClient().GetGroupMembers(t)
		calc, err = wenda.CalculateResult(v.ID)

	} else {

		d, err = wenda.SelectQuestions(&public.QuestionsTab{Creator: account}, v.Page)
		g = client.GetClient().GetGroups()

	}
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	return &kongchuanhujiao.Response{
		Message: "ok", Data: &GetQuestionsRes{d, g, n, m, calc},
	}
}

// PutQuestionsStatus 更新问题状态。
// PUT /apis/wenda/questions/status
func (a *APIs) PutQuestionsStatus(v *PutQuestionStatusReq, c *context.Context) *kongchuanhujiao.Response {

	qs, err := wenda.SelectQuestions(&public.QuestionsTab{
		ID:      v.ID,
		Creator: c.Values().Get("account").(string),
	}, 0)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	if wenda.UpdateQuestionStatus(qs[0], v.Status) != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostQuestions 新建问题。
// POST /apis/wenda/questions
func (a *APIs) PostQuestions(v *public.QuestionsTab) *kongchuanhujiao.Response {
	if wenda.InsertQuestion(v) != nil {
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
func (a *APIs) PostPraise(v *PostPraiseReq, c *context.Context) *kongchuanhujiao.Response {

	q, err := wenda.SelectQuestions(&public.QuestionsTab{
		ID:      v.ID,
		Creator: c.Values().Get("account").(string),
	}, 0)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	details, err := wenda.CalculateResult(q[0].ID)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	msg := message.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range details.Right {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&message.Target{Group: &message.Group{ID: q[0].Topic.Target}}))
	return &kongchuanhujiao.Response{Message: "ok"}
}

// PostPushcenter 推送数据到钉钉。
// POST /apis/wenda/pushcenter
func (a *APIs) PostPushcenter(v *PostPushcenterReq, c *context.Context) *kongchuanhujiao.Response {

	ac, err := accounts.SelectAccount(c.Values().Get("account").(string), 0)
	if err != nil || len(ac) == 0 {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	q, err := wenda.SelectQuestions(&public.QuestionsTab{ID: v.ID}, 0)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	if v.Target == "dingtalk" {

		err := PushDigestData("dingtalk", q[0])

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

func (a *APIs) GetAnswers(v *GetAnswersReq) *kongchuanhujiao.Response {
	ans, err := wenda.SelectAnswers(v.ID)
	if err != nil {
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}
	return &kongchuanhujiao.Response{Message: "ok", Data: ans}
}

// UploadPicture 上传图片
// POST /apis/wenda/upload
func (a *APIs) PostUploadPicture(c *context.Context) *kongchuanhujiao.Response {

	account := c.Values().Get("account").(string)

	_, fh, err := c.FormFile("file")
	if err != nil {
		logger.Warn("解析文件失败", zap.Error(err))
		return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
	}

	if fh.Size > 15*iris.MB {
		return &kongchuanhujiao.Response{Status: 1, Message: "上传的文件大小不能超过 15 MB!"}
	}

	fnamePart := strings.Split(fh.Filename, ".")
	saltedName := ""

	for i, n := range fnamePart {
		if i != len(fnamePart)-1 {
			saltedName += n
		}
	}

	saltedName += "_" + HashForSHA1(saltedName) + "." + fnamePart[len(fnamePart)-1]
	folderName := "assets/pictures/" + account

	if !Exists(folderName) {
		err = os.MkdirAll(folderName, os.ModePerm)

		if err != nil {
			logger.Warn("创建文件夹失败", zap.Error(err))

			return &kongchuanhujiao.Response{Status: 1, Message: "服务器错误"}
		}
	}

	// Upload the file to specific destination.
	dest := filepath.Join(folderName+account, saltedName)
	_, err = c.SaveFormFile(fh, dest)

	if err != nil {
		logger.Warn("解析文件失败", zap.Error(err))

		return &kongchuanhujiao.Response{
			Status:  1,
			Message: "服务器错误",
		}
	}

	return &kongchuanhujiao.Response{Status: 0, Message: "ok"}
}
