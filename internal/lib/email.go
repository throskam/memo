package lib

import (
	"bytes"
	"net/smtp"

	"github.com/resend/resend-go/v3"
)

type Email struct {
	From    string
	To      []string
	Subject string
	Body    string
	SMTP    string
	APIKey  string
}

type EmailOption func(*Email)

func WithSMTP(smtp string) EmailOption {
	return func(email *Email) {
		email.SMTP = smtp
	}
}

func WithAPIKey(apiKey string) EmailOption {
	return func(email *Email) {
		email.APIKey = apiKey
	}
}

func SendEmail(from string, to []string, subject string, body string, options ...EmailOption) error {
	email := Email{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	for _, o := range options {
		o(&email)
	}

	if email.SMTP != "" {
		return sendEmailWithSMTP(email)
	}

	return sendEmailWithResend(email)
}

func sendEmailWithSMTP(email Email) error {
	addr := email.SMTP

	content := &bytes.Buffer{}

	content.WriteString("Subject: ")
	content.WriteString(email.Subject)
	content.WriteString("\r\n")

	content.WriteString("From: ")
	content.WriteString(email.From)
	content.WriteString("\r\n")

	content.WriteString("To: ")
	content.WriteString(email.To[0])
	content.WriteString("\r\n")

	content.WriteString(email.Body)
	content.WriteString("\r\n")

	msg := content.Bytes()

	return smtp.SendMail(addr, nil, email.From, email.To, msg)
}

func sendEmailWithResend(email Email) error {
	client := resend.NewClient(email.APIKey)

	params := &resend.SendEmailRequest{
		From:    email.From,
		To:      email.To,
		Subject: email.Subject,
		Html:    email.Body,
	}

	_, err := client.Emails.Send(params)

	return err
}
