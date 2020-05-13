package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/lib/server/provider"
)

const (
	//日付フォーマット
	format = "2006/01/02 15:04:05" // 24h表現、0埋めあり
)

func TestRouter(e *echo.Echo) {
	e.Logger.Print("testUsecase")
	// out, err := exec.Command("gcloud", "config", "list").Output()
	e.GET("/user", insert)
	api := e.Group("/api")
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
