package mail_test

import (
	"net/smtp"
	"testing"

	mail "github.com/bigwhite/mail"
)

func TestSendMail(t *testing.T) {
	err := mail.SendMailWithDisclaimer("gopher mail test v1",
		"YOUR_MAILBOX",
		[]string{"DEST_MAILBOX"},
		"hello, gopher",
		"smtp.163.com:25",
		smtp.PlainAuth("", "YOUR_EMAIL_ACCOUNT", "YOUR_EMAIL_PASSWD!", "smtp.163.com"))
	if err != nil {
		t.Fatalf("want: nil, actual: %s\n", err)
	}
}
