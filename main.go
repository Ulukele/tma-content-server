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
	//teamsGroup.Get("/:id/") TODO
	//teamsGroup.Delete("/:id/") TODO
	teamsGroup.Get("/", server.HandleGetTeams)
	teamsGroup.Post("/", server.HandleCreateTeam)

	log.Fatal(app.Listen(":8081"))
}
