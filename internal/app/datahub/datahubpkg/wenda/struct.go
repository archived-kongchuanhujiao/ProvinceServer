package wenda

import (
	"coding.net/kongchuanhujiao/server/internal/app/kongchuanhujiao/public/wendapkg"
	"coding.net/kongchuanhujiao/server/internal/pkg/logger"
	"github.com/gorilla/websocket"
)

var loggerr = logger.Named("数据总线").Named("问答")

type (
	runtime map[wendapkg.QuestionID][]*websocket.Conn // runtime 运行时
)
