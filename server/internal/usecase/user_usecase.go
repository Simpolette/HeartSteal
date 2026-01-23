package usecase

import (
	"context"
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"github.com/Simpolette/HeartSteal/server/utils"

	"golang.org/x/crypto/bcrypt"
)

var _ domain.UserUsecase = &userUseCase{}

type userUseCase struct {
	userRepo          domain.UserRepository
	contextTimeout    time.Duration
	accessTokenSecret string
	accessTokenExpiry int
}

func NewUserUseCase(userRepo domain.UserRepository, timeout time.Duration, secret string, expiry int) domain.UserUsecase {
	return &userUseCase{
		userRepo:          userRepo,
		contextTimeout:    timeout,
		accessTokenSecret: secret,
		accessTokenExpiry: expiry,
	}
}

func (u *userUseCase) Register(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	_, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		return domain.ErrEmailExists
	}

	_, err = u.userRepo.GetByUsername(ctx, user.Username)
	if err == nil {
		return domain.ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return domain.ErrInternalServerError
	}
	user.Password = string(hashedPassword)

	return u.userRepo.Create(ctx, user)
}

func (u *userUseCase) Login(c context.Context, username string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := tokenutil.CreateAccessToken(user.ID.Hex(), u.accessTokenSecret, u.accessTokenExpiry)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	return accessToken, nil
}