package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/marco79423/websocket-demo-server/internal/app"
)

func main() {
	// 初始化
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// 啟動
	if err := application.Start(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := application.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	// 等待關閉訊號
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
