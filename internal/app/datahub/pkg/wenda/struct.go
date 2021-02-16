package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"github.com/gorilla/websocket"
)

var loggerr = logger.Named("数据总线").Named("问答")

type (
	runtime map[uint32][]*websocket.Conn // runtime 运行时
)
