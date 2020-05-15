package server1

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/lib"
	"github.com/taisukeyamashita/test/lib/controller"
	"github.com/taisukeyamashita/test/lib/env"
	"github.com/taisukeyamashita/test/lib/server/provider"
	"github.com/taisukeyamashita/test/routes"
)

func Run() {
	e := echo.New() // Echoのインスタンスを生成

	// アプリケーションの環境情報を設定する/ 環境情報はdeferで最終的に開放する
	application := lib.NewApplication(env.GetEnvValues(env.CreateInitializeConfig()))
	application.Initialize()
	defer application.Finalize()

	// 共通コントローラー作成/ミドルウェアの実行
	controller := controller.NewController(provider.NewAppProvider())

	//ルータを振り分けて登録させる
	routes.AddRouters(controller)

	// サーバーをポート8080で起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// ルーティングしている最中でエラーがあれば、os.Exitが呼ばれる
	e.Logger.Fatal(e.Start(":" + port))
}
