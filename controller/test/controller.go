package test

import (
	"net/http"

	"github.com/taisukeyamashita/test/lib"
	"github.com/taisukeyamashita/test/lib/controller"
	"github.com/taisukeyamashita/test/lib/logger/logrus"
	"github.com/taisukeyamashita/test/routes"
)

// Controller 管理機能用コントローラー
type Controller struct {
	// https://medium.com/eureka-engineering/golang-embedded-ac43201cf772
	// lib.Controllerインタフェースを満たすようにlib.ControllerBase構造体を埋め込む(※フィールド名を記載しないこと)
	// 埋め込んだことにより、lib.ControllerBase構造体で定義してるメソッドも実行できる
	controller.ControllerBase
	controllerName string
}

// Controller implement lib.Controller
var _ lib.Controller = Controller{}

// Controller implement http.Handler
var _ http.Handler = Controller{}

// NewController 機能用コントローラー作成
func NewController(controller controller.ControllerBase) Controller {
	return Controller{controller, "testController"}
}

// RegistControllers コントローラーを登録する
func (controller Controller) RegistControllers(mux *http.ServeMux) {
	basePath := "/api/test/"
	// 埋め込まれている'lib.ControllerBase構造体'で定義してる'フィールド'にも直接アクセスできる
	adminGroup := controller.Echo.Group(basePath)
	// 埋め込まれている'lib.ControllerBase構造体'で定義してる'メソッド'も直接実行できる
	controller.AddRoutes(adminGroup, routes.Router{})
	mux.Handle(basePath, controller)
	// TODO:  providerにloggerをもたせてログ出力を行えるようにする
	logrus.NewLogger().Info(controller.controllerName + " is registered")
}
