package sendgrid

import (
	"context"

	"github.com/taisukeyamashita/test/lib/errs"
	"github.com/taisukeyamashita/test/lib/mail"
)

// Send メール送信　EmailSenderConfig(メール送信設定)をうけとって送信
func Send(ctx context.Context, config mail.EmailSenderConfig) error {
	err := ProvideEmailSender().SendEmail(ctx, config)
	if err != nil {
		return errs.WrapXerror(err)
	}
	return nil
}
