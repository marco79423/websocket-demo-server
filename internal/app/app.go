package app

import (
	"context"
	"net"
	"net/http"

	"github.com/marco79423/websocket-demo-server/internal/config"
	"github.com/marco79423/websocket-demo-server/internal/handler"
	"github.com/marco79423/websocket-demo-server/internal/utils"
	"golang.org/x/xerrors"
)

// IApp 主應用的介面
type IApp interface {
	Start() error
	Stop() error
}

// NewApp 建立 App
func NewApp() (IApp, error) {
	// 建立 App 所需的 Context
	appCtx, err := generateAppCtx()
	if err != nil {
		return nil, xerrors.Errorf("建立 App 失敗: %w", err)
	}

	return &app{
		ctx: appCtx,
	}, nil
}

type app struct {
	ctx     context.Context // 應用 Context
	started bool            // 是否啟動

	httpServer *http.Server
}

// Start 啟動主應用 App
func (app *app) Start() error {
	conf := utils.GetCtxConfig(app.ctx)

	logger := utils.GetCtxLogger(app.ctx)
	logger.Info(app.ctx, "開始啟動主應用...")

	// 檢查是否重覆啟動
	if app.started {
		return xerrors.New("應用啟動失敗: 不能重覆啟動主應用")
	}

	// 初始化
	app.httpServer = &http.Server{
		Addr:    conf.GetAddress(),
		Handler: handler.NewRouterHandler(app.ctx),
	}

	// 啟動
	go func() {
		ln, err := net.Listen("tcp", conf.GetAddress())
		if err != nil {
			logger := utils.GetCtxLogger(app.ctx)
			logger.Error(app.ctx, xerrors.Errorf("應用啟動失敗: 啟動服務器失敗: %w", err))
		}

		_ = app.httpServer.Serve(ln)
	}()

	app.started = true
	logger.Info(app.ctx, "成功啟動應用")

	return nil
}

// Stop 關閉 App
func (app *app) Stop() error {
	logger := utils.GetCtxLogger(app.ctx)
	logger.Info(app.ctx, "開始關閉主應用...")

	// 停用所有應用功能
	err := app.httpServer.Shutdown(app.ctx)
	if err != nil {
		return xerrors.Errorf("開始關閉主應用失敗: %w", err)
	}

	app.started = false
	logger.Info(app.ctx, "主應用已關閉")

	return nil
}

// 建立主應用所需的 Context (讀取設定、建立 Logger ...)
func generateAppCtx() (context.Context, error) {
	// 初始化設定
	conf, err := config.NewConfig()
	if err != nil {
		return nil, xerrors.Errorf("建立主應用所需的 Context 失敗: %w", err)
	}

	// 初始化 Logger
	logger, err := utils.NewLogger(conf.GetName(), conf.GetLogLevel())
	if err != nil {
		return nil, xerrors.Errorf("建立主應用所需的 Context 失敗: %w", err)
	}

	// 包進 Context 裡
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", conf)
	ctx = context.WithValue(ctx, "logger", logger)

	return ctx, nil
}
