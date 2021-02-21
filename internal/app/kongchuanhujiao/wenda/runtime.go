package wenda

import (
	"net/http"

	"github.com/kongchuanhujiao/server/internal/app/datahub/pkg/wenda"
	"github.com/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

type GetRuntimeReq struct { // GetRuntimeReq 运行时请求
	ID uint32 // 唯一识别码
}

// GetRuntime 运行时。
// GET /apis/wenda/runtime 升级为 Websocket
func (a *APIs) GetRuntime(v *GetRuntimeReq, c *context.Context) {

	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := up.Upgrade(c.ResponseWriter(), c.Request(), nil)
	if err != nil {
		logger.Error("升级至 Websocket 失败", zap.Error(err))
		return
	}
	defer conn.Close()

	wenda.AddClient(v.ID, conn)
	defer wenda.RemoveClient(v.ID, conn)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
