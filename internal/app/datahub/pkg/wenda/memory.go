package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/memory"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"
)

func GetActiveGroup(k uint64) uint32 {
	return memory.ActiveGroup[k]
}

func GetCaches(k uint32) *wenda.Detail {
	return memory.Caches[k]
}

func GetAllActiveGroup() map[uint64]uint32 {
	return memory.ActiveGroup
}

func WriteCaches(k uint32, d *wenda.Detail) {
	memory.Caches[k] = d
}

func WriteActiveGroup(k uint64, u uint32) {
	memory.ActiveGroup[k] = u
}

func DeleteActiveGroup(k uint64) {
	delete(memory.ActiveGroup, k)
}
