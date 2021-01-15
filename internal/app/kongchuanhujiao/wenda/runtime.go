package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"net/http"
	"strconv"

	"coding.net/kongchuanhujiao/server/internal/app/datahub/datahubpkg/wenda"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

// GetRuntime 运行时。
// GET /apis/wenda/runtime 升级为 Websocket
func (a *APIs) GetRuntime(c *context.Context) {

	up := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := up.Upgrade(c.ResponseWriter(), c.Request(), nil)
	if err != nil {
		logger.Error("升级至 Websocket 失败", zap.Error(err))
		return
	}
	defer conn.Close()

	_, msg, err := conn.ReadMessage()
	if err != nil {
		logger.Error("读取消息失败", zap.Error(err))
		return
	}

	i, err := strconv.ParseUint(string(msg), 10, 32)
	if err != nil {
		logger.Error("解析问题 ID 失败", zap.Error(err))
		return
	}
	id := uint32(i)

	wenda.AddClient(wendapkg.QuestionID(id), conn)
	defer wenda.RemoveClient(wendapkg.QuestionID(id), conn)
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}
