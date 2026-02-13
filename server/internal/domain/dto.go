package domain

// ──────────────────────────────────────────
// Auth Request DTOs
// ──────────────────────────────────────────

type SignupRequest struct {
	Username string `json:"username"  binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"  binding:"required,email"`
}

type ChangePasswordRequest struct {
	Username string `json:"username"  binding:"required"`
}

type VerifyPinRequest struct {
	Email   string `json:"email"     binding:"required,email"`
	PinCode string `json:"pin_code"  binding:"required,len=6"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token"   binding:"required"`
	NewPassword string `json:"new_password"  binding:"required,min=8"`
}

// ──────────────────────────────────────────
// Auth Response DTOs
// ──────────────────────────────────────────

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyPinResponse struct {
	ResetToken string `json:"reset_token"`
}
