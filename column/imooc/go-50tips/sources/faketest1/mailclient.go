package mailclient

import (
	"github.com/bigwhite/mailclient/mailer"
	"github.com/bigwhite/mailclient/sign"
)

type mailClient struct {
	mlr mailer.Mailer
}

func New(mlr mailer.Mailer) *mailClient {
	return &mailClient{
		mlr: mlr,
	}
}

func (c *mailClient) ComposeAndSend(subject string,
	destinations []string, body string) (string, error) {
	signTxt := sign.Get()
	newBody := body + "\n" + signTxt

	for _, dest := range destinations {
		err := c.mlr.SendMail(subject, dest, newBody)
		if err != nil {
			return "", err
		}
	}
	return newBody, nil
}
