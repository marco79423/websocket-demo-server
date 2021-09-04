package controller

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/marco79423/websocket-demo-server/internal/utils"
)

func NewNabigodController() IController {
	return &nabigodController{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type nabigodController struct {
	upgrader websocket.Upgrader
}

func (ctrl *nabigodController) Handle(ctx *gin.Context) {
	logger := utils.GetCtxLogger(ctx)

	ws, err := ctrl.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Debug(ctx, "升級為 Websocket 失敗: %w", err)
		return
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			logger.Debug(ctx, "關閉 Websocket 失敗: %w", err)
			return
		}
	}()

	client := &http.Client{}
	for {
		req, err := http.NewRequest("GET", "https://jessigod.marco79423.net/api/sayings", nil)
		if err != nil {
			logger.Debug(ctx, "取得名言請求失敗: %w", err)
			break
		}
		req.Header.Set("Authorization", "Jessi bac")
		req.URL.Query().Add("origin", "Nabi 姐")

		resp, err := client.Do(req)
		if err != nil {
			logger.Debug(ctx, "取得名言失敗: %w", err)
			break
		}

		if resp.StatusCode != http.StatusOK {
			logger.Debug(ctx, "取得名言失敗: 狀態碼為 %v", resp.StatusCode)
			break
		}

		rawData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Debug(ctx, "解析回傳內容失敗: %w", err)
			break
		}

		var jsonData struct {
			Data []struct {
				Content string `json:"content"`
			} `json:"data"`
		}

		if err := json.Unmarshal(rawData, &jsonData); err != nil {
			logger.Debug(ctx, "解析 JSON 內容失敗: %w", err)
			break
		}

		saying := jsonData.Data[rand.Intn(len(jsonData.Data))].Content
		err = ws.WriteMessage(websocket.TextMessage, []byte(saying))
		if err != nil {
			logger.Debug(ctx, "寫入 Websocket 訊息失敗: %w", err)
			break
		}

		logger.Debug(ctx, "回傳訊息 %v", saying)

		time.Sleep(3 * time.Second)
	}
}
