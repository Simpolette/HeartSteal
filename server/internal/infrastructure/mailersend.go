package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"github.com/mailersend/mailersend-go"
)

var _ domain.EmailService = &MailerSendService{}

type MailerSendService struct {
	client    *mailersend.Mailersend
	fromEmail string
	fromName  string
}

func NewMailerSendService(apiKey string, fromEmail string) *MailerSendService {
	ms := mailersend.NewMailersend(apiKey)
	return &MailerSendService{
		client:    ms,
		fromEmail: fromEmail,
		fromName:  "HeartSteal",
	}
}

func (s *MailerSendService) SendPIN(c context.Context, toEmail string, pinCode string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	subject := "Your Verification PIN Code"
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto; padding: 24px;">
			<h2 style="color: #333;">HeartSteal Verification</h2>
			<p>Your PIN code is:</p>
			<div style="background: #f4f4f4; padding: 16px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 8px; border-radius: 8px; margin: 16px 0;">
				%s
			</div>
			<p style="color: #666; font-size: 14px;">This code expires in 15 minutes. If you did not request this, please ignore this email.</p>
		</div>
	`, pinCode)
	textBody := fmt.Sprintf("Your HeartSteal verification PIN code is: %s. This code expires in 15 minutes.", pinCode)

	from := mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	}

	recipients := []mailersend.Recipient{
		{
			Email: toEmail,
		},
	}

	message := s.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(htmlBody)
	message.SetText(textBody)

	_, err := s.client.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send PIN email: %w", err)
	}

	return nil
}
