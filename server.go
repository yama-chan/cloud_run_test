package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taisukeyamashita/test/routes"
)

func main() {
	// Echoのインスタンスを生成
	e := echo.New()
	//routes.Routerで呼び出し
	routes.Router(e)
	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// サーバーをできればポート8080で起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	loadTemplates()
}
