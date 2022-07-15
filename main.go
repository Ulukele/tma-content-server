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
	internalGroup := apiGroup.Group("/internal/") // don't forward

	internalGroup.Get("/user/", server.HandleGetUser) // for check if exist
	apiGroup.Post("/user/", server.HandleCreateUser)  // registration

	teamsGroup := apiGroup.Group("/team/")
	teamsGroup.Get("/:id/", server.HandleGetTeam)
	//teamsGroup.Delete("/:id/") TODO
	teamsGroup.Get("/", server.HandleGetTeams)    // get all user teams
	teamsGroup.Post("/", server.HandleCreateTeam) // create new team

	log.Fatal(app.Listen(":8081"))
}
