package routes

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

func TestRouter(e *echo.Echo) {
	log.Print("test")
	out, err := exec.Command("gcloud", "config", "list").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Fprint(w, string(out))

	e.GET("/test", func(c echo.Context) error {
		// return c.JSON(http.StatusOK, "アカウント")
		return c.String(http.StatusOK, string(out))
	})
}
