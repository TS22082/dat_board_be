package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/ts22082/dat_board_be/handlers"
	"github.com/ts22082/dat_board_be/middleware"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Printf("Failed to load .env: %v: ", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: false,
	}))

	app.Use(middleware.MongoConnect())

	api := app.Group("/api")

	api.Get("/github/gh_login", middleware.Logging, handlers.GhLogin)

	api.Get("/user", middleware.Logging, middleware.VerifyAuth, handlers.GetAuthedUser)

	const PORT = ":8080"
	log.Printf("Listening on port %v", PORT)
	app.Listen(PORT)
}
