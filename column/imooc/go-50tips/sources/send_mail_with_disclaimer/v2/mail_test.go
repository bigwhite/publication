package mail_test

import (
	"net/smtp"
	"testing"

	mail "github.com/bigwhite/mail"
)

type FakeEmailSender struct {
	subject string
	from    string
	to      []string
	content string
}

func (s *FakeEmailSender) Send(subject, from string,
	to []string, content string, mailserver string, a smtp.Auth) error {
	s.subject = subject
	s.from = from
	s.to = to
	s.content = content
	return nil
}

func TestSendMailWithDisclaimer(t *testing.T) {
	s := &FakeEmailSender{}
	err := mail.SendMailWithDisclaimer(s, "gopher mail test v2",
		"YOUR_MAILBOX",
		[]string{"DEST_MAILBOX"},
		"hello, gopher",
		"smtp.163.com:25",
		smtp.PlainAuth("", "YOUR_EMAIL_ACCOUNT", "YOUR_EMAIL_PASSWD!", "smtp.163.com"))
	if err != nil {
		t.Fatalf("want: nil, actual: %s\n", err)
		return
	}

	want := "hello, gopher" + "\n\n" + mail.DISCLAIMER
	if s.content != want {
		t.Fatalf("want: %s, actual: %s\n", want, s.content)
	}
}
