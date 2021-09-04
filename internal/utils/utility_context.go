package utils

import (
	"context"

	"github.com/marco79423/websocket-demo-server/internal/config"
)

func GetCtxLogger(ctx context.Context) ILogger {
	return ctx.Value("logger").(ILogger)
}

func GetCtxConfig(ctx context.Context) config.IConfig {
	return ctx.Value("config").(config.IConfig)
}
