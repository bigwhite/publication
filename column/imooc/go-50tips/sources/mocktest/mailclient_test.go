package mailclient

import (
	"errors"
	"testing"

	"github.com/bigwhite/mailclient/mailer"
	"github.com/golang/mock/gomock"
)

var senderSigns = map[string]string{
	"tonybai@example.com":  "I'm a go programmer",
	"jimxu@example.com":    "I'm a java programmer",
	"stevenli@example.com": "I'm a object-c programmer",
}

func TestComposeAndSendOk(t *testing.T) {
	old := getSign
	sender := "tonybai@example.com"
	timestamp := "Mon, 04 May 2020 11:46:12 CST"

	getSign = func(sender string) string {
		selfSignTxt := senderSigns[sender]
		return selfSignTxt + "\n" + timestamp
	}
	defer func() {
		getSign = old //测试完毕后，恢复原值
	}()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish() //go 1.14版本及以后无需调用该Finish

	mockMailer := mailer.NewMockMailer(mockCtrl)
	mockMailer.EXPECT().SendMail("hello, mock test", sender,
		"dest1@example.com",
		"the test body\n"+senderSigns[sender]+"\n"+timestamp).Return(nil).Times(1)
	mockMailer.EXPECT().SendMail("hello, mock test", sender,
		"dest2@example.com",
		"the test body\n"+senderSigns[sender]+"\n"+timestamp).Return(nil).Times(1)

	mc := New(mockMailer)
	_, err := mc.ComposeAndSend("hello, mock test",
		sender, []string{"dest1@example.com", "dest2@example.com"}, "the test body")
	if err != nil {
		t.Errorf("want nil, got %v", err)
	}
}

var errCanNotReach = errors.New("the mail server of dest address can not be reached")

func TestComposeAndSendFail(t *testing.T) {
	old := getSign
	sender := "tonybai@example.com"
	timestamp := "Mon, 04 May 2020 11:46:12 CST"

	getSign = func(sender string) string {
		selfSignTxt := senderSigns[sender]
		return selfSignTxt + "\n" + timestamp
	}
	defer func() {
		getSign = old //测试完毕后，恢复原值
	}()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish() //go 1.14版本及以后无需调用该Finish

	mockMailer := mailer.NewMockMailer(mockCtrl)
	mockMailer.EXPECT().SendMail("hello, mock test", sender,
		"dest1@example.com", "the test body\n"+senderSigns[sender]+"\n"+timestamp).Return(nil).Times(1)
	mockMailer.EXPECT().SendMail("hello, mock test", sender,
		"dest2@example.com", "the test body\n"+senderSigns[sender]+"\n"+timestamp).Return(errCanNotReach).Times(1)

	mc := New(mockMailer)
	_, err := mc.ComposeAndSend("hello, mock test", sender,
		[]string{"dest1@example.com", "dest2@example.com"}, "the test body")
	if err == nil {
		t.Errorf("want non-nil, got nil")
	}
	if err != errCanNotReach {
		t.Errorf("want %v, got %v", errCanNotReach, err)
	}
}
