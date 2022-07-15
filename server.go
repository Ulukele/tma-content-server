package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
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
	//cDBC := DBConfig{
	//	Host:     os.Getenv("POSTGRES_C_HOST"),
	//	User:     os.Getenv("POSTGRES_C_USER"),
	//	Password: os.Getenv("POSTGRES_C_PASSWORD"),
	//	Name:     os.Getenv("POSTGRES_C_NAME"),
	//	Port:     os.Getenv("POSTGRES_C_PORT"),
	//	SSLMode:  "disable",
	//	Tz:       os.Getenv("POSTGRES_C_TZ"),
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

func (s *Server) HandleGetUser(c *fiber.Ctx) error {
	log.Printf("handle get user at %s", c.Path())

	req := RequestUser{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and password")
	}
	err := validate.Struct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	user, err := s.contentDBEngine.GetUserInfo(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get user info")
	}

	type UserResp struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ur := UserResp{
		Username: user.Username,
		Password: user.Password,
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
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	user, err := s.contentDBEngine.CreateUser(req.Username, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get create user")
	}

	type UserResp struct {
		Username string `json:"username"`
	}

	ur := UserResp{
		Username: user.Username,
	}

	return c.JSON(ur)
}

func (s *Server) HandleGetTeams(c *fiber.Ctx) error {
	log.Printf("handle get teams at %s", c.Path())

	req := RequestTeams{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username")
	}
	err := validate.Struct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	teams, err := s.contentDBEngine.GetTeams(req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get teams")
	}

	type ResponseTeams struct {
		Teams []Team
	}
	teamsResp := make([]Team, 0)

	for _, teamModel := range teams {
		teamsResp = append(teamsResp, Team{Id: teamModel.Id, Name: teamModel.Name})
	}

	return c.JSON(ResponseTeams{Teams: teamsResp})
}

func (s *Server) HandleGetTeam(c *fiber.Ctx) error {
	log.Printf("handle get team at %s", c.Path())

	req := RequestTeam{}
	req.Id = c.Params("id", "")

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username")
	}

	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	teamId, err := strconv.Atoi(req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid teamId")
	}
	team, err := s.contentDBEngine.GetTeam(req.Username, uint(teamId))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get team")
	}

	resp := Team{Id: team.Id, Name: team.Name}
	return c.JSON(resp)
}

func (s *Server) HandleCreateTeam(c *fiber.Ctx) error {
	log.Printf("handle create team at %s", c.Path())

	req := RequestCreateTeam{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and team name")
	}
	err := validate.Struct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.CreateTeam(req.Username, req.TeamName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't create team")
	}

	resp := Team{Id: team.Id, Name: team.Name}

	return c.JSON(resp)
}
