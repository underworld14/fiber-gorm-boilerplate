package tests

import (
	"bytes"
	"encoding/json"
	"fiber-gorm/internal/config"
	"fiber-gorm/internal/database"
	"fiber-gorm/internal/handlers"
	"fiber-gorm/internal/middleware"
	"fiber-gorm/internal/repository"
	"fiber-gorm/internal/services"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/assert"
)

// TestApp contains all dependencies for testing the API
type TestApp struct {
	App         *fiber.App
	Config      config.Config
	AuthSvc     *services.AuthService
	UserSvc     *services.UserService
	UserRepo    *repository.UserRepository
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	TestData    map[string]interface{} // Store test data between test cases
}

// SetupTestApp creates a test instance of the application with a test database
func SetupTestApp(t *testing.T) *TestApp {
	// Load test configuration
	cfg := config.Config{
		Environment: "test",
		DBDriver:    "sqlite",
		DBSource:    ":memory:", // Use in-memory SQLite for tests
		ServerPort:  "3000",
		LogLevel:    "error",
		JWTSecret:   "test-jwt-secret",
	}

	// Connect to test database
	db, err := database.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Setup test repositories
	userRepo := &repository.UserRepository{DB: db}

	// Setup test services
	userSvc := &services.UserService{Repo: userRepo}
	authSvc := &services.AuthService{
		Cfg:      cfg, // Pass the config directly (not a pointer)
		UserRepo: userRepo,
	}

	// Setup test handlers
	userHandler := &handlers.UserHandler{Svc: userSvc}
	authHandler := &handlers.AuthHandler{AuthSvc: authSvc}

	// Create test Fiber app with required settings for testing
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		// Disable startup message for tests
		DisableStartupMessage: true,
	})

	// Setup middleware
	app.Use(recover.New())

	// Setup routes
	api := app.Group("/api")

	// User routes
	api.Post("/users", userHandler.CreateUser)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes - match the structure in main.go
	protected := api.Group("/") 
	protected.Use(middleware.JWTAuthMiddleware(&cfg))
	protected.Get("me", authHandler.Me) // Path is /api/me

	return &TestApp{
		App:         app,
		Config:      cfg,
		AuthSvc:     authSvc,
		UserSvc:     userSvc,
		UserRepo:    userRepo,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		TestData:    make(map[string]interface{}), // Initialize test data storage
	}
}

// ExecuteRequest is kept for backward compatibility but you should use MakeRequest instead
func (ta *TestApp) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	// Create a response recorder
	resp := httptest.NewRecorder()

	// Execute the request and get the response
	response, err := ta.App.Test(req)
	if err != nil {
		panic(err)
	}

	// Copy response data to our recorder for easy access
	resp.Code = response.StatusCode
	
	// Copy response body
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	resp.Write(body)

	return resp
}

// MakeRequest is a helper to create and execute requests in one step
func (ta *TestApp) MakeRequest(method, url string, body interface{}, token string) (*http.Response, error) {
	// Marshal the body to JSON if it's not nil
	var reqBody string
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = string(jsonBody)
	}

	// Create request
	req := httptest.NewRequest(method, url, bytes.NewBufferString(reqBody))
	
	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Set authorization header if token is provided
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Execute the request directly using Fiber's test method
	resp, err := ta.App.Test(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ParseResponse parses the JSON response into the provided struct
func ParseResponse(t *testing.T, resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	
	// Log the response body for debugging
	t.Logf("Response body: %s", string(body))
	
	if len(body) > 0 {
		err = json.Unmarshal(body, v)
		assert.NoError(t, err, "Failed to parse response JSON")
	}
}
