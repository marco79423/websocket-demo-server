package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/marco79423/websocket-demo-server/internal/utils"
	"golang.org/x/xerrors"
)

func NewJessigodController() IController {
	return &jessigodController{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type jessigodController struct {
	upgrader websocket.Upgrader
}

func (ctrl *jessigodController) Handle(ctx *gin.Context) {
	logger := utils.GetCtxLogger(ctx)

	ws, err := ctrl.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Debug(ctx, "升級為 Websocket 失敗: %w", err)
		return
	}
	defer ws.Close()

	go func() {
		ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if err := ws.WriteMessage(websocket.TextMessage, []byte("我是西卡神，我只想說－－我是最強而且還很愛笑的高嶺之花！")); err != nil {
			logger.Debug(ctx, "寫入 Websocket 訊息失敗: %w", err)
		}
		time.Sleep(3 * time.Second)

		for {
			saying, err := ctrl.getRandomSaying()
			if err != nil {
				logger.Debug(ctx, "取得名言請求失敗: %w", err)
				break
			}

			ws.SetWriteDeadline(time.Now().Add(1 * time.Second))
			if err := ws.WriteMessage(websocket.TextMessage, []byte(saying)); err != nil {
				logger.Debug(ctx, "寫入 Websocket 訊息失敗: %w", err)
				break
			}

			logger.Debug(ctx, "回傳訊息 %v", saying)
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			logger.Debug(ctx, "讀取 Websocket 訊息失敗: %w", err)
			break
		}
	}
}

func (ctrl *jessigodController) getRandomSaying() (string, error) {
	req, err := http.NewRequest("GET", "https://jessigod.marco79423.net/api/random-saying", nil)
	if err != nil {
		return "", xerrors.Errorf("取得隨機名言失敗: %w", err)
	}

	params := url.Values{}
	params.Add("origin", "西卡姐")
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", xerrors.Errorf("取得隨機名言失敗: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", xerrors.Errorf("取得隨機名言失敗: 狀態碼為 %v", resp.StatusCode)
	}

	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf("取得隨機名言失敗: %w", err)
	}

	var jsonData struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return "", xerrors.Errorf("取得隨機名言失敗: %w", err)
	}

	return jsonData.Data.Content, nil
}
