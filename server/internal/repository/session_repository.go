package repository

import (
	"context"
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionRepository struct {
	database   *mongo.Database
	collection string
}

func NewSessionRepository(db *mongo.Database, collection string) domain.SessionRepository {
	return &sessionRepository{
		database:   db,
		collection: collection,
	}
}

func (r *sessionRepository) Create(c context.Context, session *domain.Session) error {
	collection := r.database.Collection(r.collection)

	session.CreatedAt = time.Now()

	result, err := collection.InsertOne(c, session)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		session.ID = oid
	}

	return nil
}

func (r *sessionRepository) GetByID(c context.Context, id string) (*domain.Session, error) {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrSessionNotFound
	}

	var session domain.Session
	err = collection.FindOne(c, bson.M{"_id": objID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) GetByUserID(c context.Context, userID string) ([]*domain.Session, error) {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	cursor, err := collection.Find(c, bson.M{"user_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var sessions []*domain.Session
	if err = cursor.All(c, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRepository) DeleteByID(c context.Context, id string) error {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrSessionNotFound
	}

	_, err = collection.DeleteOne(c, bson.M{"_id": objID})
	return err
}

func (r *sessionRepository) DeleteAllByUserID(c context.Context, userID string) error {
	collection := r.database.Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	_, err = collection.DeleteMany(c, bson.M{"user_id": objID})
	return err
}
