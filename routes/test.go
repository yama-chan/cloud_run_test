package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/gcp/datastore"
	"github.com/taisukeyamashita/test/lib/times"
	utils "github.com/taisukeyamashita/test/lib/util"
	"github.com/taisukeyamashita/test/model"
)

func TestRouter(e *echo.Echo) {
	e.Logger.Print("test")
	// out, err := exec.Command("gcloud", "config", "list").Output()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Fprint(w, string(out))
	api := e.Group("/api")
	api.GET("/test", func(c echo.Context) error {
		// return c.JSON(http.StatusOK, "アカウント")
		ctx := c.Request().Context()
		out := "Hello World"
		e.Logger.Print(string(out))
		e.Logger.Print("http://localhost:8080/api/test")
		// err := firestore.NewFirestore(ctx).Insert(ctx)
		// if err != nil {
		// 	log.Println("failed to insert data to storage.")
		// }

		//日付フォーマット
		const format = "2006/01/02 15:04:05" // 24h表現、0埋めあり
		user := model.UserInf{
			ID:               utils.CreateUniqueId(),
			Fullname:         "test２",
			LastModifiedDate: times.CurrentTime().Format(format),
		}
		err := datastore.NewDatastore(ctx).Put(ctx, &user)
		if err != nil {
			e.Logger.Print("failed to insert data to storage.")
		}
		return c.String(http.StatusOK, string(out))
	})
}
