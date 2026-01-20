package domain

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       		primitive.ObjectID 	 `bson:"_id,omitempty"   json:"id"`
	Username     	string             	 `bson:"username"        json:"username"`
	Email    		string             	 `bson:"email"           json:"email"`
	PhoneNumber		string				 `bson:"phone_number"    json:"phoneNumber"`
	Password 		string             	 `bson:"password"        json:"-"`
	AvatarUrl		string				 `bson:"avatar_url"      json:"avatar_url"`
	FriendsList 	[]primitive.ObjectID `bson:"friends_list"    json:"friends_list"`
	CreatedAt 		time.Time 			 `bson:"created_at"      json:"created_at"`
	UpdatedAt 		time.Time 			 `bson:"updated_at"      json:"updated_at"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (*User, error)
	GetByID(c context.Context, id string) (*User, error)
}

type UserUsecase interface {
	Register(c context.Context, user *User) error
	Login(c context.Context, email string, password string) (string, error)
}