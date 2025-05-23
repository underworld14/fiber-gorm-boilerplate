package main

import (
	"fiber-gorm/internal/config"
	"fiber-gorm/internal/database"
	"fiber-gorm/internal/handlers"
	"fiber-gorm/internal/logger"
	"fiber-gorm/internal/middleware"
	"fiber-gorm/internal/repository"
	"fiber-gorm/internal/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Set default JWT secret if not provided
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "your-default-secret-key-change-in-production"
	}

	// Setup logger
	logger.Setup(cfg)
	log := logger.New("main")
	log.Info().Msg("Starting application...")

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		logger.Fatal(err, "Failed to connect to database")
	}
	log.Info().Msg("Database connected successfully")

	// Setup repositories
	userRepo := repository.NewUserRepository(db)

	// Setup services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(cfg, userRepo)

	// Setup handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Create Fiber app with custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Setup middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.RequestLogger())
	app.Use(limiter.New(limiter.Config{
		Max:               3,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Storage: nil,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Too many request")
		},
	}))

	// Setup API routes
	api := app.Group("/api")

	// User routes
	api.Post("/users", userHandler.CreateUser)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh-token", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Profile routes
	profile := api.Group("/profile", middleware.JWTAuthMiddleware(&cfg))
	profile.Get("/", authHandler.Me)

	// Add health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not found",
		})
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Info().Msgf("Server starting on %s", serverAddr)

	if err := app.Listen(serverAddr); err != nil {
		logger.Fatal(err, "Server failed to start")
	}
}
