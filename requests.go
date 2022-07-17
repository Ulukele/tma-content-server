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

type RequestCreateTeam struct {
	RequestFrom
	TeamName string `json:"name" validate:"required"`
}

type RequestBoards struct {
	RequestFrom
	TeamID uint `json:"teamId" validate:"required"`
}

type RequestCreateBoard struct {
	RequestFrom
	TeamID uint   `json:"teamId" validate:"required"`
	Name   string `json:"name" validate:"required"`
}
