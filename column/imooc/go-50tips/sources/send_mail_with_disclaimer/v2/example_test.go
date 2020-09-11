package mail_test

import (
	"fmt"
	"net/smtp"

	mail "github.com/bigwhite/mail"
	email "github.com/jordan-wright/email"
)

type EmailSenderAdapter struct {
	e *email.Email
}

func (adapter *EmailSenderAdapter) Send(subject, from string,
	to []string, content string, mailserver string, a smtp.Auth) error {
	adapter.e.Subject = subject
	adapter.e.From = from
	adapter.e.To = to
	adapter.e.Text = []byte(content)
	return adapter.e.Send(mailserver, a)
}

func ExampleSendMailWithDisclaimer() {
	adapter := &EmailSenderAdapter{
		e: email.NewEmail(),
	}
	err := mail.SendMailWithDisclaimer(adapter, "gopher mail test v2",
		"YOUR_MAILBOX",
		[]string{"DEST_MAILBOX"},
		"hello, gopher",
		"smtp.163.com:25",
		smtp.PlainAuth("", "YOUR_EMAIL_ACCOUNT", "YOUR_EMAIL_PASSWD!", "smtp.163.com"))
	if err != nil {
		fmt.Printf("SendMail error: %s\n", err)
		return
	}
	fmt.Println("SendMail ok")

	// OutPut:
	// SendMail ok
}
