// Package main provides the server implementation for the dat board application.
// It includes setting up the server, loading environment variables, and defining routes and middleware.
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/ts22082/dat_board_be/handlers"
	"github.com/ts22082/dat_board_be/middleware"
)

// main is the entry point for the server application.
// It loads environment variables, sets up the Fiber app with necessary middleware, and defines API routes.
func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env: %v: ", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Configure CORS middleware
	// After deploying will restrict to deployed address
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: false,
	}))

	// Use MongoConnect middleware for database connection
	app.Use(middleware.MongoConnect())

	// Set up API routes
	api := app.Group("/api")

	// GitHub login route with logging middleware
	api.Get("/github/gh_login", middleware.Logging, handlers.GhLogin)

	// User route with authentication and logging middleware
	api.Get("/user", middleware.Logging, middleware.VerifyAuth, handlers.GetAuthedUser)

	// A test route with logging and delay middleware
	api.Get("/delay/:delay", middleware.Logging, middleware.Delay, handlers.TestDelay)

	// Creates an Item route with authentication and logging middleware
	api.Post("/item", middleware.Logging, middleware.VerifyAuth, handlers.CreateItem)

	// Get Item route with authentication and logging middleware
	api.Delete("/item/:id", middleware.Logging, middleware.VerifyAuth, handlers.DeleteItem)

	// Get Item route with authentication and logging middleware
	api.Get("/item/:id", middleware.Logging, middleware.VerifyAuth, handlers.GetItem)

	// Get Items route with authentication and logging middleware
	api.Get("/items", middleware.Logging, middleware.VerifyAuth, handlers.GetItems)

	// Start the server on port 8080
	const PORT = ":8080"
	log.Printf("Listening on port %v", PORT)
	app.Listen(PORT)

	// Gracefully shutdown the server
	err = app.Shutdown()

	// Log any errors if graceful exit is unsuccessful
	if err != nil {
		log.Fatal(err)
	}
}
