//go:generate mockgen -source=./mailer.go -destination=./mock_mailer.go -package=mailer Mailer

package mailer

type Mailer interface {
	SendMail(subject, sender, destination, body string) error
}
