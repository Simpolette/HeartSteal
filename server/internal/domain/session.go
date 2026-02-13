package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionBlocked  = errors.New("session is blocked")
)

const (
	CollectionSession = "sessions"
)

type Session struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"        json:"id"`
	UserID           primitive.ObjectID `bson:"user_id"              json:"user_id"`
	RefreshTokenHash string             `bson:"refresh_token_hash"   json:"-"`
	UserAgent        string             `bson:"user_agent"           json:"user_agent"`
	ClientIP         string             `bson:"client_ip"            json:"client_ip"`
	IsBlocked        bool               `bson:"is_blocked"           json:"is_blocked"`
	ExpiresAt        time.Time          `bson:"expires_at"           json:"expires_at"`
	CreatedAt        time.Time          `bson:"created_at"           json:"created_at"`
}

type SessionRepository interface {
	Create(c context.Context, session *Session) error
	GetByID(c context.Context, id string) (*Session, error)
	GetByUserID(c context.Context, userID string) ([]*Session, error)
	DeleteByID(c context.Context, id string) error
	DeleteAllByUserID(c context.Context, userID string) error
}
