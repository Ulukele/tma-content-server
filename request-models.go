package main

// RequestFrom all requests contains user
// that requires that data
type RequestFrom struct {
	UserId uint `json:"userId" validate:"required"`
}

type RequestUser struct {
	RequestFrom
}

type RequestCreateUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RequestTeams struct {
	RequestFrom
}

type RequestTeam struct {
	RequestFrom
	Id uint `json:"id" validate:"required"`
}

type RequestJoinTeam struct {
	RequestTeam
	Password string `json:"password" validate:"required"`
}

type RequestCreateTeam struct {
	RequestFrom
	TeamName     string `json:"name" validate:"required"`
	TeamPassword string `json:"password" validate:"required"`
}

type RequestBoards struct {
	RequestFrom
	TeamID uint `json:"teamId" validate:"required"`
}

type RequestBoard struct {
	RequestFrom
	Id     uint `json:"id" validate:"required"`
	TeamID uint `json:"teamId" validate:"required"`
}

type RequestCreateBoard struct {
	RequestFrom
	TeamID uint   `json:"teamId" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type RequestTasks struct {
	RequestFrom
	TeamID  uint `json:"teamId" validate:"required"`
	BoardID uint `json:"boardId" validate:"required"`
}

type RequestTask struct {
	RequestFrom
	TeamID  uint `json:"teamId" validate:"required"`
	BoardID uint `json:"boardId" validate:"required"`
	Id      uint `json:"taskId" validate:"required"`
}

type RequestCreateTask struct {
	RequestFrom
	TeamID  uint   `json:"teamId" validate:"required"`
	BoardID uint   `json:"boardId" validate:"required"`
	Title   string `json:"title" validate:"required"`
}
