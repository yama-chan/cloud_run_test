package main

import (
	"github.com/taisukeyamashita/test/controller/test"
	"github.com/taisukeyamashita/test/lib"
	"github.com/taisukeyamashita/test/lib/controller"
	"github.com/taisukeyamashita/test/lib/env"
	"github.com/taisukeyamashita/test/lib/server"
	provider2 "github.com/taisukeyamashita/test/lib/server/provider"
)

func main2() {

	application := lib.NewApplication(env.GetEnvValues(env.CreateInitializeConfig()))
	application.Initialize()
	defer application.Finalize()

	var provider server.Provider
	if env.OnLocalDevServer {
		// ローカルサーバの場合
	} else {
		provider = provider2.NewAppProvider()
	}

	controllerBase := controller.NewController(provider)
	application.Run(
		test.NewController(controllerBase),
	)
}

// 初期化を行います。
//echoインスタンスが初期化される前、main関数が実行される後に実行される
func init() {
	// initializeTime()
}

// 初期化
// func initializeTime() {
// 	// UTCになるので明示的にJST変換する
// 	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
// 	log.Println("server.go = initialTime() called !!")
// }
