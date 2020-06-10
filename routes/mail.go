package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taisukeyamashita/test/lib/config"
	"github.com/taisukeyamashita/test/lib/controller"
	"github.com/taisukeyamashita/test/lib/mail"
	"github.com/taisukeyamashita/test/lib/mail/sendgrid"
	"github.com/taisukeyamashita/test/lib/server"
)

type MailHandler struct {
	provider server.Provider
	echo     *echo.Echo
}

func MailRouter(c controller.ControllerBase) {
	handler := &MailHandler{
		provider: c.Provider,
		echo:     c.Echo,
	}
	mailHandle(handler)
}
func mailHandle(h *MailHandler) {
	h.echo.Logger.Print("Mailrouter")
	api := h.echo.Group("/api")
	api.GET("/mail", h.send)
}

func (h *MailHandler) send(c echo.Context) error {
	ctx := h.provider.Context(c.Request())
	err := sendgrid.Send(ctx, mail.EmailSenderConfig{
		Message: mail.Message{
			Title: "テスト タイトル",
			// Body:  "テスト　ボディ",
			// RegisteredAt: time.Time `validate:"required"` // 作成日時
		},
		FromName:  config.EmailFromName,
		FromEmail: config.EmailFromAddress,
		To: []string{
			"taisuke.yamashita+systest1@topgate.co.jp",
			// "string",
		},
	})
	return c.String(http.StatusOK, fmt.Sprintf("err: %v", err))
}
