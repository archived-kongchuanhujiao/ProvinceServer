package wenda

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kongchuanhujiao/server/internal/app/client"
	"github.com/kongchuanhujiao/server/internal/app/client/message"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/account"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/app/kongchuanhujiao"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

type APIs struct{} // APIs 问答 APIs

type (
	// GetQuestionsReq 获取问题列表或问题 请求结构
	GetQuestionsReq struct {
		Page uint32 // 页面
		ID   uint32 // 标识号
	}

	// GetQuestionsReq 获取问题列表或问题 数据响应结构
	GetQuestionsRes struct {
		Questions []*public.QuestionsTab `json:"questions"`  // 问题
		Groups    *public.Groups         `json:"groups"`     // 群
		GroupName string                 `json:"group_name"` // 群名称
		Members   *public.GroupMembers   `json:"members"`    // 群成员
		Result    *public.Result         `json:"result"`     // 结果
	}
)

// GetQuestions 获取问题列表或问题 APIs。
// 调用方法：GET /apis/wenda/questions
func (a *APIs) GetQuestions(v *GetQuestionsReq, c *context.Context) *kongchuanhujiao.Response {

	// FIXME 需要拆分出更细的颗粒密度
	var (
		d    []*public.QuestionsTab
		g    *public.Groups
		n    string // 群名称
		m    *public.GroupMembers
		calc *public.Result
		err  error
		acct = c.Values().Get("account").(string)
	)

	if v.ID != 0 {

		d, err = wenda.SelectQuestions(&public.QuestionsTab{
			ID:      v.ID,
			Creator: acct,
		}, 0)
		if err != nil {
			return kongchuanhujiao.DefaultErrResp
		}
		t := d[0].Topic.Target
		n = client.GetClient().GetGroupName(t)
		m = client.GetClient().GetGroupMembers(t)
		calc, err = wenda.CalculateResult(v.ID)

	} else {

		d, err = wenda.SelectQuestions(&public.QuestionsTab{Creator: acct}, v.Page)
		g = client.GetClient().GetGroups()

	}
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	return kongchuanhujiao.GenerateSuccResp(
		&GetQuestionsRes{d, g, n, m, calc},
	)
}

// ====================================================================================================================

// PutQuestionStatusReq 更新问题状态 请求结构
type PutQuestionStatusReq struct {
	ID     uint32 // 标识号
	Status uint8  // 状态
}

