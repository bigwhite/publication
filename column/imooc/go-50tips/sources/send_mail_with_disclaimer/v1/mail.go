package mail

import (
	"net/smtp"

	email "github.com/jordan-wright/email"
)

const DISCLAIMER = `--------------------------------------------------------
免责声明：此电子邮件和任何附件可能包含特权和机密信息，仅供指定的收件人使用。如果您错误收到此电子邮件，请通知发件人 并立即删除此电子邮件。任何保密性，特权或版权都不会被放弃或丢失，因为此电子邮件是错误地发送给您的。您有责任检查此电子邮件和任何附件是否包含病毒。不保证此材料不含计算机病毒或任何其他缺陷或错误。使用本材料引起的任何损失/损坏不由寄件人负责。发件人的全部责任将仅限于重新提供材料。
--------------------------------------------------------`

func attachDisclaimer(content string) string {
	return content + "\n\n" + DISCLAIMER
}

func SendMailWithDisclaimer(subject, from string, to []string,
	content string, mailserver string,
	a smtp.Auth) error {

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(attachDisclaimer(content))
	return e.Send(mailserver, a)
}
