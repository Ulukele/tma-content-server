package main

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type UserFull struct {
	User
	Password string `json:"password"`
}

type Team struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	OwnerId  uint   `json:"ownerId"`
	Password string `json:"password"`
}

type TeamExtended struct {
	Team
	Users []User `json:"users"`
}

type Board struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	TeamId uint   `json:"teamId"`
}

type BoardExtended struct {
	Board
	Tasks []Task
}

type Task struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Solved  bool   `json:"solved"`
	BoardId uint   `json:"boardId"`
}
