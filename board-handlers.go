package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

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
		res := s.SerializeBoard(&boardModel)
		resp = append(resp, res)
	}

	return c.JSON(resp)
}

func (s *Server) HandleGetBoard(c *fiber.Ctx) error {
	log.Printf("handle get board at %s", c.Path())

	req := RequestBoard{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.Id = c.Locals("boardId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	board, err := s.contentDBEngine.GetBoard(req.UserId, req.TeamID, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get board")
	}

	res, err := s.SerializeBoardExtended(board)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "can't serialize board")
	}
	return c.JSON(res)
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

	return c.JSON(s.SerializeBoard(board))
}

func (s *Server) HandleDeleteBoard(c *fiber.Ctx) error {
	log.Printf("handle delete board at %s", c.Path())

	req := RequestBoard{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.Id = c.Locals("boardId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	_, err := s.contentDBEngine.DeleteBoard(req.UserId, req.TeamID, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't delete board")
	}

	return c.JSON("")
}
