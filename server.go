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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//routes.Routerで呼び出し
	routes.Router(e)
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
