package repository

import (
	"context"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (r *userRepository) Create(c context.Context, user *domain.User) error {
	collection := r.database.Collection(r.collection)

	result, err := collection.InsertOne(c, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

func (r *userRepository) GetByEmail(c context.Context, email string) (*domain.User, error) {
	collection := r.database.Collection(r.collection)

	var user domain.User

	filter := bson.M{"email": email}

	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}


func (r *userRepository) GetByUsername(c context.Context, username string) (*domain.User, error) {
	collection := r.database.Collection(r.collection)

	var user domain.User

	filter := bson.M{"username": username}

	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}


func (r *userRepository) GetByID(c context.Context, id string) (*domain.User, error) {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	var user domain.User

	filter := bson.M{"_id": objID}
	
	err = collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdatePassword(c context.Context, userID string, hashedPassword string) error {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}

	_, err = collection.UpdateOne(c, filter, update)
	return err
}