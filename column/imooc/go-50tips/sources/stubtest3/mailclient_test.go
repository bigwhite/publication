package mailclient

import (
	"strings"
	"testing"

	"github.com/prashantv/gostub"
)

type fakeOkMailer struct{}

func (m *fakeOkMailer) SendMail(subject string, sender string, dest string, body string) error {
	return nil
}

var senderSigns = map[string]string{
	"tonybai@example.com":  "I'm a go programmer",
	"jimxu@example.com":    "I'm a java programmer",
	"stevenli@example.com": "I'm a object-c programmer",
}

func TestComposeAndSendWithSign(t *testing.T) {
	sender := "tonybai@example.com"
	timestamp := "Mon, 04 May 2020 11:46:12 CST"

	stubs := gostub.Stub(&getSign, func(sender string) string {
		selfSignTxt := senderSigns[sender]
		return selfSignTxt + "\n" + timestamp
	})
	defer stubs.Reset()

	m := &fakeOkMailer{}
	mc := New(m)
	body, err := mc.ComposeAndSend("hello, stub test", sender,
		[]string{"xxx@example.com"}, "the test body")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}

	if !strings.Contains(body, timestamp) {
		t.Errorf("the sign of the mail does not contain [%s]", timestamp)
	}

	if !strings.Contains(body, senderSigns[sender]) {
		t.Errorf("the sign of the mail does not contain [%s]", senderSigns[sender])
	}

	sender = "jimxu@example.com"
	body, err = mc.ComposeAndSend("hello, stub test", sender,
		[]string{"xxx@example.com"}, "the test body")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}

	if !strings.Contains(body, senderSigns[sender]) {
		t.Errorf("the sign of the mail does not contain [%s]", senderSigns[sender])
	}
}
