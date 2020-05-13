package lib

import (
	"fmt"
	"log"
	"net/http"

	"github.com/taisukeyamashita/test/lib/env"
)

// Application アプリケーション
type Application struct {
	env env.EnvValues
}

// NewApplication 新しいアプリケーションを作成
func NewApplication(env env.EnvValues) Application {
	return Application{env: env}
}

// Initialize Application初期化
func (application Application) Initialize() {
	application.env.Initialize()
}

// Run Application実行
func (application Application) Run(controllers ...Controller) {
	registControllers := application.BuildRouter(controllers...)
	port := application.env.Port()
	if port == "" {
		port = "8080"
	}
	listenPort := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(listenPort, registControllers))
}

// BuildRouter ルーターを構築する
func (application Application) BuildRouter(controllers ...Controller) http.Handler {
	mux := http.NewServeMux()
	for _, controller := range controllers {
		controller.RegistControllers(mux)
	}
	return mux
}

// Finalize Application終了処理
func (application Application) Finalize() {
	application.env.Finalize()
}
