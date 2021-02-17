package memory

import "coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"

var (
	Caches      = map[uint32]*wenda.Detail{} // Caches 缓存
	ActiveGroup = map[uint64]uint32{}        // ActiveGroup 活动的群
)
