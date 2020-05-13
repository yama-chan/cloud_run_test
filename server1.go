package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taisukeyamashita/test/lib/times"
	"github.com/taisukeyamashita/test/routes"
)

func main() {
	// Echoのインスタンスを生成
	e := echo.New()
	// ミドルウェアを設定
	//全てのリクエストについてアクセスログを取得
	e.Use(middleware.Logger())
	//アプリケーションの内部でpanicが発生した場合でも、一律共通のエラーハンドラに処理を飛ばす
	e.Use(middleware.Recover())

	//routes.AddRoutersで呼び出し
	//ルータを振り分けて登録させる
	routes.AddRouters(e)
	// サーバーをポート8080で起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Info(e.Start(":" + port))
}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	times.SetJSTTime()
}
