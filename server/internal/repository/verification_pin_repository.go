package repository

import (
	"context"
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type verificationPINRepository struct {
	database   *mongo.Database
	collection string
}

func NewVerificationPINRepository(db *mongo.Database, collection string) domain.VerificationPINRepository {
	return &verificationPINRepository{
		database:   db,
		collection: collection,
	}
}

func (r *verificationPINRepository) Create(c context.Context, pin *domain.VerificationPIN) error {
	collection := r.database.Collection(r.collection)

	pin.CreatedAt = time.Now()

	result, err := collection.InsertOne(c, pin)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		pin.ID = oid
	}

	return nil
}

func (r *verificationPINRepository) GetByUserIDAndType(c context.Context, userID string, pinType string) (*domain.VerificationPIN, error) {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	var pin domain.VerificationPIN
	filter := bson.M{
		"user_id": objID,
		"type":    pinType,
	}

	err = collection.FindOne(c, filter).Decode(&pin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrInvalidPIN
		}
		return nil, err
	}

	return &pin, nil
}

func (r *verificationPINRepository) DeleteByUserIDAndType(c context.Context, userID string, pinType string) error {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	filter := bson.M{
		"user_id": objID,
		"type":    pinType,
	}

	_, err = collection.DeleteMany(c, filter)
	return err
}