// PutQuestionsStatus 更新问题状态 APIs。
// 调用方法：PUT /apis/wenda/questions/status
func (a *APIs) PutQuestionsStatus(v *PutQuestionStatusReq, c *context.Context) *kongchuanhujiao.Response {

	qs, err := wenda.SelectQuestions(&public.QuestionsTab{
		ID:      v.ID,
		Creator: c.Values().Get("account").(string),
	}, 0)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	if wenda.UpdateQuestionStatus(qs[0], v.Status) != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// PostQuestions 新建问题 APIs。
// 调用方法：POST /apis/wenda/questions
func (a *APIs) PostQuestions(v *public.QuestionsTab) *kongchuanhujiao.Response {
	if wenda.InsertQuestion(v) != nil {
		return kongchuanhujiao.DefaultErrResp
	}
	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// PutQuestions 更新问题 APIs。
// 调用方法：PUT /apis/wenda/questions
func (a *APIs) PutQuestions(v *public.QuestionsTab) *kongchuanhujiao.Response {
	err := wenda.UpdateQuestion(v)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}
	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// PostPraiseReq 推送表扬列表 请求结构
type PostPraiseReq struct {
	ID uint32 // 标识号
}

// PostPraise 推送表扬列表 APIs。
// 调用方法：POST /apis/wenda/praise
func (a *APIs) PostPraise(v *PostPraiseReq, c *context.Context) *kongchuanhujiao.Response {

	q, err := wenda.SelectQuestions(&public.QuestionsTab{
		ID:      v.ID,
		Creator: c.Values().Get("account").(string),
	}, 0)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	details, err := wenda.CalculateResult(q[0].ID)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	msg := message.NewTextMessage("表扬下列答对的同学：\n")
	for _, mem := range details.Right {
		msg.AddAt(mem)
	}
	client.GetClient().SendMessage(msg.SetTarget(&message.Target{Group: &message.Group{ID: q[0].Topic.Target}}))

	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// PostPushcenterReq 推送数据到钉钉 请求结构
type PostPushcenterReq struct {
	ID     uint32 // 标识号
	Target string // 目标
}

// PostPushcenter 推送数据到钉钉 APIs。
// 调用方法：POST /apis/wenda/pushcenter
func (a *APIs) PostPushcenter(v *PostPushcenterReq, c *context.Context) *kongchuanhujiao.Response {

	ac, err := account.SelectAccount(c.Values().Get("account").(string), 0)
	if err != nil || len(ac) == 0 {
		return kongchuanhujiao.DefaultErrResp
	}

	q, err := wenda.SelectQuestions(&public.QuestionsTab{ID: v.ID}, 0)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}

	err = PushDigestData(q[0])
	if err != nil {
		zap.L().Error("推送作答数据失败", zap.Error(err))
		return kongchuanhujiao.DefaultErrResp
	}

	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// DeleteQuestionsReq 删除问题 请求结构
type DeleteQuestionsReq struct {
	ID uint32 // 标识号
}

// DeleteQuestions 删除问题 APIs。
// 调用方法：Delete /apis/wenda/questions
func (a *APIs) DeleteQuestions(v *DeleteQuestionsReq) *kongchuanhujiao.Response {
	err := wenda.DeleteQuestion(v.ID)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}
	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// GetAnswersReq 获取作答 请求结构
type GetAnswersReq struct {
	ID uint32 // 标识号
}

// GetAnswers 获取作答 APIs。
// 调用方法：GET /apis/wenda/answers
func (a *APIs) GetAnswers(v *GetAnswersReq) *kongchuanhujiao.Response {
	ans, err := wenda.SelectAnswers(v.ID)
	if err != nil {
		return kongchuanhujiao.DefaultErrResp
	}
	return kongchuanhujiao.GenerateSuccResp(ans)
}

// ====================================================================================================================

// PostUploadPicture 上传图片 APIs。
// 调用方法：POST /apis/wenda/upload
func (a *APIs) PostUploadPicture(c *context.Context) *kongchuanhujiao.Response {

	acct := c.Values().Get("account").(string)

	_, fh, err := c.FormFile("file")
	if err != nil {
		zap.L().Warn("解析文件失败", zap.Error(err))
		return kongchuanhujiao.DefaultErrResp
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
	folderName := "assets/pictures/" + acct

	if !Exists(folderName) {
		err = os.MkdirAll(folderName, os.ModePerm)

		if err != nil {
			zap.L().Warn("创建文件夹失败", zap.Error(err))

			return kongchuanhujiao.DefaultErrResp
		}
	}

	// Upload the file to specific destination.
	dest := filepath.Join(folderName+acct, saltedName)
	_, err = c.SaveFormFile(fh, dest)

	if err != nil {
		zap.L().Warn("解析文件失败", zap.Error(err))

		return kongchuanhujiao.DefaultErrResp
	}

	return kongchuanhujiao.DefaultSuccResp
}

// ====================================================================================================================

// GetAnswersReq 获取作答 请求结构
type GetAnswerCSVReq struct {
	ID uint32 // 标识号
}

// GetCsv 获取 CSV APIs。
// 调用方法：GET /apis/wenda/csv
func (a *APIs) GetCsv(v *GetAnswerCSVReq, c *context.Context) {

	ans, err := wenda.SelectAnswers(v.ID)
	if err != nil {
		return
	}
	csv, err := AnswerToCSV(ans)
	if err != nil {
		zap.L().Error("转换答题数据至 CSV 二进制流失败", zap.Error(err))
		return
	}

	c.ContentType("application/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="答题数据详情（`+
		strconv.FormatUint(uint64(v.ID), 10)+`）.csv"`)

	_, _ = c.Write([]byte{239, 187, 191}) // UTF-8 BOM
	_, _ = c.Write(csv)
}
