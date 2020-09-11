package mailclient_test

import (
	"fmt"
	"testing"

	"github.com/bigwhite/mailclient"
)

type fakeOkMailer struct{}

func (m *fakeOkMailer) SendMail(subject string, dest string, body string) error {
	return nil
}

func TestComposeAndSendOk(t *testing.T) {
	m := &fakeOkMailer{}
	mc := mailclient.New(m)
	_, err := mc.ComposeAndSend("hello, fake test", []string{"xxx@example.com"}, "the test body")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
}

type fakeFailMailer struct{}

func (m *fakeFailMailer) SendMail(subject string, dest string, body string) error {
	return fmt.Errorf("can not reach the mail server of dest [%s]", dest)
}

func TestComposeAndSendFail(t *testing.T) {
	m := &fakeFailMailer{}
	mc := mailclient.New(m)
	_, err := mc.ComposeAndSend("hello, fake test", []string{"xxx@example.com"}, "the test body")
	if err == nil {
		t.Errorf("want non-nil, got nil")
	}
}
