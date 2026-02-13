package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/Simpolette/HeartSteal/server/internal/domain"
	"github.com/Simpolette/HeartSteal/server/internal/domain/mocks"
	"github.com/Simpolette/HeartSteal/server/internal/usecase"
	tokenutil "github.com/Simpolette/HeartSteal/server/utils"
)

// ──────────────────────────────────────────
// Test helpers
// ──────────────────────────────────────────

type testDeps struct {
	userRepo    *mocks.MockUserRepository
	sessionRepo *mocks.MockSessionRepository
	pinRepo     *mocks.MockVerificationPINRepository
	emailSvc    *mocks.MockEmailService
	uc          domain.UserUsecase
}

func setup() *testDeps {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	pinRepo := new(mocks.MockVerificationPINRepository)
	emailSvc := new(mocks.MockEmailService)

	uc := usecase.NewUserUseCase(
		userRepo, sessionRepo, pinRepo, emailSvc,
		2*time.Second,
		"test_access_secret", 1,
		"test_refresh_secret", 168,
	)

	return &testDeps{userRepo, sessionRepo, pinRepo, emailSvc, uc}
}

// ──────────────────────────────────────────
// Register Tests
// ──────────────────────────────────────────

func TestUserUseCase_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		d := setup()
		user := &domain.User{
			Username: "test",
			Email:    "new@example.com",
			Password: "password123",
		}

		d.userRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, domain.ErrUserNotFound)
		d.userRepo.On("GetByUsername", mock.Anything, "test").Return(nil, domain.ErrUserNotFound)
		d.userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == "new@example.com" && u.Username == "test" && u.Password != "password123"
		})).Return(nil)

		err := d.uc.Register(context.Background(), user)

		assert.NoError(t, err)
		d.userRepo.AssertExpectations(t)
	})

	t.Run("ErrorEmailExists", func(t *testing.T) {
		d := setup()
		user := &domain.User{Email: "existing@example.com", Password: "123"}

		existingUser := &domain.User{Email: "existing@example.com"}
		d.userRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

		err := d.uc.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrEmailExists, err)
		d.userRepo.AssertNotCalled(t, "Create")
	})

	t.Run("ErrorUsernameExists", func(t *testing.T) {
		d := setup()
		user := &domain.User{Email: "new@example.com", Username: "exist", Password: "123"}

		d.userRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, domain.ErrUserNotFound)
		d.userRepo.On("GetByUsername", mock.Anything, "exist").Return(&domain.User{Username: "exist"}, nil)

		err := d.uc.Register(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUsernameExists, err)
		d.userRepo.AssertNotCalled(t, "Create")
	})
}

// ──────────────────────────────────────────
// Login Tests
// ──────────────────────────────────────────

func TestUserUseCase_Login(t *testing.T) {
	plainPass := "secret123"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(plainPass), 10)
	hashedPass := string(hashedBytes)

	t.Run("Success", func(t *testing.T) {
		d := setup()
		username := "test"

		foundUser := &domain.User{
			Username: username,
			Password: hashedPass,
		}

		d.userRepo.On("GetByUsername", mock.Anything, username).Return(foundUser, nil)
		d.sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)

		accessToken, refreshToken, err := d.uc.Login(context.Background(), username, plainPass)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		d.userRepo.AssertExpectations(t)
		d.sessionRepo.AssertExpectations(t)
	})

	t.Run("ErrorUserNotFound", func(t *testing.T) {
		d := setup()

		d.userRepo.On("GetByUsername", mock.Anything, "ghost").Return(nil, domain.ErrUserNotFound)

		accessToken, refreshToken, err := d.uc.Login(context.Background(), "ghost", "anyPass")

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})

	t.Run("ErrorWrongPassword", func(t *testing.T) {
		d := setup()
		username := "testErr"

		foundUser := &domain.User{
			Username: username,
			Password: hashedPass,
		}
		d.userRepo.On("GetByUsername", mock.Anything, username).Return(foundUser, nil)

		accessToken, refreshToken, err := d.uc.Login(context.Background(), username, "wrong_password")

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})
}

// ──────────────────────────────────────────
// ForgotPassword Tests
// ──────────────────────────────────────────

func TestUserUseCase_ForgotPassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		d := setup()
		email := "user@example.com"

		foundUser := &domain.User{Email: email}
		d.userRepo.On("GetByEmail", mock.Anything, email).Return(foundUser, nil)
		d.pinRepo.On("DeleteByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(nil)
		d.pinRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.VerificationPIN")).Return(nil)
		d.emailSvc.On("SendPIN", mock.Anything, email, mock.AnythingOfType("string")).Return(nil)

		err := d.uc.ForgotPassword(context.Background(), email)

		assert.NoError(t, err)
		d.userRepo.AssertExpectations(t)
		d.pinRepo.AssertExpectations(t)
		d.emailSvc.AssertExpectations(t)
	})

	t.Run("ErrorUserNotFound", func(t *testing.T) {
		d := setup()
		email := "nobody@example.com"

		d.userRepo.On("GetByEmail", mock.Anything, email).Return(nil, domain.ErrUserNotFound)

		err := d.uc.ForgotPassword(context.Background(), email)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserNotFound, err)
		d.pinRepo.AssertNotCalled(t, "Create")
		d.emailSvc.AssertNotCalled(t, "SendPIN")
	})
}

