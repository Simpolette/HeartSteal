package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidPIN = errors.New("invalid PIN code")
	ErrPINExpired = errors.New("PIN code has expired")
)

const (
	CollectionVerificationPIN = "verification_pins"

	PINTypePasswordReset  = "password_reset"
	PINTypePasswordChange = "password_change"

	PINLength   = 6
	PINLifetime = 15 * time.Minute
)

type VerificationPIN struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"   json:"id"`
	UserID    primitive.ObjectID `bson:"user_id"         json:"user_id"`
	CodeHash  string             `bson:"code_hash"       json:"-"`
	Type      string             `bson:"type"            json:"type"`
	ExpiresAt time.Time          `bson:"expires_at"      json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"      json:"created_at"`
}

type VerificationPINRepository interface {
	Create(c context.Context, pin *VerificationPIN) error
	GetByUserIDAndType(c context.Context, userID string, pinType string) (*VerificationPIN, error)
	DeleteByUserIDAndType(c context.Context, userID string, pinType string) error
}
