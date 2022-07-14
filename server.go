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
	contentDBEngine  *DBEngine
	sessionsDBEngine *DBEngine
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

	// configure content db engine
	// from environment
	//sDBC := DBConfig{
	//	Host:     os.Getenv("POSTGRES_S_HOST"),
	//	User:     os.Getenv("POSTGRES_S_USER"),
	//	Password: os.Getenv("POSTGRES_S_PASSWORD"),
	//	Name:     os.Getenv("POSTGRES_S_NAME"),
	//	Port:     os.Getenv("POSTGRES_S_PORT"),
	//	SSLMode:  "disable",
	//	Tz:       os.Getenv("POSTGRES_S_TZ"),
	//}
	sDBC := DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Name:     "postgres",
		Port:     "5432",
		SSLMode:  "disable",
		Tz:       "Asia/Novosibirsk",
	}

	log.Print("Try to create sessions db engine")
	sessionsEngine, err := NewDBEngine(sDBC)
	if err != nil {
		return nil, err
	}
	log.Print("Create sessions db engine")

	s := &Server{}
	s.contentDBEngine = contentEngine
	s.sessionsDBEngine = sessionsEngine

	// init tables with models
	err = s.contentDBEngine.initTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) fetchUser(username string) error {

	// check if already in content db
	var serviceUser = ServiceUserModel{Username: username}
	var exists bool
	err := s.contentDBEngine.DB.Model(&serviceUser).
		Select("count(*) > 0").
		Where("Username = ?", username).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// check if in sessions db
	var user = UserModel{Username: username}
	err = s.sessionsDBEngine.DB.Model(&user).
		Select("count(*) > 0").
		Where("Username = ?", username).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		serviceUser := ServiceUserModel{Username: username}
		if err = s.contentDBEngine.DB.Save(&serviceUser).Error; err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("no such user %s", username)
	}
}

// service user handlers

func (s *Server) HandleGetUser(c *fiber.Ctx) error {
	log.Printf("handle get user at %s", c.Path())

	req := RequestUser{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect username and sessionID")
	}
	err := validate.Struct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	if err = s.fetchUser(req.Username); err != nil {
		log.Printf("error while fetch: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "can't find such user")
	}

	user, err := s.contentDBEngine.GetUserInfo(req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get user info")
	}

	// looks like useless (get username by username)
	// but may be added new fields (avatar, status, ...)
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

	if err = s.fetchUser(req.Username); err != nil {
		log.Printf("error while fetch: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "can't find such user")
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

	if err = s.fetchUser(req.Username); err != nil {
		log.Printf("error while fetch: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "can't find such user")
	}

	team, err := s.contentDBEngine.CreateTeam(req.Username, req.TeamName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't create team")
	}

	resp := Team{Id: team.Id, Name: team.Name}

	return c.JSON(resp)
}
