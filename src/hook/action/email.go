package action

import (
	"fmt"
	"github.com/go-mail/mail"
	"strconv"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

type Email struct {
	Hook        *hook.Hook
	Action      *hook.EmailAction
	SmtpProfile *model.SmtpProfile
	LogModel    *model.ActionLogClient
	EmailModel  *model.EmailClient
}

func NewEmailAction(profile *model.SmtpProfile, action *hook.EmailAction, hook *hook.Hook, email *model.EmailClient, log *model.ActionLogClient) *Email {
	return &Email{
		Action:      action,
		SmtpProfile: profile,
		Hook:        hook,
		LogModel:    log,
		EmailModel:  email,
	}
}

func (email *Email) Send(args map[string]string) {
	start := time.Now().UnixMilli()

	m := mail.NewMessage()
	from := email.SmtpProfile.Sender
	if email.SmtpProfile.Sender == "" {
		email.LogModel.AddDebugLog("Sender is empty, use username instead")
		from = email.SmtpProfile.Username
	}

	body := email.Action.Body
	subject := email.Action.Subject
	for k, v := range args {
		body = strings.ReplaceAll(body, fmt.Sprintf("${%s}", k), v)
		subject = strings.ReplaceAll(subject, fmt.Sprintf("${%s}", k), v)
	}

	m.SetHeader("From", from)
	m.SetHeader("To", email.Action.To...)
	m.SetHeader("Cc", email.Action.Cc...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := mail.NewDialer(email.SmtpProfile.Host, email.SmtpProfile.Port, email.SmtpProfile.Username, email.SmtpProfile.Password)
	d.SSL = email.SmtpProfile.Tls

	if err := d.DialAndSend(m); err != nil {
		email.LogModel.AddErrorLog(err.Error())
		return
	}
	msg := fmt.Sprintf("Use profile %s, send %s to %s, CC to: %s", email.SmtpProfile.Name, email.Action.Body, strings.Join(email.Action.To, ";"), strings.Join(email.Action.Cc, ";"))
	email.LogModel.AddDebugLog(msg)

	end := time.Now().UnixMilli()
	email.LogModel.AddLog("Send successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
}
