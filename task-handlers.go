package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func (s *Server) HandleGetTasks(c *fiber.Ctx) error {
	log.Printf("handle get tasks at %s", c.Path())

	req := RequestTasks{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.BoardID = c.Locals("boardId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	tasks, err := s.contentDBEngine.GetTasks(req.UserId, req.TeamID, req.BoardID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get tasks")
	}

	resp := make([]Task, 0)

	for _, taskModel := range tasks {
		res := s.SerializeTask(&taskModel)
		resp = append(resp, res)
	}

	return c.JSON(resp)
}

func (s *Server) HandleCreateTask(c *fiber.Ctx) error {
	log.Printf("handle create task at %s", c.Path())

	req := RequestCreateTask{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.BoardID = c.Locals("boardId").(uint)

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "expect title")
	}

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	taskEntity := Task{Title: req.Title, Importance: req.Importance, WorkerId: req.WorkerId}
	task, err := s.contentDBEngine.CreateTask(req.UserId, req.TeamID, req.BoardID, taskEntity)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't create task")
	}

	return c.JSON(s.SerializeTask(task))
}

func (s *Server) HandleGetTask(c *fiber.Ctx) error {
	log.Printf("handle get task at %s", c.Path())

	req := RequestTask{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.BoardID = c.Locals("boardId").(uint)
	req.Id = c.Locals("taskId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	task, err := s.contentDBEngine.GetTask(req.UserId, req.TeamID, req.BoardID, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't get task")
	}

	return c.JSON(s.SerializeTask(task))
}

func (s *Server) HandleDeleteTask(c *fiber.Ctx) error {
	log.Printf("handle delete task at %s", c.Path())

	req := RequestTask{}
	req.UserId = c.Locals("userId").(uint)
	req.TeamID = c.Locals("teamId").(uint)
	req.BoardID = c.Locals("boardId").(uint)
	req.Id = c.Locals("taskId").(uint)

	if err := validate.Struct(req); err != nil {
		log.Printf("validation error: %s", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "validation error")
	}

	_, err := s.contentDBEngine.DeleteTask(req.UserId, req.TeamID, req.BoardID, req.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "can't delete task")
	}

	return c.JSON("")
}
