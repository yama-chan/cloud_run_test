package sendgrid

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	sgMail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/taisukeyamashita/test/lib/config"
	"github.com/taisukeyamashita/test/lib/errs"
	"github.com/taisukeyamashita/test/lib/mail"
)

// EmailSender SendGridAPIでメールの送信を行う
type EmailSender struct {
	// l *logger.Logger
}

// ProvideEmailSender EmailSenderを生成する
// func ProvideEmailSender(l *logger.Logger) *EmailSender {
// 	return &EmailSender{l: l}
// }
func ProvideEmailSender() *EmailSender {
	return &EmailSender{}
}

type sendGridEmailConf mail.EmailSenderConfig

/**
**************************************************************************************
"github.com/sendgrid/sendgrid-go"の実装内容抜粋

// NewSendClient constructs a new Twilio SendGrid client given an API key
func NewSendClient(key string) *Client {
	request := GetRequest(key, "/v3/mail/send", "")
	request.Method = "POST"
	return &Client{request}
}


// Send sends an email through Twilio SendGrid
func (cl *Client) Send(email *mail.SGMailV3) (*rest.Response, error) {
	cl.Body = mail.GetRequestBody(email)
	return MakeRequest(cl.Request)
}

// SGMailV3 contains mail struct
type SGMailV3 struct {
	From             *Email             `json:"from,omitempty"`
	Subject          string             `json:"subject,omitempty"`
	Personalizations []*Personalization `json:"personalizations,omitempty"`
	Content          []*Content         `json:"content,omitempty"`
	Attachments      []*Attachment      `json:"attachments,omitempty"`
	TemplateID       string             `json:"template_id,omitempty"`
	Sections         map[string]string  `json:"sections,omitempty"`
	Headers          map[string]string  `json:"headers,omitempty"`
	Categories       []string           `json:"categories,omitempty"`
	CustomArgs       map[string]string  `json:"custom_args,omitempty"`
	SendAt           int                `json:"send_at,omitempty"`
	BatchID          string             `json:"batch_id,omitempty"`
	Asm              *Asm               `json:"asm,omitempty"`
	IPPoolID         string             `json:"ip_pool_name,omitempty"`
	MailSettings     *MailSettings      `json:"mail_settings,omitempty"`
	TrackingSettings *TrackingSettings  `json:"tracking_settings,omitempty"`
	ReplyTo          *Email             `json:"reply_to,omitempty"`
}

// Personalization holds mail body struct
type Personalization struct {
	To                  []*Email               `json:"to,omitempty"`
	CC                  []*Email               `json:"cc,omitempty"`
	BCC                 []*Email               `json:"bcc,omitempty"`
	Subject             string                 `json:"subject,omitempty"`
	Headers             map[string]string      `json:"headers,omitempty"`
	Substitutions       map[string]string      `json:"substitutions,omitempty"`
	CustomArgs          map[string]string      `json:"custom_args,omitempty"`
	DynamicTemplateData map[string]interface{} `json:"dynamic_template_data,omitempty"`
	Categories          []string               `json:"categories,omitempty"`
	SendAt              int                    `json:"send_at,omitempty"`
}

**************************************************************************************
**/

// SendEmail Eメール送信
func (e EmailSender) SendEmail(ctx context.Context, conf mail.EmailSenderConfig) error {
	sConf := sendGridEmailConf(conf)

	var (
		client = sendgrid.NewSendClient(config.SendGridAPIKey)
		from   = sConf.createFrom()
		tos    = sConf.createToPersonalizations()
	)
	message := sConf.createMailMessage(from, tos)
	log.Print("email is about to send", message)
	resp, err := client.Send(message)
	log.Print("email response: ", resp)
	if err != nil {
		return errs.WrapXerror(err)
	}
	if err != nil {
		return err
	} else if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("SendGrid送信エラー: code=%v", resp.StatusCode)
	}
	log.Print("email send success")
	return nil
}

func (s sendGridEmailConf) createFrom() *sgMail.Email {
	return sgMail.NewEmail(s.FromName, s.FromEmail)
}

func (s sendGridEmailConf) createToPersonalizations() []*sgMail.Personalization {
	personalizations := make([]*sgMail.Personalization, len(s.To))
	for i, to := range s.To {
		p := sgMail.NewPersonalization()
		p.AddTos(sgMail.NewEmail("", to)) // とりあえず名前は未設定で
		data := struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			Email: "test@example.com",
			Name:  to,
		}
		p.SetDynamicTemplateData("data", data)
		personalizations[i] = p
	}
	return personalizations
}

func (s sendGridEmailConf) createMailMessage(from *sgMail.Email, tos []*sgMail.Personalization) *sgMail.SGMailV3 {
	message := sgMail.NewV3Mail()
	// ※トランザクショナルテンプレートを使用してリクエストにtemplate_IDを指定している場合はContentは必須ではありません。
	// message.AddContent(sgsgMail.NewContent(string(s.ContentType), s.Message.Body))
	message.Subject = s.Message.Title
	message.SetFrom(from)
	message.AddPersonalizations(tos...)
	message.SetTemplateID(config.SendGridTemplateID)
	return message
}
