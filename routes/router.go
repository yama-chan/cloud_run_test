package routes

import (
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	TestRouter(e)
	// UserRouter(e)
}
