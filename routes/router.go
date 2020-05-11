package routes

import (
	"github.com/labstack/echo/v4"
)

// Router セキュア機能用ルーター
type Router struct{}

// EchoにRouteを追加
func AddRouters(e *echo.Echo) {
	TestRouter(e)
	// UserRouter1(e)
	// UserRouter2(e)
	// UserRouter3(e)
	// UserRouter4(e)
}
