package internal

import (
	public "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"go.uber.org/zap"
	"strings"
)

// GetGroups 获取群
func (q *QQ) GetGroups() *public.Groups {
	g := public.Groups{}
	for _, v := range q.client.GroupList {
		g[uint64(v.Code)] = v.Name
	}
	return &g
}

// GetGroupName 获取群名称
func (q *QQ) GetGroupName(i uint64) string { return q.client.FindGroup(int64(i)).Name }

// GetGroupMembers 获取群成员
func (q *QQ) GetGroupMembers(i uint64) *public.GroupMembers {
	m := public.GroupMembers{}
	for _, v := range q.client.FindGroup(int64(i)).Members {
		m[uint64(v.Uin)] = v.DisplayName()
	}
	return &m
}

// ExtractWords 分词
func (q *QQ) ExtractWords(s string) (w []string) {

	w, err := q.client.GetWordSegmentation(s)
	if err != nil {
		zap.L().Error("分词失败", zap.Error(err))
	}

	for k, v := range w {
		w[k] = strings.ReplaceAll(v, "\u0000", "")
	}
	return
}
