package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"fiber-gorm/internal/config"
	"fiber-gorm/internal/models"
	"fiber-gorm/internal/repository"
	"fiber-gorm/internal/validators"
)

// Custom token claims to include additional data
type TokenClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

// Error types for authentication
var (
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidToken       = errors.New("Invalid or expired token")
	ErrPasswordMismatch   = errors.New("Passwords do not match")
)

// AuthService handles authentication logic
type AuthService struct {
	Cfg      config.Config
	UserRepo *repository.UserRepository
}

func NewAuthService(cfg config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		Cfg:      cfg,
		UserRepo: userRepo,
	}
}

// CreateTokens generates both access and refresh tokens for a user
func (s *AuthService) CreateTokens(user *models.User) (accessToken string, refreshToken string, err error) {
	// Create access token with custom claims
	accessExp := time.Now().Add(15 * time.Minute)
	accessClaims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		},
		Email:    user.Email,
		Username: user.Name,
		Role:     "user", // You can add role-based auth later
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString([]byte(s.Cfg.JWTSecret))
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign access token")
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token with minimal claims (long-lived)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		Subject:   user.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString([]byte(s.Cfg.JWTSecret))
	if err != nil {
		log.Error().Err(err).Msg("Failed to sign refresh token")
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *AuthService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (s *AuthService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Cfg.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("failed to parse token claims")
	}

	return claims.Subject, nil
}

// HashPassword creates a bcrypt hash of a password
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword checks if the provided password matches the hashed password
func (s *AuthService) ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

// LoginUser authenticates a user and returns access and refresh tokens
func (s *AuthService) LoginUser(payload *models.LoginUserPayload) (user *models.User, accessToken string, refreshToken string, err error) {
	// Find the user by email
	user, err = s.UserRepo.FindUserByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Str("email", payload.Email).Msg("User not found during login")
		return nil, "", "", ErrInvalidCredentials
	}

	// Compare the password with the stored hash
	if err = s.ComparePassword(user.Password, payload.Password); err != nil {
		log.Debug().Err(err).Str("email", payload.Email).Msg("Password mismatch during login")
		return nil, "", "", ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err = s.CreateTokens(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// RegisterUser creates a new user account
func (s *AuthService) RegisterUser(payload *models.CreateUserPayload) (*models.User, error) {
	// Validate the payload
	if err := validators.ValidateUserCreation(payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if user already exists
	_, err := s.UserRepo.FindUserByEmail(payload.Email)
	if err == nil {
		// User already exists
		return nil, errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(payload.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user := models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Hobby:    payload.Hobby,
	}

	// Save the user to the database
	if err := s.UserRepo.CreateUser(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// RefreshTokens generates new access and refresh tokens using a valid refresh token
func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	// Validate the refresh token
	userID, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get the user
	user, err := s.UserRepo.FindUserById(userID)
	if err != nil {
		return "", "", ErrUserNotFound
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.CreateTokens(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to create tokens: %w", err)
	}

	return accessToken, newRefreshToken, nil
}
