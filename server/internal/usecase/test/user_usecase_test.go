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
)

func Test_Register(t *testing.T) {
	// Setup
	setup := func() (*mocks.MockUserRepository, domain.UserUsecase) {
        mockRepo := new(mocks.MockUserRepository)
        timeout := 2 * time.Second
        u := usecase.NewUserUseCase(mockRepo, timeout, "secret", 3600)
        return mockRepo, u
    }
	
	t.Run("Success", func(t *testing.T) {
		mockRepo, u := setup()
		user := &domain.User{
			Username: "test",
			Email:    "new@example.com",
			Password: "password123",
		}

		mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, domain.ErrUserNotFound)

		mockRepo.On("GetByUsername", mock.Anything, "test").Return(nil, domain.ErrUserNotFound)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			// verify password was hashed (it shouldn't match the plain text anymore)
			return u.Email == "new@example.com" && u.Username == "test" && u.Password != "password123"
		})).Return(nil)

		// Execute
		err := u.Register(context.Background(), user)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_EmailExists", func(t *testing.T) {
		mockRepo, u := setup()
		user := &domain.User{Email: "existing@example.com", Password: "123"}

		existingUser := &domain.User{Email: "existing@example.com"}
		mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

		// Execute
		err := u.Register(context.Background(), user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrEmailExists, err)
		// Ensure Create was NEVER called
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Error_UsernameExists", func(t *testing.T) {
		mockRepo, u := setup()
		user := &domain.User{Email: "existing@example.com", Username: "exist", Password: "123"}

		existingUser := &domain.User{Username: "exist"}

		mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(nil, domain.ErrUserNotFound)

		mockRepo.On("GetByUsername", mock.Anything, "exist").Return(existingUser, nil)

		// Execute
		err := u.Register(context.Background(), user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUsernameExists, err)
		// Ensure Create was NEVER called
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func Test_Login(t *testing.T) {
	setup := func() (*mocks.MockUserRepository, domain.UserUsecase) {
        mockRepo := new(mocks.MockUserRepository)
        timeout := 2 * time.Second
        u := usecase.NewUserUseCase(mockRepo, timeout, "my_secret_key", 3600)
        return mockRepo, u
    }
	// Helper: Pre-hash a password so bcrypt.Compare works
	plainPass := "secret123"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(plainPass), 10)
	hashedPass := string(hashedBytes)

	t.Run("Success", func(t *testing.T) {
		mockRepo, u := setup()
		username := "test"
		
		// Mock returns a user with the REAL hashed password
		foundUser := &domain.User{
			Username:    username,
			Password: hashedPass, 
		}

		mockRepo.On("GetByUsername", mock.Anything, username).Return(foundUser, nil)

		// Execute
		token, err := u.Login(context.Background(), username, plainPass)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token) // JWT should be generated
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_UserNotFound", func(t *testing.T) {
		mockRepo, u := setup()
		username := "ghost"
		
		mockRepo.On("GetByUsername", mock.Anything, username).Return(nil, domain.ErrUserNotFound)

		token, err := u.Login(context.Background(), username, "anyPass")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})

	t.Run("Error_WrongPassword", func(t *testing.T) {
		mockRepo, u := setup()
		username := "testErr"
		
		// User exists
		foundUser := &domain.User{
			Username:    username,
			Password: hashedPass, 
		}
		mockRepo.On("GetByUsername", mock.Anything, username).Return(foundUser, nil)

		// Login with WRONG password
		token, err := u.Login(context.Background(), username, "wrong_password")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
	})
}
