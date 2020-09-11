package mailer

type Mailer interface {
	SendMail(subject, destination, body string) error
}
