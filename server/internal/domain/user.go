package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       		primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName     	string             `bson:"name"          json:"name"`
	Email    		string             `bson:"email"         json:"email"`
	
	Password 		string             `bson:"password"      json:"-"`
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