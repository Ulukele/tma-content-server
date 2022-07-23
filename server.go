package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"strconv"
)

// Validator
var validate = validator.New()

type Server struct {
	contentDBEngine *DBEngine
}

func NewServer() (*Server, error) {
	// configure content db engine
	// from environment
	cDBC := DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_NAME"),
		Port:     os.Getenv("POSTGRES_PORT"),
		SSLMode:  "disable",
		Tz:       os.Getenv("POSTGRES_TZ"),
	}

	contentEngine, err := NewDBEngine(cDBC)
	if err != nil {
		return nil, err
	}

	s := &Server{}
	s.contentDBEngine = contentEngine

	// init tables with models
	err = s.contentDBEngine.initTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) StartApp() error {
	app := fiber.New()
	app.Use(cors.New())

	apiGroup := app.Group("/api/v1/content/", func(c *fiber.Ctx) error {
		userId, err := strconv.Atoi(c.Get("UserId", "not a number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect userId")
		}
		c.Locals("userId", uint(userId))
		return c.Next()
	})

	teamsGroup := apiGroup.Group("/team/")
	teamsGroup.Get("/", s.HandleGetTeams)
	teamsGroup.Post("/", s.HandleCreateTeam)

	concreteTeamGroup := teamsGroup.Group("/:teamId/", func(c *fiber.Ctx) error {
		teamId, err := strconv.Atoi(c.Params("teamId", "not a number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect teamId")
		}
		c.Locals("teamId", uint(teamId))
		return c.Next()
	})
	concreteTeamGroup.Get("/", s.HandleGetTeam)
	concreteTeamGroup.Delete("/", s.HandleDeleteTeam)
	concreteTeamGroup.Post("/join/", s.HandleJoinTeam)
	concreteTeamGroup.Post("/leave/", s.HandleLeaveTeam)

	boardsGroup := concreteTeamGroup.Group("/board/")
	boardsGroup.Get("/", s.HandleGetBoards)
	boardsGroup.Post("/", s.HandleCreateBoard)

	concreteBoardGroup := boardsGroup.Group("/:boardId/", func(c *fiber.Ctx) error {
		boardId, err := strconv.Atoi(c.Params("boardId", "not a number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect boardId")
		}
		c.Locals("boardId", uint(boardId))
		return c.Next()
	})
	concreteBoardGroup.Get("/", s.HandleGetBoard)
	concreteBoardGroup.Delete("/", s.HandleDeleteBoard)

	tasksGroup := concreteBoardGroup.Group("/task/")
	tasksGroup.Get("/", s.HandleGetTasks)
	tasksGroup.Post("/", s.HandleCreateTask)

	concreteTaskGroup := tasksGroup.Group("/:taskId/", func(c *fiber.Ctx) error {
		taskId, err := strconv.Atoi(c.Params("taskId", "not a number"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "expect taskId")
		}
		c.Locals("taskId", uint(taskId))
		return c.Next()
	})
	concreteTaskGroup.Get("/", s.HandleGetTask)
	concreteTaskGroup.Delete("/", s.HandleDeleteTask)

	return app.Listen(os.Getenv("LISTEN_ON"))
}
