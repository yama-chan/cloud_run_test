package routes

import (
	"github.com/taisukeyamashita/test/lib/controller"
)

// Router セキュア機能用ルーター
type Router struct{}

// AddRouters EchoにRouteを追加
func AddRouters(c controller.ControllerBase) {
	TestRouter(c)
	// UserRouter1(e)
	// UserRouter2(e)
	// UserRouter3(e)
	// UserRouter4(e)
}
