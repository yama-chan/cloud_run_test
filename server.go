package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taisukeyamashita/test/routes"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

func main() {

	// Echoのインスタンスを生成
	e := echo.New()
	//routes.Routerで呼び出し
	routes.Router(e)

	//セッションを設定
	store := sessions.NewCookieStore([]byte("secret"))
	//セッション保持時間
	store.MaxAge(86400)
	//sessionオプション設定
	store.Options = &sessions.Options{
		Path: "/",
		// MaxAge=0 means no Max-Age attribute specified and the cookie will be
		// deleted after the browser session ends.
		// MaxAge<0 means delete cookie immediately.
		// MaxAge>0 means Max-Age attribute present and given in seconds.

		//ログイン(/login)画面遷移時にセッションを削除
		// MaxAge: -1,
		//７日間セッション(ログイン状態)を維持
		// MaxAge:   86400 * 7,
		//５分間セッション(ログイン状態)を維持
		MaxAge:   3000 * 7,
		HttpOnly: true,
	}
	//session
	e.Use(session.Middleware(store))

	// renderer := &TemplateRenderer{
	// 	temp: template.Must(template.ParseGlob("view/*.html")),
	// }
	// テンプレートを利用するためのRendererの設定
	// e.Renderer = renderer

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静的ファイルのパスを設定
	e.Use(middleware.Static(""))
	e.Use(middleware.Static("/static"))

	// サーバーをできればポート8080で起動
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Print("main.go = init() >>>>>>>> loadTemplates() called !! \n")
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	loadTemplates()
}
