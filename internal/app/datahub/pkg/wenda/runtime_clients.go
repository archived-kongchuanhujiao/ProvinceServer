package wenda

import (
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// wrapper 封装
type wrapper struct {
	Result *wenda.Result `json:"result"`
}

var loggerr = logger.Named("数据总线").Named("问答")
var wendaRuntime = wenda.Runtime{}

// AddClient 新增客户端
func AddClient(id uint32, gid uint64, conn *websocket.Conn) {
	wendaRuntime = append(wendaRuntime, &wenda.WSConn{
		Conn:           conn,
		ListenQuestion: id,
		ListenGroup:    gid,
	})
	loggerr.Info("新增连接", zap.String("连接", conn.RemoteAddr().String()))
}

// RemoveClient 移除客户端
func RemoveClient(id uint32, gid uint64, conn *websocket.Conn) {
	for k, v := range wendaRuntime {
		if v.Conn == conn {
			switch {
			case v.ListenQuestion == id && id != 0:
			case v.ListenGroup == gid && gid != 0:
				wendaRuntime = append(wendaRuntime[:k], wendaRuntime[k+1:]...)
			}
		}
	}
	loggerr.Info("移除连接", zap.String("连接", conn.RemoteAddr().String()))
}

// PushWords 推送消息
func PushWords(gid uint64, data []string) {
	for _, v := range wendaRuntime {
		if v.ListenGroup == gid && v.ListenQuestion == 0 {
			err := v.Conn.WriteJSON(data)
			if err != nil {
				loggerr.Error("推送数据失败", zap.Error(err))
				continue
			}
		}
	}
	loggerr.Info("推送词云数据成功", zap.Uint64("群", gid))
}

// PushData 推送数据
func PushData(id uint32, data *wenda.Result) {
	for _, v := range wendaRuntime {
		if v.ListenQuestion == id && v.ListenGroup == 0 {
			err := v.Conn.WriteJSON(wrapper{data})
			if err != nil {
				loggerr.Error("推送数据失败", zap.Error(err))
				continue
			}
		}
	}
	loggerr.Info("推送数据成功", zap.Uint32("ID", id))
}
