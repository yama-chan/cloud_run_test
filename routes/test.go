package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/lib/controller"
	"github.com/taisukeyamashita/test/lib/server"
)

type TestHandler struct {
	provider server.Provider
	echo     *echo.Echo
}

func TestRouter(c controller.ControllerBase) {
	handler := &TestHandler{
		provider: c.Provider,
		echo:     c.Echo,
	}
	handle(handler)
}
func handle(h *TestHandler) {
	h.echo.Logger.Print("testUsecase")
	api := h.echo.Group("/api")
	api.GET("/user", h.insert)
	api.GET("/hello", h.helloWorld)
}

func (h *TestHandler) insert(c echo.Context) error {
	ctx := h.provider.Context(c.Request())
	return h.provider.TestUsecase(ctx).Insert(ctx)
}

func (h *TestHandler) helloWorld(c echo.Context) error {
	out := "Hello World"
	return c.String(http.StatusOK, string(out))
}
