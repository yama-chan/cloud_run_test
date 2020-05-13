package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/lib/server/provider"
)

func TestRouter(e *echo.Echo) {
	e.Logger.Print("testUsecase")
	api := e.Group("/api")
	api.GET("/user", insert)
	api.GET("/hello", helloWorld)
}

func insert(c echo.Context) error {
	provider := provider.NewAppProvider()
	ctx := provider.Context(c.Request())
	return provider.TestUsecase(ctx).Insert(ctx)
}

func helloWorld(c echo.Context) error {
	out := "Hello World"
	return c.String(http.StatusOK, string(out))
}
