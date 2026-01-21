package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionFriendRequest = "friend_requests"
)

type RequestStatus string

const (
	StatusPending  RequestStatus = "pending"
	StatusAccepted RequestStatus = "accepted"
	StatusRejected RequestStatus = "rejected"
)

func (s RequestStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusRejected:
		return true
	}
	return false
}

var (
	ErrFriendRequestNotFound = errors.New("friend request not found")
	ErrFriendRequestExists   = errors.New("friend request already sent or users are already friends")
	ErrSelfRequest           = errors.New("cannot send friend request to yourself")
	ErrNotYourRequest        = errors.New("you are not authorized to handle this request")
)

type FriendRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"    json:"id"`
	FromUserID primitive.ObjectID `bson:"from_user_id"     json:"fromUserId"`
	ToUserID   primitive.ObjectID `bson:"to_user_id"       json:"toUserId"`
	Status     RequestStatus      `bson:"status"           json:"status"`     // pending, accepted, rejected
	CreatedAt  time.Time          `bson:"created_at"       json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updated_at"       json:"updatedAt"`
}

type FriendRequestRepository interface {
	Create(c context.Context, request *FriendRequest) error
	GetByID(c context.Context, id string) (*FriendRequest, error)
	FetchSent(c context.Context, userID string) ([]FriendRequest, error)
	FetchPending(c context.Context, userID string) ([]FriendRequest, error)
	UpdateStatus(c context.Context, id string, status string) error
	Delete(c context.Context, id string) error
}

type FriendRequestUsecase interface {
	SendRequest(c context.Context, fromUserID string, toUserID string) error
	AcceptRequest(c context.Context, requestID string, currentUserID string) error
	RejectRequest(c context.Context, requestID string, currentUserID string) error
	GetSentRequests(c context.Context, userID string) ([]FriendRequest, error)
	GetPendingRequests(c context.Context, userID string) ([]FriendRequest, error)
}