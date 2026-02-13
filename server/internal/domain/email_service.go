package domain

import "context"

// EmailService defines the contract for sending transactional emails.
type EmailService interface {
	SendPIN(c context.Context, toEmail string, pinCode string) error
}
