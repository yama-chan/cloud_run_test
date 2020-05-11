package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	server "github.com/taisukeyamashita/test/lib/server/provider"
	"github.com/taisukeyamashita/test/lib/times"
	utils "github.com/taisukeyamashita/test/lib/util"
	"github.com/taisukeyamashita/test/model"
)

const (
	//日付フォーマット
	format = "2006/01/02 15:04:05" // 24h表現、0埋めあり
)

func TestRouter(e *echo.Echo) {
	e.Logger.Print("testUsecase")
	// out, err := exec.Command("gcloud", "config", "list").Output()
	e.GET("/user", func(c echo.Context) error {
		provider := server.NewAppProvider()
		ctx := provider.Context(c.Request())
		user := model.UserInf{
			ID:               utils.CreateUniqueId(),
			Fullname:         "test２",
			LastModifiedDate: times.CurrentTime().Format(format),
		}
		return provider.TestUsecase(ctx).DatastoreOperator.Put(ctx, &user)
	})
	api := e.Group("/api")
	api.GET("/hello", func(c echo.Context) error {
		// return c.JSON(http.StatusOK, "アカウント")
		ctx := c.Request().Context()
		out := "Hello World"
		e.Logger.Print("http://localhost:8080/api/hello")
		return c.String(http.StatusOK, string(out))
	})
}
