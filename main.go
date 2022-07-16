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

	userGroup := apiGroup.Group("/user/")

	userGroup.Get("/", server.HandleGetUser)     // for check if exist
	userGroup.Post("/", server.HandleCreateUser) // registration

	teamsGroup := apiGroup.Group("/team/")
	teamsGroup.Get("/:id/", server.HandleGetTeam)
	teamsGroup.Delete("/:id/", server.HandleDeleteTeam)
	teamsGroup.Get("/", server.HandleGetTeams)    // get all user teams
	teamsGroup.Post("/", server.HandleCreateTeam) // create new team

	boardsGroup := apiGroup.Group("/board/")
	//boardsGroup.Get("/:id/") // TODO
	//boardsGroup.Delete("/:id/") // TODO
	boardsGroup.Get("/", server.HandleGetBoards)    // get all team boards
	boardsGroup.Post("/", server.HandleCreateBoard) // create new board

	log.Fatal(app.Listen(":8081"))
}
