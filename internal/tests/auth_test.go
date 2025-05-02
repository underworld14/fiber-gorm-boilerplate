package tests

import (
	"fiber-gorm/internal/handlers"
	"fiber-gorm/internal/models"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Seed the random number generator for test emails
	rand.Seed(time.Now().UnixNano())
}

// Generate a unique email for tests
func randomEmail() string {
	return fmt.Sprintf("test-%d@example.com", rand.Int())
}

func TestAuthFlow(t *testing.T) {
	// Setup test application
	app := SetupTestApp(t)

	// Test registration
	t.Run("Register User", func(t *testing.T) {
		// Create registration payload with random email
		testEmail := randomEmail()
		t.Logf("Using test email: %s", testEmail)
		app.TestData["registered_email"] = testEmail

		payload := models.CreateUserPayload{
			Name:     "Test User",
			Email:    testEmail,
			Password: "Password123!",
		}

		// Make registration request
		resp, err := app.MakeRequest(http.MethodPost, "/api/auth/register", payload, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response
		var tokenResp handlers.TokenResponse
		ParseResponse(t, resp, &tokenResp)

		// Verify token response
		assert.NotEmpty(t, tokenResp.AccessToken)
		assert.NotEmpty(t, tokenResp.RefreshToken)
		assert.Equal(t, "bearer", tokenResp.TokenType)
	})

	// Test login
	t.Run("Login User", func(t *testing.T) {
		// Get the email from the registration test
		testEmail := app.TestData["registered_email"].(string)

		// Create login payload
		payload := models.LoginUserPayload{
			Email:    testEmail,
			Password: "Password123!",
		}

		// Make login request
		resp, err := app.MakeRequest(http.MethodPost, "/api/auth/login", payload, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var tokenResp handlers.TokenResponse
		ParseResponse(t, resp, &tokenResp)

		// Verify token response
		assert.NotEmpty(t, tokenResp.AccessToken)
		assert.NotEmpty(t, tokenResp.RefreshToken)
		assert.Equal(t, "bearer", tokenResp.TokenType)

		// Test accessing protected route with token
		t.Run("Access Protected Route", func(t *testing.T) {
			// Make request to protected endpoint
			meResp, err := app.MakeRequest(http.MethodGet, "/api/me", nil, tokenResp.AccessToken)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, meResp.StatusCode)

			// Parse user response
			var user models.User
			ParseResponse(t, meResp, &user)

			// Verify user details
			assert.Equal(t, "Test User", user.Name)
			assert.Equal(t, testEmail, user.Email)
			assert.Empty(t, user.Password) // Password should not be returned
		})

		// Test token refresh
		t.Run("Refresh Token", func(t *testing.T) {
			// Create refresh token payload
			refreshPayload := map[string]string{
				"refresh_token": tokenResp.RefreshToken,
			}

			// Make refresh request
			refreshResp, err := app.MakeRequest(http.MethodPost, "/api/auth/refresh", refreshPayload, "")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, refreshResp.StatusCode)

			// Parse response
			var newTokenResp handlers.TokenResponse
			ParseResponse(t, refreshResp, &newTokenResp)

			// Verify new tokens
			assert.NotEmpty(t, newTokenResp.AccessToken)
			assert.NotEmpty(t, newTokenResp.RefreshToken)
			assert.NotEqual(t, tokenResp.AccessToken, newTokenResp.AccessToken)
		})
	})

	// Test invalid login
	t.Run("Invalid Login", func(t *testing.T) {
		// Get the email from the registration test
		testEmail := app.TestData["registered_email"].(string)

		// Create invalid login payload
		payload := models.LoginUserPayload{
			Email:    testEmail,
			Password: "WrongPassword123!",
		}

		// Make login request
		resp, err := app.MakeRequest(http.MethodPost, "/api/auth/login", payload, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test unauthorized access
	t.Run("Unauthorized Access", func(t *testing.T) {
		// Try to access protected route without token
		resp, err := app.MakeRequest(http.MethodGet, "/api/me", nil, "")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
