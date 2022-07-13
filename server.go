package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
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
	cDBC := DBConfig{
		Host:     os.Getenv("POSTGRES_C_HOST"),
		User:     os.Getenv("POSTGRES_C_USER"),
		Password: os.Getenv("POSTGRES_C_PASSWORD"),
		Name:     os.Getenv("POSTGRES_C_NAME"),
		Port:     os.Getenv("POSTGRES_C_PORT"),
		Tz:       os.Getenv("POSTGRES_C_TZ"),
	}

	contentEngine, err := NewDBEngine(cDBC)
	if err != nil {
		return nil, err
	}

	// configure content db engine
	// from environment
	sDBC := DBConfig{
		Host:     os.Getenv("POSTGRES_S_HOST"),
		User:     os.Getenv("POSTGRES_S_USER"),
		Password: os.Getenv("POSTGRES_S_PASSWORD"),
		Name:     os.Getenv("POSTGRES_S_NAME"),
		Port:     os.Getenv("POSTGRES_S_PORT"),
		Tz:       os.Getenv("POSTGRES_S_TZ"),
	}

	sessionsEngine, err := NewDBEngine(sDBC)
	if err != nil {
		return nil, err
	}

	s := &Server{}
	s.contentDBEngine = contentEngine
	s.sessionsDBEngine = sessionsEngine
	return s, nil
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

	if ok := s.sessionsDBEngine.CheckSessionId(req.Username, req.SessionID); !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user or sessionID")
	}

	user, err := s.contentDBEngine.GetUserInfo(req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't get user info")
	}

	type UserResp struct {
		Username string   `json:"username"`
		Tables   []string `json:"tables"`
	}

	ur := UserResp{
		Username: user.Username,
		Tables:   []string{},
	}

	return c.JSON(ur)
}
