package ciyun

import (
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/ciyun"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

type APIs struct{} // APIs 词云 APIs

// GetWordStatReq 词云请求
type GetWordStatReq struct {
	GID uint64 // 监听群聊号
}

// GetWordStat 词云
// GET /apis/wenda/wordstat
func (a *APIs) GetWordstat(v *GetWordStatReq, c *context.Context) {
	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := up.Upgrade(c.ResponseWriter(), c.Request(), nil)
	if err != nil {
		logger.Error("升级至 Websocket 失败", zap.Error(err))
		return
	}
	defer conn.Close()

	ciyun.AddClient(v.GID, conn)
	defer ciyun.RemoveClient(v.GID, conn)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
