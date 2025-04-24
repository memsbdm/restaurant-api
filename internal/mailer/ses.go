package mailer

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/memsbdm/restaurant-api/config"
)

type Mailer interface {
	Send(mail *Mail) error
}

type mailer struct {
	cfg     *config.Container
	session *ses.SES
}

type Mail struct {
	To      []string
	Subject string
	Body    string
}

func NewSES(cfg *config.Container) *mailer {
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Mailer.Region),
		Credentials: credentials.NewStaticCredentials(cfg.Mailer.AccessKey, cfg.Mailer.SecretKey, ""),
	})
	if err != nil {
		log.Fatalf("error during mailer initialization: %v", err)
	}

	return &mailer{
		cfg:     cfg,
		session: ses.New(awsSession),
	}
}

func (m *mailer) Send(mail *Mail) error {
	if m.cfg.App.Env != config.EnvProduction {
		var builder strings.Builder
		for i := range len(mail.To) {
			builder.WriteString(fmt.Sprintf("- %s", mail.To[i]))
		}
		mail.To = []string{m.cfg.Mailer.DebugTo}
		mail.Body = mail.Body + "Initial email addresses " + builder.String()
	}

	toAddresses := make([]*string, len(mail.To))
	for i, addr := range mail.To {
		toAddresses[i] = aws.String(addr)
	}

	sesInput := &ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: toAddresses},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{Data: aws.String(mail.Body)},
			},
			Subject: &ses.Content{Data: aws.String(mail.Subject)},
		},
		Source: aws.String(m.cfg.Mailer.From),
	}

	msgID, err := m.session.SendEmail(sesInput)
	if err != nil {
		log.Printf("error sending email: %v - msgId: %s", err, msgID.String())
		return err
	}

	log.Printf("mail sent: %s", msgID.String())

	return nil
}
