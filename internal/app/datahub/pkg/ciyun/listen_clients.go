package ciyun

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// runtime 使用词云功能的客户端
type runtime map[uint64][]*websocket.Conn

var loggerr = zap.L().Named("数据总线").Named("词云")
var wordCloudRuntime = runtime{}

// AddClient 新增客户端
func AddClient(id uint64, conn *websocket.Conn) {
	wordCloudRuntime[id] = append(wordCloudRuntime[id], conn)
	loggerr.Info("新增连接", zap.String("连接", conn.RemoteAddr().String()))
}

// RemoveClient 移除客户端
func RemoveClient(id uint64, conn *websocket.Conn) {
	for k, v := range wordCloudRuntime[id] {
		if v == conn {
			wordCloudRuntime[id] = append(wordCloudRuntime[id][:k], wordCloudRuntime[id][k+1:]...)
		}
	}
	loggerr.Info("移除连接", zap.String("连接", conn.RemoteAddr().String()))
}

// PushWord 推送词云数据
func PushWord(gid uint64, data []string) {
	if len(wordCloudRuntime) == 0 || len(wordCloudRuntime[gid]) == 0 {
		return
	}

	for _, v := range wordCloudRuntime[gid] {
		err := v.WriteJSON(data)
		if err != nil {
			loggerr.Error("推送词云数据失败", zap.Error(err))
			continue
		}
	}
	loggerr.Info("推送词云数据成功", zap.Uint64("群ID", gid))
}
