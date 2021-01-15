package wenda

import "coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"

var (
	caches      = map[uint32]*wendapkg.WendaDetails{} // caches 缓存
	activeGroup = map[uint64]wendapkg.QuestionID{}    // activeGroup 活动的群
)
