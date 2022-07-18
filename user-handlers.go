package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

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

	return c.JSON(s.SerializeUser(user))
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

	return c.JSON(s.SerializeUser(user))
}
