package service

import (
	"bytes"
	"html/template"
	"log"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/mailer"
	"github.com/memsbdm/restaurant-api/internal/mailer/mailtemplates"
)

type MailerService interface {
	Send(mail *mailer.Mail) error
	RenderTemplate(name string, data any) (string, error)
}

type mailerService struct {
	cfg    *config.Mailer
	mailer mailer.Mailer
	tmpl   *template.Template
}

func NewMailerService(cfg *config.Mailer, mailer mailer.Mailer) *mailerService {
	return &mailerService{
		cfg:    cfg,
		mailer: mailer,
		tmpl:   loadTemplates(),
	}
}

func (s *mailerService) Send(mail *mailer.Mail) error {
	err := s.mailer.Send(&mailer.Mail{
		To:      mail.To,
		Subject: mail.Subject,
		Body:    mail.Body,
	})

	return err
}

func (s *mailerService) RenderTemplate(name string, data any) (string, error) {
	var buf bytes.Buffer
	err := s.tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("failed to render template %s: %v", name, err)
		return "", err
	}

	return buf.String(), nil
}

func loadTemplates() *template.Template {
	tmpl, err := template.ParseFS(mailtemplates.FS, "*.tmpl")
	if err != nil {
		log.Fatalf("failed to parse mail templates: %v", err)
	}

	return tmpl
}
