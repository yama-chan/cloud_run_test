package main

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	//簡単なマイクロサービスならこのルーティングを使用できそう
	routes.AddRouters(e)
	// サーバーをポート8080で起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	initializeTime()
}

// ローカル時間を初期化
func initializeTime() {
	// UTCになるので明示的にJST変換する
	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	log.Println("server.go = initialTime() called !!")
}
