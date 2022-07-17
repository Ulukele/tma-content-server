package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

// Validator
var validate = validator.New()

type Server struct {
	contentDBEngine *DBEngine
}

func NewServer() (*Server, error) {
	// configure content db engine
	// from environment
	//cDBC := DBConfig{
	//	Host:     os.Getenv("POSTGRES_HOST"),
	//	User:     os.Getenv("POSTGRES_USER"),
	//	Password: os.Getenv("POSTGRES_PASSWORD"),
	//	Name:     os.Getenv("POSTGRES_NAME"),
	//	Port:     os.Getenv("POSTGRES_PORT"),
	//	SSLMode:  "disable",
	//	Tz:       os.Getenv("POSTGRES_TZ"),
	//}
	cDBC := DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Name:     "postgres",
		Port:     "5432",
		SSLMode:  "disable",
		Tz:       "Asia/Novosibirsk",
	}

	log.Print("Try to create content db engine")
	contentEngine, err := NewDBEngine(cDBC)
	if err != nil {
		return nil, err
	}
	log.Print("Create content db engine")

	s := &Server{}
	s.contentDBEngine = contentEngine

	// init tables with models
	err = s.contentDBEngine.initTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// service user handlers

func (s *Server) HandleInternalGetUser(c *fiber.Ctx) error {
	log.Printf("handle internal get user at %s", c.Path())
	if !c.Locals("internal").(bool) {
		return fiber.NewError(fiber.StatusNotFound, "call internal method not from internal path")
	}

	username := c.Get("Username", "")
	if username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "specify username in handler")
	}

	user := UserFull{}
	userModel, err := s.contentDBEngine.InternalGetUser(username)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("no such user %s", username))
	}

	user.Id = userModel.Id
	user.Username = userModel.Username
	user.Password = userModel.Password

	return c.JSON(user)
}

func (s *Server) HandleGetUser(c *fiber.Ctx) error {
	log.Printf("handle get user at %s", c.Path())

	req := RequestUser{}
	req.UserId = c.Locals("userId").(uint)

	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	user, err := s.contentDBEngine.GetUser(req.UserId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get user info")
	}

	ur := User{
		Id:       user.Id,
		Username: user.Username,
	}

	return c.JSON(ur)
}

func (s *Server) HandleCreateUser(c *fiber.Ctx) error {
	log.Printf("handle create user at %s", c.Path())

	req := RequestCreateUser{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and password")
	}
	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	user, err := s.contentDBEngine.CreateUser(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get create user")
	}

	return c.JSON(User{
		Id:       user.Id,
		Username: user.Username,
	})
}

func (s *Server) HandleGetTeams(c *fiber.Ctx) error {
	log.Printf("handle get teams at %s", c.Path())

	req := RequestTeams{}
	req.UserId = c.Locals("userId").(uint)

	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	teams, err := s.contentDBEngine.GetTeams(req.UserId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get teams")
	}

	resp := make([]Team, 0)

	for _, teamModel := range teams {
		resp = append(resp, Team{
			Id:      teamModel.Id,
			Name:    teamModel.Name,
			OwnerId: teamModel.OwnerId,
		})
	}

	return c.JSON(resp)
}

func (s *Server) HandleGetTeam(c *fiber.Ctx) error {
	log.Printf("handle get team at %s", c.Path())

	req := RequestTeam{}
	req.UserId = c.Locals("userId").(uint)
	req.Id = c.Locals("teamId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.GetTeam(req.UserId, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get team")
	}

	resp := Team{
		Id:      team.Id,
		Name:    team.Name,
		OwnerId: team.OwnerId,
	}
	return c.JSON(resp)
}

func (s *Server) HandleCreateTeam(c *fiber.Ctx) error {
	log.Printf("handle create team at %s", c.Path())

	req := RequestCreateTeam{}
	req.UserId = c.Locals("userId").(uint)

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect team name")
	}
	err := validate.Struct(req)
	if err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.CreateTeam(req.UserId, req.TeamName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't create team")
	}

	resp := Team{
		Id:      team.Id,
		Name:    team.Name,
		OwnerId: team.OwnerId,
	}

	return c.JSON(resp)
}

func (s *Server) HandleDeleteTeam(c *fiber.Ctx) error {
	log.Printf("handle delete team at %s", c.Path())

	req := RequestTeam{}
	req.UserId = c.Locals("userId").(uint)
	req.Id = c.Locals("teamId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.DeleteTeam(req.UserId, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get team")
	}

	resp := Team{
		Id:      team.Id,
		Name:    team.Name,
		OwnerId: team.OwnerId,
	}
	return c.JSON(resp)
}

func (s *Server) HandleGetBoards(c *fiber.Ctx) error {
	log.Printf("handle get boards at %s", c.Path())

	req := RequestBoards{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	boards, err := s.contentDBEngine.GetBoards(req.UserId, req.TeamID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get boards")
	}

	resp := make([]Board, 0)

	for _, boardModel := range boards {
		resp = append(resp, Board{
			Id:     boardModel.Id,
			Name:   boardModel.Name,
			TeamId: boardModel.TeamId,
		})
	}

	return c.JSON(resp)
}

func (s *Server) HandleCreateBoard(c *fiber.Ctx) error {
	log.Printf("handle create board at %s", c.Path())

	req := RequestCreateBoard{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect name")
	}

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	board, err := s.contentDBEngine.CreateBoard(req.UserId, req.TeamID, req.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't create board")
	}

	resp := Board{
		Id:     board.Id,
		Name:   board.Name,
		TeamId: board.TeamId,
	}

	return c.JSON(resp)
}
