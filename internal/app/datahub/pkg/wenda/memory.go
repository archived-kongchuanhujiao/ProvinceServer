package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/datahub/internal/memory"
	"coding.net/kongchuanhujiao/server/internal/app/datahub/public/wenda"
)

// GetActiveGroup 获取 ActiveGroup 字段
func GetActiveGroup(k uint64) uint32 {
	return memory.ActiveGroup[k]
}

// GetCaches 获取 Caches 字段
func GetCaches(k uint32) *wenda.Detail {
	return memory.Caches[k]
}

// GetAllActiveGroup 获取 ActiveGroup 数据
func GetAllActiveGroup() map[uint64]uint32 {
	return memory.ActiveGroup
}

// WriteCaches 写入 Cahces 字段
func WriteCaches(k uint32, d *wenda.Detail) {
	memory.Caches[k] = d
}

// WriteActiveGroup 写入 ActiveGroup 字段
func WriteActiveGroup(k uint64, u uint32) {
	memory.ActiveGroup[k] = u
}

// DeleteActiveGroup 删除 ActiveGroup 字段
func DeleteActiveGroup(k uint64) {
	delete(memory.ActiveGroup, k)
}
