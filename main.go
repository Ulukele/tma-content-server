package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	apiGroup := app.Group("/api/")

	apiGroup.Get("/user/", server.HandleGetUser)

	teamsGroup := apiGroup.Group("/team/")
	teamsGroup.Get("/:id/")
	teamsGroup.Delete("/:id/")
	teamsGroup.Get("/")
	teamsGroup.Post("/")

	log.Fatal(app.Listen(":8081"))
}
