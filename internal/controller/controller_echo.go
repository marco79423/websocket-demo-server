package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/marco79423/websocket-demo-server/internal/utils"
)

func NewEchoController() IController {
	return &echoController{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type echoController struct {
	upgrader websocket.Upgrader
}

func (ctrl *echoController) Handle(ctx *gin.Context) {
	logger := utils.GetCtxLogger(ctx)

	ws, err := ctrl.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Debug(ctx, "升級為 Websocket 失敗: %w", err)
		return
	}
	defer ws.Close()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			logger.Debug(ctx, "讀取 Websocket 訊息失敗: %w", err)
			break
		}

		err = ws.WriteMessage(messageType, message)
		if err != nil {
			logger.Debug(ctx, "寫入 Websocket 訊息失敗: %w", err)
			break
		}

		logger.Debug(ctx, "回傳訊息 [%v] %v", messageType, message)
	}
}
