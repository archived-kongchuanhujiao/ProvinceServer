package wenda

import "coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"

var (
	Caches      = map[wendapkg.QuestionID]*wendapkg.WendaDetails{} // Caches 缓存
	ActiveGroup = map[uint64]wendapkg.QuestionID{}                 // ActiveGroup 活动的群
)
