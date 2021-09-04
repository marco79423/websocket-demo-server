package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/marco79423/websocket-demo-server/internal/utils"
)

func NewRouterHandler(ctx context.Context) http.Handler {
	handler := gin.New()
	handler.Use(
		gin.Recovery(),

		// 注入基本資訊
		func(ginCtx *gin.Context) {
			ginCtx.Set("logger", utils.GetCtxLogger(ctx))
			ginCtx.Set("config", utils.GetCtxConfig(ctx))
		},
	)

	setEchoRoute(handler)
	setJessigodRoute(handler)

	return handler
}

func setEchoRoute(handler *gin.Engine) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	handler.GET("/echo", func(ctx *gin.Context) {
		logger := utils.GetCtxLogger(ctx)

		ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
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
	})
}

func setJessigodRoute(handler *gin.Engine) {
	rand.Seed(time.Now().Unix())
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	handler.GET("/jessigod", func(ctx *gin.Context) {
		logger := utils.GetCtxLogger(ctx)

		ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
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
			req.URL.Query().Add("origin", "西卡姐")

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
	})
}