// ──────────────────────────────────────────
// VerifyPIN Tests
// ──────────────────────────────────────────

func TestUserUseCase_VerifyPIN(t *testing.T) {
	pinCode := "123456"
	pinHash, _ := bcrypt.GenerateFromPassword([]byte(pinCode), 10)

	t.Run("Success", func(t *testing.T) {
		d := setup()
		email := "user@example.com"
		foundUser := &domain.User{Email: email}

		storedPIN := &domain.VerificationPIN{
			UserID:    foundUser.ID,
			CodeHash:  string(pinHash),
			Type:      domain.PINTypePasswordReset,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		d.userRepo.On("GetByEmail", mock.Anything, email).Return(foundUser, nil)
		d.pinRepo.On("GetByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(storedPIN, nil)
		d.pinRepo.On("DeleteByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(nil)

		resetToken, err := d.uc.VerifyPIN(context.Background(), email, pinCode)

		assert.NoError(t, err)
		assert.NotEmpty(t, resetToken)
		d.userRepo.AssertExpectations(t)
		d.pinRepo.AssertExpectations(t)
	})

	t.Run("ErrorExpiredPIN", func(t *testing.T) {
		d := setup()
		email := "user@example.com"
		foundUser := &domain.User{Email: email}

		expiredPIN := &domain.VerificationPIN{
			UserID:    foundUser.ID,
			CodeHash:  string(pinHash),
			Type:      domain.PINTypePasswordReset,
			ExpiresAt: time.Now().Add(-5 * time.Minute), // expired
		}

		d.userRepo.On("GetByEmail", mock.Anything, email).Return(foundUser, nil)
		d.pinRepo.On("GetByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(expiredPIN, nil)
		d.pinRepo.On("DeleteByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(nil)

		resetToken, err := d.uc.VerifyPIN(context.Background(), email, pinCode)

		assert.Error(t, err)
		assert.Empty(t, resetToken)
		assert.Equal(t, domain.ErrPINExpired, err)
	})

	t.Run("ErrorWrongPIN", func(t *testing.T) {
		d := setup()
		email := "user@example.com"
		foundUser := &domain.User{Email: email}

		storedPIN := &domain.VerificationPIN{
			UserID:    foundUser.ID,
			CodeHash:  string(pinHash),
			Type:      domain.PINTypePasswordReset,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		d.userRepo.On("GetByEmail", mock.Anything, email).Return(foundUser, nil)
		d.pinRepo.On("GetByUserIDAndType", mock.Anything, foundUser.ID.Hex(), domain.PINTypePasswordReset).Return(storedPIN, nil)

		resetToken, err := d.uc.VerifyPIN(context.Background(), email, "000000")

		assert.Error(t, err)
		assert.Empty(t, resetToken)
		assert.Equal(t, domain.ErrInvalidPIN, err)
	})
}

// ──────────────────────────────────────────
// ResetPassword Tests
// ──────────────────────────────────────────

func TestUserUseCase_ResetPassword(t *testing.T) {
	userID := "507f1f77bcf86cd799439011"

	t.Run("Success", func(t *testing.T) {
		d := setup()

		// Generate a valid reset token using the same secret as the use case
		validToken, err := tokenutil.CreateAccessToken(userID, "test_access_secret", 1)
		assert.NoError(t, err)

		d.userRepo.On("UpdatePassword", mock.Anything, userID, mock.AnythingOfType("string")).Return(nil)
		d.sessionRepo.On("DeleteAllByUserID", mock.Anything, userID).Return(nil)

		err = d.uc.ResetPassword(context.Background(), validToken, "newSecurePass123")

		assert.NoError(t, err)
		d.userRepo.AssertExpectations(t)
		d.sessionRepo.AssertExpectations(t)
	})

	t.Run("ErrorInvalidToken", func(t *testing.T) {
		d := setup()

		err := d.uc.ResetPassword(context.Background(), "not.a.valid.token", "newSecurePass123")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidResetToken, err)
		d.userRepo.AssertNotCalled(t, "UpdatePassword")
	})

	t.Run("ErrorExpiredToken", func(t *testing.T) {
		d := setup()

		// Generate a token that is already expired (0 hours = expires immediately with IssuedAt in past)
		expiredToken, err := tokenutil.CreateAccessToken(userID, "test_access_secret", 0)
		assert.NoError(t, err)

		// Wait briefly to ensure expiry
		time.Sleep(10 * time.Millisecond)

		err = d.uc.ResetPassword(context.Background(), expiredToken, "newSecurePass123")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidResetToken, err)
		d.userRepo.AssertNotCalled(t, "UpdatePassword")
	})
}

// ──────────────────────────────────────────
// SignOut Tests
// ──────────────────────────────────────────

func TestUserUseCase_SignOut(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		d := setup()
		userID := "507f1f77bcf86cd799439011"

		d.sessionRepo.On("DeleteAllByUserID", mock.Anything, userID).Return(nil)

		err := d.uc.SignOut(context.Background(), userID)

		assert.NoError(t, err)
		d.sessionRepo.AssertExpectations(t)
	})
}
