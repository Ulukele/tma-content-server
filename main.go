package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"strconv"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	apiGroup := app.Group("/api/")

	internalGroup := apiGroup.Group("/internal/", func(c *fiber.Ctx) error {
		c.Locals("internal", true)
		return c.Next()
	}) // for backend
	publicGroup := apiGroup.Group("/public/", func(c *fiber.Ctx) error {
		c.Locals("internal", false)
		return c.Next()
	}) // for frontend

	internalGroup.Get("/user/", server.HandleInternalGetUser) // to receive username and password

	userGroup := publicGroup.Group("/user/")
	userGroup.Post("/", server.HandleCreateUser)

	concreteUserGroup := userGroup.Group("/:userId/", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Params("userId", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect userId")
		}
		c.Locals("userId", uint(userId))
		return c.Next()
	})
	concreteUserGroup.Get("/", server.HandleGetUser)

	teamsGroup := concreteUserGroup.Group("/team/")
	teamsGroup.Get("/", server.HandleGetTeams)    // get all user teams
	teamsGroup.Post("/", server.HandleCreateTeam) // create new team

	concreteTeamGroup := teamsGroup.Group("/:teamId/", func(c *fiber.Ctx) error {
		teamId, err := strconv.Atoi(c.Params("teamId", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect teamId")
		}
		c.Locals("teamId", uint(teamId))
		return c.Next()
	})
	concreteTeamGroup.Get("/", server.HandleGetTeam)
	concreteTeamGroup.Delete("/", server.HandleDeleteTeam)
	concreteTeamGroup.Post("/join/", server.HandleJoinTeam)
	concreteTeamGroup.Post("/leave/", server.HandleLeaveTeam)

	boardsGroup := concreteTeamGroup.Group("/board/")
	boardsGroup.Get("/", server.HandleGetBoards)    // get all team boards
	boardsGroup.Post("/", server.HandleCreateBoard) // create new board

	concreteBoardGroup := boardsGroup.Group("/:boardId/", func(c *fiber.Ctx) error {
		boardId, err := strconv.Atoi(c.Params("boardId", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect boardId")
		}
		c.Locals("boardId", uint(boardId))
		return c.Next()
	})
	concreteBoardGroup.Get("/", server.HandleGetBoard)
	concreteBoardGroup.Delete("/", server.HandleDeleteBoard)

	tasksGroup := concreteBoardGroup.Group("/task/")
	tasksGroup.Get("/", server.HandleGetTasks)    // get all board tasks
	tasksGroup.Post("/", server.HandleCreateTask) // create new task

	concreteTaskGroup := tasksGroup.Group("/:taskId/", func(c *fiber.Ctx) error {
		taskId, err := strconv.Atoi(c.Params("taskId", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect taskId")
		}
		c.Locals("taskId", uint(taskId))
		return c.Next()
	})
	concreteTaskGroup.Get("/", server.HandleGetTask)
	concreteTaskGroup.Delete("/", server.HandleDeleteTask)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
