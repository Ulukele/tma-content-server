package main

type RequestUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RequestCreateUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RequestTeams struct {
	Username string `json:"username" validate:"required"`
}

type RequestTeam struct {
	Username string `json:"username" validate:"required"`
	Id       string `json:"id" validate:"required"`
}

type RequestCreateTeam struct {
	Username string `json:"username" validate:"required"`
	TeamName string `json:"name" validate:"required"`
}

type RequestBoards struct {
	TeamID uint `json:"teamID" validate:"required"`
}
