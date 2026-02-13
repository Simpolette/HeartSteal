package domain

import (
	"errors"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid access token")
	ErrInvalidResetToken  = errors.New("invalid reset token")
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       		primitive.ObjectID 	 `bson:"_id,omitempty"   json:"id"`
	Username     	string             	 `bson:"username"        json:"username"`
	Email    		string             	 `bson:"email"           json:"email"`
	Password 		string             	 `bson:"password"        json:"-"`
	AvatarUrl		string				 `bson:"avatar_url"      json:"avatar_url"`
	FriendsList 	[]primitive.ObjectID `bson:"friends_list"    json:"friends_list"`
	CreatedAt 		time.Time 			 `bson:"created_at"      json:"created_at"`
	UpdatedAt 		time.Time 			 `bson:"updated_at"      json:"updated_at"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByUsername(c context.Context, username string) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
	UpdatePassword(c context.Context, userID string, hashedPassword string) error
}

type UserUsecase interface {
	Register(c context.Context, user *User) error
	Login(c context.Context, username string, password string) (accessToken string, refreshToken string, err error)
	ForgotPassword(c context.Context, email string) error
	VerifyPIN(c context.Context, email string, pinCode string) (resetToken string, err error)
	ResetPassword(c context.Context, resetToken string, newPassword string) error
	SignOut(c context.Context, userID string) error
}