package mailer

type Mailer interface {
	SendMail(subject, sender string, destination string, body string) error
}
