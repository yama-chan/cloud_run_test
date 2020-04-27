package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TestRouter(e *echo.Echo) {
	log.Print("test")
	// out, err := exec.Command("gcloud", "config", "list").Output()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Fprint(w, string(out))

	e.GET("/test", func(c echo.Context) error {
		// return c.JSON(http.StatusOK, "アカウント")
		out := "Hello World"
		fmt.Println(string(out))
		fmt.Println("http://localhost:8080/test")
		return c.String(http.StatusOK, string(out))
	})
}
