package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

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
		resp = append(resp, s.SerializeTeam(&teamModel))
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

	res, err := s.SerializeTeamExtended(team)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't serialize team")
	}
	return c.JSON(res)
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

	team, err := s.contentDBEngine.CreateTeam(req.UserId, req.TeamName, req.TeamPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't create team")
	}
	return c.JSON(s.SerializeTeam(team))
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
	return c.JSON(s.SerializeTeam(team))
}

func (s *Server) HandleJoinTeam(c *fiber.Ctx) error {
	log.Printf("handle join team at %s", c.Path())

	req := RequestJoinTeam{}
	req.UserId = c.Locals("userId").(uint)
	req.Id = c.Locals("teamId").(uint)

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect team password")
	}

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.JoinTeam(req.UserId, req.Id, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't join team")
	}
	return c.JSON(s.SerializeTeam(team))
}

func (s *Server) HandleLeaveTeam(c *fiber.Ctx) error {
	log.Printf("handle leave team at %s", c.Path())

	req := RequestTeam{}
	req.UserId = c.Locals("userId").(uint)
	req.Id = c.Locals("teamId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	team, err := s.contentDBEngine.LeaveTeam(req.UserId, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't leave team")
	}
	return c.JSON(s.SerializeTeam(team))
}
