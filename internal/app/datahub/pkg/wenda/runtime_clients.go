package wenda

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var wendaRuntime = runtime{}

// AddClient 新增客户端
func AddClient(id uint32, conn *websocket.Conn) {
	wendaRuntime[id] = append(wendaRuntime[id], conn)
	loggerr.Info("新增连接", zap.String("连接", conn.RemoteAddr().String()))
}

// RemoveClient 移除客户端
func RemoveClient(id uint32, conn *websocket.Conn) {
	for k, v := range wendaRuntime[id] {
		if v == conn {
			wendaRuntime[id] = append(wendaRuntime[id][:k], wendaRuntime[id][k+1:]...)
		}
	}
	loggerr.Info("移除连接", zap.String("连接", conn.RemoteAddr().String()))
}

// PushData 推送数据 TODO 数据结构
func PushData(id uint32, data []*AnswersTab) {
	for _, v := range wendaRuntime[id] {
		err := v.WriteJSON(data)
		if err != nil {
			loggerr.Error("推送数据失败", zap.Error(err))
			continue
		}
	}
	loggerr.Info("推送数据成功", zap.Uint32("ID", uint32(id)))
}
