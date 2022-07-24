package main

type User struct {
	Id uint `json:"id"`
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
	Id         uint   `json:"id"`
	Title      string `json:"title"`
	Solved     bool   `json:"solved"`
	Importance uint   `json:"importance"`
	BoardId    uint   `json:"boardId"`
	WorkerId   uint   `json:"workerId"`
}
