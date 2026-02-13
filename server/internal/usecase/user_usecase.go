package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	tokenutil "github.com/Simpolette/HeartSteal/server/utils"

	"golang.org/x/crypto/bcrypt"
)

var _ domain.UserUsecase = &userUseCase{}

type userUseCase struct {
	userRepo           domain.UserRepository
	sessionRepo        domain.SessionRepository
	pinRepo            domain.VerificationPINRepository
	emailService       domain.EmailService
	contextTimeout     time.Duration
	accessTokenSecret  string
	accessTokenExpiry  int
	refreshTokenSecret string
	refreshTokenExpiry int
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	pinRepo domain.VerificationPINRepository,
	emailService domain.EmailService,
	timeout time.Duration,
	accessSecret string,
	accessExpiry int,
	refreshSecret string,
	refreshExpiry int,
) domain.UserUsecase {
	return &userUseCase{
		userRepo:           userRepo,
		sessionRepo:        sessionRepo,
		pinRepo:            pinRepo,
		emailService:       emailService,
		contextTimeout:     timeout,
		accessTokenSecret:  accessSecret,
		accessTokenExpiry:  accessExpiry,
		refreshTokenSecret: refreshSecret,
		refreshTokenExpiry: refreshExpiry,
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

func (u *userUseCase) Login(c context.Context, username string, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := tokenutil.CreateAccessToken(user.ID.Hex(), u.accessTokenSecret, u.accessTokenExpiry)
	if err != nil {
		return "", "", domain.ErrInternalServerError
	}

	refreshToken, err := tokenutil.CreateRefreshToken(user.ID.Hex(), u.refreshTokenSecret, u.refreshTokenExpiry)
	if err != nil {
		return "", "", domain.ErrInternalServerError
	}

	// Hash refresh token for secure storage (SHA-256, since JWTs exceed bcrypt's 72-byte limit)
	refreshHash := hashToken(refreshToken)

	session := &domain.Session{
		UserID:           user.ID,
		RefreshTokenHash: string(refreshHash),
		ExpiresAt:        time.Now().Add(time.Hour * time.Duration(u.refreshTokenExpiry)),
	}

	err = u.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", "", domain.ErrInternalServerError
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) ForgotPassword(c context.Context, email string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Delete any existing PINs for this user/type
	_ = u.pinRepo.DeleteByUserIDAndType(ctx, user.ID.Hex(), domain.PINTypePasswordReset)

	// Generate 6-digit PIN
	pinCode, err := generatePIN(domain.PINLength)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Hash PIN before storing
	pinHash, err := bcrypt.GenerateFromPassword([]byte(pinCode), 10)
	if err != nil {
		return domain.ErrInternalServerError
	}

	pin := &domain.VerificationPIN{
		UserID:    user.ID,
		CodeHash:  string(pinHash),
		Type:      domain.PINTypePasswordReset,
		ExpiresAt: time.Now().Add(domain.PINLifetime),
	}

	err = u.pinRepo.Create(ctx, pin)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Send raw PIN via email
	err = u.emailService.SendPIN(ctx, email, pinCode)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (u *userUseCase) VerifyPIN(c context.Context, email string, pinCode string) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	pin, err := u.pinRepo.GetByUserIDAndType(ctx, user.ID.Hex(), domain.PINTypePasswordReset)
	if err != nil {
		return "", domain.ErrInvalidPIN
	}

	// Check expiration
	if time.Now().After(pin.ExpiresAt) {
		// Clean up expired PIN
		_ = u.pinRepo.DeleteByUserIDAndType(ctx, user.ID.Hex(), domain.PINTypePasswordReset)
		return "", domain.ErrPINExpired
	}

	// Compare hashes
	err = bcrypt.CompareHashAndPassword([]byte(pin.CodeHash), []byte(pinCode))
	if err != nil {
		return "", domain.ErrInvalidPIN
	}

	// PIN is valid — delete it so it can't be reused
	_ = u.pinRepo.DeleteByUserIDAndType(ctx, user.ID.Hex(), domain.PINTypePasswordReset)

	// Generate a short-lived reset session token (1 hour)
	resetToken, err := tokenutil.CreateAccessToken(user.ID.Hex(), u.accessTokenSecret, 1)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	return resetToken, nil
}

func (u *userUseCase) ResetPassword(c context.Context, resetToken string, newPassword string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Validate the reset token
	authorized, err := tokenutil.IsAuthorized(resetToken, u.accessTokenSecret)
	if err != nil || !authorized {
		return domain.ErrInvalidResetToken
	}

	userID, err := tokenutil.ExtractIDFromToken(resetToken, u.accessTokenSecret)
	if err != nil {
		return domain.ErrInvalidResetToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return domain.ErrInternalServerError
	}

	err = u.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Invalidate all existing sessions for security
	_ = u.sessionRepo.DeleteAllByUserID(ctx, userID)

	return nil
}

func (u *userUseCase) SignOut(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.sessionRepo.DeleteAllByUserID(ctx, userID)
}

// generatePIN creates a cryptographically secure random numeric PIN of the given length.
func generatePIN(length int) (string, error) {
	pin := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		pin += fmt.Sprintf("%d", n.Int64())
	}
	return pin, nil
}

// hashToken returns a hex-encoded SHA-256 hash of the given token.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}